package events

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/dipdup-net/go-lib/tzkt/events/signalr"

	tzktData "github.com/dipdup-net/go-lib/tzkt/data"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// TzKT - struct that used for connection to TzKT events server
type TzKT struct {
	s            *signalr.SignalR
	invokationID int

	subscriptions []signalr.Invocation

	log zerolog.Logger

	msgs chan Message
	wg   sync.WaitGroup
}

// NewTzKT - constructor of `TzKT`. `url` is TzKT events base URL. If it's empty https://api.tzkt.io/v1/ws is set.
func NewTzKT(url string) *TzKT {
	if url == "" {
		url = tzktData.BaseEventsURL
	}
	return &TzKT{
		s:             signalr.NewSignalR(url),
		msgs:          make(chan Message, 1024),
		subscriptions: make([]signalr.Invocation, 0),
		log:           log.Logger,
	}
}

// SetLogger -
func (tzkt *TzKT) SetLogger(logger zerolog.Logger) {
	tzkt.log = logger
	tzkt.s.SetLogger(logger)
}

// Connect - connect to events SignalR server
func (tzkt *TzKT) Connect(ctx context.Context) error {
	if err := tzkt.s.Connect(ctx, signalr.Version1); err != nil {
		return err
	}
	tzkt.s.SetOnReconnect(tzkt.onReconnect)
	tzkt.listen(ctx)
	return nil
}

// Close - closing all connections
func (tzkt *TzKT) Close() error {
	tzkt.wg.Wait()

	if err := tzkt.s.Close(); err != nil {
		return err
	}
	close(tzkt.msgs)
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

// SubscribeToTokenBalances - sends token balances when they are updated.
func (tzkt *TzKT) SubscribeToTokenBalances(account, contract, tokenID string) error {
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
	return tzkt.subscribe(MethodTokenBalances, args)
}

// SubscribeToCycles - notifies of the start of a new cycle with a specified delay.
// delayBlocks is the number of blocks (2 by default) to delay a new cycle notification. It should be >= 2 (to not worry abour reorgs) and < cycle size
func (tzkt *TzKT) SubscribeToCycles(delayBlocks uint64) error {
	if delayBlocks < 2 {
		return errors.Errorf("delayBocks should be >= 2: %d", delayBlocks)
	}
	args := map[string]any{
		"delayBocks": delayBlocks,
	}
	return tzkt.subscribe(MethodCycles, args)
}

func (tzkt *TzKT) subscribe(channel string, args ...interface{}) error {
	tzkt.invokationID += 1
	msg := signalr.NewInvocation(fmt.Sprintf("%d", tzkt.invokationID), channel, args...)
	tzkt.subscriptions = append(tzkt.subscriptions, msg)
	return tzkt.s.Send(msg)
}

func (tzkt *TzKT) listen(ctx context.Context) {
	tzkt.wg.Add(1)

	go func() {
		defer tzkt.wg.Done()

		for {
			select {
			case <-ctx.Done():
				tzkt.log.Debug().Msg("listening was stopped")
				return
			case msg := <-tzkt.s.Messages():
				switch typ := msg.(type) {
				case signalr.Invocation:
					if len(typ.Arguments) == 0 {
						tzkt.log.Warn().Msgf("empty arguments of invocation: %v", typ)
						continue
					}

					var packet Packet
					if err := json.Unmarshal(typ.Arguments[0], &packet); err != nil {
						tzkt.log.Err(err).Msg("invalid invocation argument")
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
							tzkt.log.Err(err).Msg("error during parsing data")
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

func parseData(channel string, data []byte) (any, error) {
	switch channel {
	case ChannelAccounts:
		var acc []tzktData.Account
		err := json.Unmarshal(data, &acc)
		return acc, err
	case ChannelBigMap:
		var updates []tzktData.BigMapUpdate
		err := json.Unmarshal(data, &updates)
		return updates, err
	case ChannelBlocks:
		var block []tzktData.Block
		err := json.Unmarshal(data, &block)
		return block, err
	case ChannelHead:
		var head tzktData.Head
		err := json.Unmarshal(data, &head)
		return head, err
	case ChannelOperations:
		return parseOperations(data)
	case ChannelTransfers:
		var transfer []tzktData.Transfer
		err := json.Unmarshal(data, &transfer)
		return transfer, err
	case ChannelTokenBalances:
		var balances []tzktData.TokenBalance
		err := json.Unmarshal(data, &balances)
		return balances, err
	case ChannelCycles:
		var cycle tzktData.Cycle
		err := json.Unmarshal(data, &cycle)
		return cycle, err
	default:
		return nil, errors.Errorf("unknown channel: %s", channel)
	}
}

func parseOperations(data []byte) (any, error) {
	var operations []tzktData.Operation
	if err := json.Unmarshal(data, &operations); err != nil {
		return nil, err
	}
	result := make([]any, 0)
	for i := range operations {
		switch operations[i].Type {
		case tzktData.KindDelegation:
			result = append(result, &tzktData.Delegation{})
		case tzktData.KindOrigination:
			result = append(result, &tzktData.Origination{})
		case tzktData.KindReveal:
			result = append(result, &tzktData.Reveal{})
		case tzktData.KindTransaction:
			result = append(result, &tzktData.Transaction{})
		case tzktData.KindMigration:
			result = append(result, &tzktData.Migration{})
		case tzktData.KindActivation:
			result = append(result, &tzktData.Activation{})
		case tzktData.KindBallot:
			result = append(result, &tzktData.Ballot{})
		case tzktData.KindDoubleBaking:
			result = append(result, &tzktData.DoubleBaking{})
		case tzktData.KindDoubleConsensus:
			result = append(result, &tzktData.DoubleConsensus{})
		case tzktData.KindAttestation:
			result = append(result, &tzktData.Attestation{})
		case tzktData.KindNonceRevelation:
			result = append(result, &tzktData.NonceRevelation{})
		case tzktData.KindProposal:
			result = append(result, &tzktData.Proposal{})
		case tzktData.KindPreattestations:
			result = append(result, &tzktData.Preattestation{})
		case tzktData.KindRegisterGlobalConstant:
			result = append(result, &tzktData.RegisterConstant{})
		case tzktData.KindSetDepositsLimit:
			result = append(result, &tzktData.SetDepositsLimit{})
		case tzktData.KindSetDelegateParameters:
			result = append(result, &tzktData.SetDelegateParameters{})
		case tzktData.KindRollupDispatchTickets:
			result = append(result, &tzktData.TxRollupDispatchTicket{})
		case tzktData.KindRollupFinalizeCommitment:
			result = append(result, &tzktData.TxRollupFinalizeCommitment{})
		case tzktData.KindRollupReturnBond:
			result = append(result, &tzktData.TxRollupReturnBond{})
		case tzktData.KindRollupSubmitBatch:
			result = append(result, &tzktData.TxRollupSubmitBatch{})
		case tzktData.KindTransferTicket:
			result = append(result, &tzktData.TransferTicket{})
		case tzktData.KindTxRollupCommit:
			result = append(result, &tzktData.TxRollupCommit{})
		case tzktData.KindTxRollupOrigination:
			result = append(result, &tzktData.TxRollupOrigination{})
		case tzktData.KindTxRollupRejection:
			result = append(result, &tzktData.TxRollupRejection{})
		case tzktData.KindTxRollupRemoveCommitment:
			result = append(result, &tzktData.TxRollupRemoveCommitment{})
		case tzktData.KindRevelationPenalty:
			result = append(result, &tzktData.RevelationPenalty{})
		case tzktData.KindAttestationReward:
			result = append(result, &tzktData.AttestationReward{})
		case tzktData.KindBaking:
			result = append(result, &tzktData.Baking{})
		case tzktData.KindIncreasePaidStorage:
			result = append(result, &tzktData.IncreasePaidStorage{})
		case tzktData.KindVdfRevelation:
			result = append(result, &tzktData.VdfRevelation{})
		case tzktData.KindUpdateSecondaryKey:
			result = append(result, &tzktData.UpdateSecondaryKey{})
		case tzktData.KindDrainDelegate:
			result = append(result, &tzktData.DrainDelegate{})
		case tzktData.KindSrAddMessages:
			result = append(result, &tzktData.SmartRollupAddMessage{})
		case tzktData.KindSrCement:
			result = append(result, &tzktData.SmartRollupCement{})
		case tzktData.KindSrExecute:
			result = append(result, &tzktData.SmartRollupExecute{})
		case tzktData.KindSrOriginate:
			result = append(result, &tzktData.SmartRollupOriginate{})
		case tzktData.KindSrPublish:
			result = append(result, &tzktData.SmartRollupPublish{})
		case tzktData.KindSrRecoverBond:
			result = append(result, &tzktData.SmartRollupRecoverBond{})
		case tzktData.KindSrRefute:
			result = append(result, &tzktData.SmartRollupRefute{})
		case tzktData.KindDalPublishCommitment:
			result = append(result, &tzktData.DalPublishCommitment{})
		case tzktData.KindDalAttestationReward:
			result = append(result, &tzktData.DalAttestationReward{})
		case tzktData.KindStaking:
			result = append(result, &tzktData.Staking{})
		default:
			result = append(result, make(map[string]interface{}))
		}
	}

	err := json.Unmarshal(data, &result)
	return result, err
}
