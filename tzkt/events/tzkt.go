package events

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/dipdup-net/go-lib/tzkt/events/signalr"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// TzKT - struct that used for connection to TzKT events server
type TzKT struct {
	s            *signalr.SignalR
	invokationID int

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
		msgs:          make(chan Message, 1024),
		stop:          make(chan struct{}, 1),
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

// IsConnected - reports whether the connection to TzKT events is established
func (tzkt *TzKT) IsConnected() bool {
	return tzkt.s != nil && tzkt.s.IsConnected()
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
	return tzkt.subscribe(MethodOperations, args)
}

// SubscribeToBigMaps - subscribe to bigmaps channel. Sends bigmap updates.
func (tzkt *TzKT) SubscribeToBigMaps(ptr *int64, contract, path string, tags ...string) error {
	args := make(map[string]interface{})
	if len(tags) > 0 {
		args["tags"] = tags
	}
	if contract != "" {
		args["contract"] = contract
	}
	if path != "" {
		args["path"] = path
	}
	if ptr != nil {
		args["ptr"] = *ptr
	}
	return tzkt.subscribe(MethodBigMap, args)
}

// SubscribeToAccounts - subscribe to accounts channel. Sends touched accounts (affected by any operation in any way)..
func (tzkt *TzKT) SubscribeToAccounts(addresses ...string) error {
	args := make(map[string]interface{})
	if len(addresses) > 0 {
		args["addresses"] = addresses
	}
	return tzkt.subscribe(MethodAccounts, args)
}

// SubscribeToTokenTransfers - subscribe to transfers channel. Sends token transfers.
func (tzkt *TzKT) SubscribeToTokenTransfers(account, contract, tokenID string) error {
	args := make(map[string]interface{})
	if account != "" {
		args["account"] = account
	}
	if contract != "" {
		args["contract"] = contract
	}
	if tokenID != "" {
		args["tokenID"] = tokenID
	}
	return tzkt.subscribe(MethodTokenTransfers, args)
}

func (tzkt *TzKT) subscribe(channel string, args ...interface{}) error {
	tzkt.invokationID += 1
	msg := signalr.NewInvocation(fmt.Sprintf("%d", tzkt.invokationID), channel, args...)
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
				switch typ := msg.(type) {
				case signalr.Invocation:
					if len(typ.Arguments) == 0 {
						log.Warn().Msgf("empty arguments of invocation: %v", typ)
						continue
					}

					var packet Packet
					if err := json.Unmarshal(typ.Arguments[0], &packet); err != nil {
						log.Err(err).Msg("invalid invocation argument")
						continue
					}

					message := Message{
						Channel: typ.Target,
						Type:    packet.Type,
						State:   packet.State,
					}

					if packet.Data != nil {
						data, err := parseData(typ.Target, packet.Data)
						if err != nil {
							log.Err(err).Msg("error during parsing data")
							continue
						}
						message.Body = data
					}

					tzkt.msgs <- message
				case signalr.Completion:
					for i := range tzkt.subscriptions {
						if tzkt.subscriptions[i].ID != typ.ID {
							continue
						}
						tzkt.msgs <- Message{
							Channel: tzkt.subscriptions[i].Target,
							Type:    MessageTypeSubscribed,
							State:   typ.Result,
						}
						break
					}
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

func parseData(channel string, data []byte) (interface{}, error) {
	switch channel {
	case ChannelAccounts:
		var acc []Account
		err := json.Unmarshal(data, &acc)
		return acc, err
	case ChannelBigMap:
		var updates []BigMapUpdate
		err := json.Unmarshal(data, &updates)
		return updates, err
	case ChannelBlocks:
		var block []Block
		err := json.Unmarshal(data, &block)
		return block, err
	case ChannelHead:
		var head Head
		err := json.Unmarshal(data, &head)
		return head, err
	case ChannelOperations:
		return parseOperations(data)
	case ChannelTransfers:
		var transfer []Transfer
		err := json.Unmarshal(data, &transfer)
		return transfer, err
	default:
		return nil, errors.Errorf("unknown channel: %s", channel)
	}
}

func parseOperations(data []byte) (interface{}, error) {
	var operations []Operation
	if err := json.Unmarshal(data, &operations); err != nil {
		return nil, err
	}
	result := make([]interface{}, 0)
	for i := range operations {
		switch operations[i].Type {
		case KindDelegation:
			result = append(result, &Delegation{})
		case KindOrigination:
			result = append(result, &Origination{})
		case KindReveal:
			result = append(result, &Reveal{})
		case KindTransaction:
			result = append(result, &Transaction{})
		default:
			result = append(result, make(map[string]interface{}))
		}
	}

	err := json.Unmarshal(data, &result)
	return result, err
}
