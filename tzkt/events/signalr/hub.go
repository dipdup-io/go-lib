package signalr

import (
	"net"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
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
	stop    chan struct{}
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
		stop:    make(chan struct{}, 1),
	}, nil
}

// Connect -
func (hub *Hub) Connect() error {
	if err := hub.handshake(); err != nil {
		return err
	}

	hub.listen()
	return nil
}

func (hub *Hub) handshake() error {
	log.Infof("connecting to %s...", hub.url.String())

	c, _, err := websocket.DefaultDialer.Dial(hub.url.String(), nil)
	if err != nil {
		return errors.Wrap(err, "Connect Dial")
	}
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
	log.Info("connected")

	return nil
}

// Close -
func (hub *Hub) Close() error {
	hub.stop <- struct{}{}
	hub.wg.Wait()

	if err := hub.Send(newCloseMessage()); err != nil {
		return err
	}

	if err := hub.conn.Close(); err != nil {
		return err
	}

	close(hub.stop)
	close(hub.msgs)
	return nil
}

func (hub *Hub) reconnect() error {
	log.Warn("reconnecting...")

	if err := hub.Send(newCloseMessage()); err != nil {
		log.Error(err)
	}

	if err := hub.conn.Close(); err != nil {
		log.Error(err)
	}
	log.Info("connection closed")
	if err := hub.handshake(); err != nil {
		return err
	}
	if hub.onReconnect != nil {
		return hub.onReconnect()
	}
	return nil
}

func (hub *Hub) listen() {
	hub.wg.Add(1)

	go func() {
		defer hub.wg.Done()

		for {
			select {
			case <-hub.stop:
				return
			default:
				if err := hub.readAllMessages(); err != nil {
					switch {
					case errors.Is(err, ErrTimeout) || errors.Is(err, &websocket.CloseError{}):
						if err := hub.reconnect(); err != nil {
							log.Errorf("reconnect: %s", err.Error())
							log.Warn("retry after 5 seconds")
							time.Sleep(time.Second * 5)
						}
					case errors.Is(err, ErrEmptyResponse):
					default:
						log.Errorf("readAllMessages: %s", err.Error())
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
		log.Warn("No messages during read timeout")
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
		log.Error(msg.Error)
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

		return nil, errors.Wrap(err, "NextReader")
	}
	return NewJSONReader(r), nil
}

func getDeadline() time.Time {
	return time.Now().Add(keepAlive)
}
