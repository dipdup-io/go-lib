package events

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/dipdup-net/go-lib/tzkt/events/signalr"
	log "github.com/sirupsen/logrus"
)

// TzKT - struct that used for connection to TzKT events server
type TzKT struct {
	s            *signalr.SignalR
	invokationID string

	subscriptions []signalr.Invocation

	msgs chan Message
	stop chan struct{}
	wg   sync.WaitGroup
}

// NewTzKT - constructor of `TzKT`. `url` is TzKT events base URL. If it's empty https://api.tzkt.io/v1/events is set.
func NewTzKT(url string) *TzKT {
	if url == "" {
		url = BaseURL
	}
	return &TzKT{
		s:             signalr.NewSignalR(url),
		invokationID:  fmt.Sprintf("%d", time.Now().UnixNano()),
		msgs:          make(chan Message),
		stop:          make(chan struct{}),
		subscriptions: make([]signalr.Invocation, 0),
	}
}

// Connect - connect to events SignalR server
func (tzkt *TzKT) Connect() error {
	if err := tzkt.s.Connect(signalr.Version1); err != nil {
		return err
	}
	tzkt.s.SetOnReconnect(tzkt.onReconnect)
	tzkt.listen()
	return nil
}

// Close - closing all connections
func (tzkt *TzKT) Close() error {
	tzkt.stop <- struct{}{}
	tzkt.wg.Wait()

	if err := tzkt.s.Close(); err != nil {
		return err
	}
	close(tzkt.msgs)
	close(tzkt.stop)
	return nil
}

// Listen - listen channel with all received messages
func (tzkt *TzKT) Listen() <-chan Message {
	return tzkt.msgs
}

// SubscribeToHead - subscribe to head channel. Sends the blockchain head every time it has been updated.
func (tzkt *TzKT) SubscribeToHead() error {
	return tzkt.subscribe(MethodHead)
}

// SubscribeToBlocks - subscribe to blocks channel. Sends blocks added to the blockchain.
func (tzkt *TzKT) SubscribeToBlocks() error {
	return tzkt.subscribe(MethodBlocks)
}

// SubscribeToOperations - subscribe to operations channel.
// Sends operations of specified types or related to specified accounts, included into the blockchain.
// Filters by `address` and list of `types` is appliable.
func (tzkt *TzKT) SubscribeToOperations(address string, types ...string) error {
	args := make(map[string]interface{})
	if len(types) > 0 {
		args["types"] = strings.Join(types, ",")
	}
	if address != "" {
		args["address"] = address
	}
	if len(args) > 0 {
		return tzkt.subscribe(MethodOperations, args)
	}
	return tzkt.subscribe(MethodOperations)
}

func (tzkt *TzKT) subscribe(channel string, args ...interface{}) error {
	msg := signalr.NewInvocation(tzkt.invokationID, channel, args...)
	tzkt.subscriptions = append(tzkt.subscriptions, msg)
	return tzkt.s.Send(msg)
}

func (tzkt *TzKT) listen() {
	tzkt.wg.Add(1)

	go func() {
		defer tzkt.wg.Done()

		for {
			select {
			case <-tzkt.stop:
				return
			case msg := <-tzkt.s.Messages():
				typ, ok := msg.(signalr.Invocation)
				if !ok {
					continue
				}
				if len(typ.Arguments) == 0 {
					log.Warnf("Empty arguments of invokation: %v", typ)
					continue
				}
				args, ok := typ.Arguments[0].(map[string]interface{})
				if !ok {
					log.Warnf("Invalid arguments type: %v", typ)
					continue
				}
				msgType, ok := args["type"]
				if !ok {
					log.Warnf("Empty tzkt message type: %v", args)
					continue
				}
				msgState, ok := args["state"]
				if !ok {
					log.Warnf("Empty tzkt message state: %v", args)
					continue
				}
				data, ok := args["data"]
				if !ok {
					data = nil
				}

				tzkt.msgs <- Message{
					Channel: typ.Target,
					Type:    MessageType(msgType.(float64)),
					State:   uint64(msgState.(float64)),
					Body:    data,
				}
			}
		}
	}()
}

func (tzkt *TzKT) onReconnect() error {
	for i := range tzkt.subscriptions {
		if err := tzkt.s.Send(tzkt.subscriptions[i]); err != nil {
			return err
		}
	}
	return nil
}
