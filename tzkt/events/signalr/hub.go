package signalr

import (
	"context"
	"net"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	keepAlive = time.Second * 15
)

// Hub -
type Hub struct {
	url  *url.URL
	conn *websocket.Conn

	encoder Encoding
	msgs    chan interface{}
	log     zerolog.Logger
	mx      sync.Mutex
	wg      sync.WaitGroup

	onReconnect func() error
}

// NewHub -
func NewHub(address, connectionToken string) (*Hub, error) {
	u, err := url.Parse(address)
	if err != nil {
		return nil, err
	}
	switch u.Scheme {
	case "https":
		u.Scheme = "wss"
	case "http":
		u.Scheme = "ws"
	default:
		return nil, errors.Wrap(ErrInvalidScheme, u.Scheme)
	}

	return &Hub{
		url:     u,
		encoder: NewJSONEncoding(),
		msgs:    make(chan interface{}, 1024),
		log:     log.Logger,
	}, nil
}

// Connect -
func (hub *Hub) Connect(ctx context.Context) error {
	if err := hub.handshake(); err != nil {
		return err
	}

	hub.listen(ctx)
	return nil
}

func (hub *Hub) handshake() error {
	hub.log.Debug().Msgf("connecting to %s...", hub.url.String())

	c, response, err := websocket.DefaultDialer.Dial(hub.url.String(), nil)
	if err != nil {
		return errors.Wrap(err, "Connect Dial")
	}
	defer response.Body.Close()

	hub.conn = c

	if err := hub.Send(newHandshakeRequest()); err != nil {
		return errors.Wrap(err, "Connect handshake send message")
	}

	var resp Error
	if err := hub.readOneMessage(&resp); err != nil {
		return errors.Wrap(err, "readOneMessage")
	}

	if resp.Error != "" {
		return errors.Wrap(ErrHandshake, resp.Error)
	}
	hub.log.Debug().Msg("connected")

	return nil
}

// Close -
func (hub *Hub) Close() error {
	hub.wg.Wait()

	if err := hub.Send(newCloseMessage()); err != nil {
		return err
	}

	if err := hub.conn.Close(); err != nil {
		return err
	}

	close(hub.msgs)
	return nil
}

func (hub *Hub) reconnect() error {
	hub.log.Warn().Msg("reconnecting...")

	if err := hub.Send(newCloseMessage()); err != nil {
		hub.log.Err(err).Msg("send")
	}

	if err := hub.conn.Close(); err != nil {
		hub.log.Err(err).Msg("close")
	}
	hub.log.Debug().Msg("connection closed")
	if err := hub.handshake(); err != nil {
		return err
	}
	if hub.onReconnect != nil {
		return hub.onReconnect()
	}
	return nil
}

func (hub *Hub) listen(ctx context.Context) {
	hub.wg.Add(1)

	go func() {
		defer hub.wg.Done()

		for {
			select {
			case <-ctx.Done():
				hub.log.Debug().Msg("stop hub listenning...")
				return
			default:
				if err := hub.readAllMessages(); err != nil {
					switch {
					case errors.Is(err, ErrTimeout) || websocket.IsCloseError(err, websocket.CloseAbnormalClosure):
						if err := hub.reconnect(); err != nil {
							hub.log.Err(err).Msg("reconnect")
							hub.log.Warn().Msg("retry after 5 seconds")
							time.Sleep(time.Second * 5)
						}
					case errors.Is(err, ErrEmptyResponse):
					default:
						hub.log.Err(err).Msg("readAllMessages")
					}
				}
			}
		}
	}()
}

// Send - send message
func (hub *Hub) Send(msg interface{}) error {
	data, err := hub.encoder.Encode(msg)
	if err != nil {
		return err
	}

	hub.log.Trace().Str("data", string(data)).Msg("==> TzKT server")

	hub.mx.Lock()
	defer hub.mx.Unlock()
	return hub.conn.WriteMessage(websocket.TextMessage, data)
}

func (hub *Hub) readOneMessage(msg interface{}) error {
	scanner, err := hub.getScanner()
	if err != nil {
		return err
	}
	if scanner == nil {
		return nil
	}
	if err := scanner.Scan(); err != nil {
		return err
	}
	data, err := scanner.Bytes()
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return ErrEmptyResponse
	}
	hub.log.Trace().Str("data", string(data)).Msg("<== TzKT server")

	if err := json.Unmarshal(data, msg); err != nil {
		return err
	}

	if err := hub.conn.SetReadDeadline(getDeadline()); err != nil {
		return errors.Wrap(err, "SetReadDeadline")
	}

	return nil
}

func (hub *Hub) readAllMessages() error {
	scanner, err := hub.getScanner()
	if err != nil {
		return err
	}
	if scanner == nil {
		hub.log.Warn().Msg("no messages during read timeout")
		return ErrEmptyResponse
	}
	if err := scanner.Scan(); err != nil {
		return err
	}

	for {
		data, err := scanner.Bytes()
		if err != nil {
			return err
		}
		if len(data) == 0 {
			break
		}

		hub.log.Trace().Str("data", string(data)).Msg("<== TzKT server")

		msg, err := hub.encoder.Decode(data)
		if err != nil {
			return err
		}
		hub.msgs <- msg

		if closeMsg, ok := msg.(CloseMessage); ok {
			return hub.closeMessageHandler(closeMsg)
		}
	}

	if err := hub.conn.SetReadDeadline(getDeadline()); err != nil {
		return errors.Wrap(err, "SetReadDeadline")
	}

	return nil
}

func (hub *Hub) closeMessageHandler(msg CloseMessage) error {
	if msg.Error != "" {
		hub.log.Error().Msg(msg.Error)
	}
	if !msg.AllowReconnect {
		return ErrConnectionClose
	}
	return hub.reconnect()
}

func (hub *Hub) getScanner() (*JSONReader, error) {
	_, r, err := hub.conn.NextReader()
	if err != nil {
		if e, ok := err.(net.Error); ok && e.Timeout() {
			return nil, ErrTimeout
		}
		if websocket.IsCloseError(err, websocket.CloseAbnormalClosure) {
			return nil, err
		}

		return nil, errors.Wrap(err, "NextReader")
	}
	return NewJSONReader(r), nil
}

func getDeadline() time.Time {
	return time.Now().Add(keepAlive)
}
