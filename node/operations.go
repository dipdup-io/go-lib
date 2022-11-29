package node

import (
	"bytes"
	stdJSON "encoding/json"

	"github.com/pkg/errors"
)

// OperationConstraint -
type OperationConstraint interface {
	AccountActivation | Ballot | Delegation | DoubleBakingEvidence |
		DoubleEndorsementEvidence | Endorsement | EndorsementWithSlot |
		Origination | Proposal | Reveal | SeedNonceRevelation | Transaction |
		RegisterGlobalConstant | DoublePreendorsementEvidence | SetDepositsLimit |
		Preendorsement | Event | VdfRevelation | TxRollupCommit | TxRollupOrigination |
		TxRollupDispatchTickets | TxRollupFinalizeCommitment | TxRollupRejection |
		TxRollupRemoveCommitment | TxRollupSubmitBatch | UpdateConsensusKey |
		DrainDelegate
}

// OperationGroup -
type OperationGroup struct {
	Protocol  string      `json:"protocol"`
	ChainID   string      `json:"chain_id"`
	Hash      string      `json:"hash"`
	Branch    string      `json:"branch"`
	Signature string      `json:"signature"`
	Contents  []Operation `json:"contents"`
}

// Operation -
type Operation struct {
	Kind string      `json:"kind"`
	Body interface{} `json:"-"`
}

// UnmarshalJSON -
func (op *Operation) UnmarshalJSON(data []byte) error {
	type buf Operation
	if err := json.Unmarshal(data, (*buf)(op)); err != nil {
		return err
	}

	var err error
	switch op.Kind {
	case KindActivation:
		err = parseOperation[AccountActivation](data, op)
	case KindBallot:
		err = parseOperation[Ballot](data, op)
	case KindDelegation:
		err = parseOperation[Delegation](data, op)
	case KindDoubleBaking:
		err = parseOperation[DoubleBakingEvidence](data, op)
	case KindDoubleEndorsing:
		err = parseOperation[DoubleEndorsementEvidence](data, op)
	case KindEndorsement:
		err = parseOperation[Endorsement](data, op)
	case KindEndorsementWithSlot:
		err = parseOperation[EndorsementWithSlot](data, op)
	case KindOrigination:
		err = parseOperation[Origination](data, op)
	case KindProposal:
		err = parseOperation[Proposal](data, op)
	case KindReveal:
		err = parseOperation[Reveal](data, op)
	case KindNonceRevelation:
		err = parseOperation[SeedNonceRevelation](data, op)
	case KindTransaction:
		err = parseOperation[Transaction](data, op)
	case KindRegisterGlobalConstant:
		err = parseOperation[RegisterGlobalConstant](data, op)
	case KindDoublePreendorsement:
		err = parseOperation[DoublePreendorsementEvidence](data, op)
	case KindSetDepositsLimit:
		err = parseOperation[SetDepositsLimit](data, op)
	case KindPreendorsement:
		err = parseOperation[Preendorsement](data, op)
	case KindEvent:
		err = parseOperation[Event](data, op)
	case KindVdfRevelation:
		err = parseOperation[VdfRevelation](data, op)
	case KindTxRollupOrigination:
		err = parseOperation[TxRollupOrigination](data, op)
	case KindTxRollupCommit:
		err = parseOperation[TxRollupCommit](data, op)
	case KindTxRollupDispatchTickets:
		err = parseOperation[TxRollupDispatchTickets](data, op)
	case KindTxRollupFinalizeCommitment:
		err = parseOperation[TxRollupFinalizeCommitment](data, op)
	case KindTxRollupRejection:
		err = parseOperation[TxRollupRejection](data, op)
	case KindTxRollupRemoveCommitment:
		err = parseOperation[TxRollupRemoveCommitment](data, op)
	case KindTxRollupSubmitBatch:
		err = parseOperation[TxRollupSubmitBatch](data, op)
	case KindUpdateConsensusKey:
		err = parseOperation[UpdateConsensusKey](data, op)
	case KindDrainDelegate:
		err = parseOperation[DrainDelegate](data, op)

	}
	return err
}

func parseOperation[M OperationConstraint](data []byte, operation *Operation) error {
	var model M
	if err := json.Unmarshal(data, &model); err != nil {
		return err
	}
	operation.Body = model
	return nil
}

// NewTypedOperation -
func NewTypedOperation[M OperationConstraint](operation Operation) (M, error) {
	if operation.Body == nil {
		var m M
		return m, errors.New("nil operation body")
	}
	model, ok := operation.Body.(M)
	if !ok {
		var m M
		return m, errors.Errorf("invalid body type: %T", operation.Body)
	}
	return model, nil
}

// AccountActivation -
type AccountActivation struct {
	Pkh      string                      `json:"pkh"`
	Secret   string                      `json:"secret"`
	Metadata *OnlyBalanceUpdatesMetadata `json:"metadata,omitempty"`
}

// Ballot -
type Ballot struct {
	Source   string `json:"source"`
	Period   uint64 `json:"period"`
	Proposal string `json:"proposal"`
	Ballot   string `json:"ballot"`
}

// Endorsement -
type Endorsement struct {
	Level    uint64               `json:"level"`
	Metadata *EndorsementMetadata `json:"metadata,omitempty"`
}

// EndorsementWithSlot -
type EndorsementWithSlot struct {
	Endorsement EndorsementWithSlotEntity `json:"endorsement"`
	Slot        uint64                    `json:"slot"`
	Metadata    *EndorsementMetadata      `json:"metadata,omitempty"`
}

// Preendorsement -
type Preendorsement struct {
	Slot             uint64 `json:"slot"`
	Level            uint64 `json:"level"`
	Round            int64  `json:"round"`
	BlockPayloadHash string `json:"block_payload_hash"`
	Metadata         struct {
		BalanceUpdates      []interface{} `json:"balance_updates"`
		Delegate            string        `json:"delegate"`
		PreendorsementPower int           `json:"preendorsement_power"`
	} `json:"metadata"`
}

// Delegation -
type Delegation struct {
	Source       string                    `json:"source"`
	Fee          string                    `json:"fee"`
	Counter      string                    `json:"counter"`
	GasLimit     string                    `json:"gas_limit"`
	StorageLimit string                    `json:"storage_limit"`
	Delegate     string                    `json:"delegate,omitempty"`
	Metadata     *ManagerOperationMetadata `json:"metadata,omitempty"`
}

// DoubleBakingEvidence -
type DoubleBakingEvidence struct {
	Bh1      *Header                     `json:"bh1,omitempty"`
	Bh2      *Header                     `json:"bh2,omitempty"`
	Metadata *OnlyBalanceUpdatesMetadata `json:"metadata,omitempty"`
}

// DoubleEndorsementEvidence -
type DoubleEndorsementEvidence struct {
	Op1      *InlinedEndorsement         `json:"op1"`
	Op2      *InlinedEndorsement         `json:"op2"`
	Metadata *OnlyBalanceUpdatesMetadata `json:"metadata,omitempty"`
}

// DoublePreendorsementEvidence -
type DoublePreendorsementEvidence struct {
	Op1      *InlinedEndorsement         `json:"op1"`
	Op2      *InlinedEndorsement         `json:"op2"`
	Metadata *OnlyBalanceUpdatesMetadata `json:"metadata,omitempty"`
}

// Origination -
type Origination struct {
	Source        string                    `json:"source"`
	Fee           string                    `json:"fee"`
	Counter       string                    `json:"counter"`
	GasLimit      string                    `json:"gas_limit"`
	StorageLimit  string                    `json:"storage_limit"`
	Balance       string                    `json:"balance"`
	Delegate      string                    `json:"delegate,omitempty"`
	Script        stdJSON.RawMessage        `json:"script"`
	ManagerPubkey string                    `json:"managerPubkey,omitempty"`
	Metadata      *ManagerOperationMetadata `json:"metadata,omitempty"`
}

// Proposal -
type Proposal struct {
	Source    string   `json:"source"`
	Period    uint64   `json:"period"`
	Proposals []string `json:"proposals"`
}

// Reveal -
type Reveal struct {
	Source       string                    `json:"source"`
	Fee          string                    `json:"fee"`
	Counter      string                    `json:"counter"`
	GasLimit     string                    `json:"gas_limit"`
	StorageLimit string                    `json:"storage_limit"`
	PublicKey    string                    `json:"public_key"`
	Metadata     *ManagerOperationMetadata `json:"metadata,omitempty"`
}

// RegisterGlobalConstant -
type RegisterGlobalConstant struct {
	Source       string             `json:"source"`
	Fee          string             `json:"fee"`
	Counter      string             `json:"counter"`
	GasLimit     string             `json:"gas_limit"`
	StorageLimit string             `json:"storage_limit"`
	Value        stdJSON.RawMessage `json:"value"`
}

// SeedNonceRevelation -
type SeedNonceRevelation struct {
	Level    uint64                      `json:"level"`
	Nonce    string                      `json:"nonce"`
	Metadata *OnlyBalanceUpdatesMetadata `json:"metadata,omitempty"`
}

// SetDepositsLimit -
type SetDepositsLimit struct {
	Kind         string                    `json:"kind"`
	Source       string                    `json:"source"`
	Fee          string                    `json:"fee"`
	Counter      string                    `json:"counter"`
	GasLimit     string                    `json:"gas_limit"`
	StorageLimit string                    `json:"storage_limit"`
	Limit        *string                   `json:"limit,omitempty"`
	Metadata     *ManagerOperationMetadata `json:"metadata"`
}

// Event -
type Event struct {
	Kind    string             `json:"kind"`
	Source  string             `json:"source"`
	Nonce   int                `json:"nonce"`
	Type    stdJSON.RawMessage `json:"type"`
	Tag     string             `json:"tag"`
	Payload stdJSON.RawMessage `json:"payload"`
	Result  OperationResult    `json:"result"`
}

// VdfRevelation -
type VdfRevelation struct {
	Kind     string                    `json:"kind"`
	Solution []string                  `json:"solution"`
	Metadata *ManagerOperationMetadata `json:"metadata,omitempty"`
}

// Transaction -
type Transaction struct {
	Source       string                    `json:"source"`
	Fee          string                    `json:"fee"`
	Counter      string                    `json:"counter"`
	GasLimit     string                    `json:"gas_limit"`
	StorageLimit string                    `json:"storage_limit"`
	Amount       string                    `json:"amount"`
	Destination  string                    `json:"destination"`
	Parameters   *Parameters               `json:"parameters,omitempty"`
	Metadata     *ManagerOperationMetadata `json:"metadata,omitempty"`
}

// Parameters -
type Parameters struct {
	Entrypoint string              `json:"entrypoint"`
	Value      *stdJSON.RawMessage `json:"value"`
}

// ManagerOperationMetadata -
type ManagerOperationMetadata struct {
	BalanceUpdates           []BalanceUpdate `json:"balance_updates"`
	OperationResult          OperationResult `json:"operation_result"`
	InternalOperationResults []Operation     `json:"internal_operation_results,omitempty"`
}

// OnlyBalanceUpdatesMetadata -
type OnlyBalanceUpdatesMetadata struct {
	BalanceUpdates []BalanceUpdate `json:"balance_updates"`
}

// EndorsementMetadata -
type EndorsementMetadata struct {
	BalanceUpdates []BalanceUpdate `json:"balance_updates"`
	Delegate       string          `json:"delegate"`
	Slots          []int           `json:"slots"`
}

// OperationResult -
type OperationResult struct {
	Status                       string              `json:"status"`
	Storage                      *stdJSON.RawMessage `json:"storage,omitempty"`
	BigMapDiff                   []BigMapDiff        `json:"big_map_diff,omitempty"`
	LazyStorageDiff              []LazyStorageDiff   `json:"lazy_storage_diff,omitempty"`
	BalanceUpdates               []BalanceUpdate     `json:"balance_updates,omitempty"`
	OriginatedContracts          []string            `json:"originated_contracts,omitempty"`
	ConsumedGas                  string              `json:"consumed_gas,omitempty"`
	ConsumedMilligas             string              `json:"consumed_milligas,omitempty"`
	StorageSize                  string              `json:"storage_size,omitempty"`
	OriginatedRollup             string              `json:"originated_rollup"`
	PaidStorageSizeDiff          string              `json:"paid_storage_size_diff,omitempty"`
	AllocatedDestinationContract bool                `json:"allocated_destination_contract,omitempty"`
	Level                        *uint64             `json:"level,omitempty"`
	Errors                       []ResultError       `json:"errors,omitempty"`
}

// ResultError -
type ResultError struct {
	Kind           string              `json:"kind"`
	ID             string              `json:"id,omitempty"`
	With           *stdJSON.RawMessage `json:"with,omitempty"`
	Msg            string              `json:"msg,omitempty"`
	Location       int64               `json:"location,omitempty"`
	ContractHandle string              `json:"contract_handle,omitempty"`
	ContractCode   *stdJSON.RawMessage `json:"contract_code,omitempty"`
}

// BigMapDiff is an element of BigMapDiffs
type BigMapDiff struct {
	Action            string              `json:"action,omitempty"`
	BigMap            string              `json:"big_map,omitempty"`
	KeyHash           string              `json:"key_hash,omitempty"`
	Key               *stdJSON.RawMessage `json:"key,omitempty"`
	Value             *stdJSON.RawMessage `json:"value,omitempty"`
	SourceBigMap      string              `json:"source_big_map,omitempty"`
	DestinationBigMap string              `json:"destination_big_map,omitempty"`
	KeyType           *stdJSON.RawMessage `json:"key_type,omitempty"`
	ValueType         *stdJSON.RawMessage `json:"value_type,omitempty"`
}

// LazyStorageDiff -
type LazyStorageDiff struct {
	Kind string                      `json:"kind"`
	ID   string                      `json:"id"`
	Diff LazyStorageDiffBigMapEntity `json:"diff"`
}

// LazyStorageDiffBigMapEntity -
type LazyStorageDiffBigMapEntity struct {
	Action    string                  `json:"action"`
	Updates   []LazyStorageDiffUpdate `json:"updates,omitempty"`
	Source    string                  `json:"source,omitempty"`
	KeyType   *stdJSON.RawMessage     `json:"key_type,omitempty"`
	ValueType *stdJSON.RawMessage     `json:"value_type,omitempty"`
}

// LazyStorageDiffUpdate -
type LazyStorageDiffUpdate struct {
	Action string              `json:"action,omitempty"`
	Key    *stdJSON.RawMessage `json:"key,omitempty"`
	Value  *stdJSON.RawMessage `json:"value,omitempty"`
}

// LazyStorageDiffSaplingEntity -
type LazyStorageDiffSaplingEntity struct {
	Action   string                  `json:"action"`
	Updates  []LazyStorageDiffUpdate `json:"updates,omitempty"`
	Source   string                  `json:"source,omitempty"`
	MemoSize uint64                  `json:"memo_size,omitempty"`
}

// LazyStorageDiffUpdatesSaplingState -
type LazyStorageDiffUpdatesSaplingState struct {
	CommitmentsAndCiphertexts CommitmentsAndCiphertexts `json:"commitments_and_ciphertexts,omitempty"`
	Nullifiers                []string                  `json:"nullifiers,omitempty"`
}

// CommitmentsAndCiphertexts-
type CommitmentsAndCiphertexts struct {
	Commitments string
	Ciphertexts SaplingTransactionCiphertext
}

// UnmarshalJSON -
func (cc *CommitmentsAndCiphertexts) UnmarshalJSON(data []byte) error {
	raw := []interface{}{
		cc.Commitments, cc.Ciphertexts,
	}
	return json.Unmarshal(data, &raw)
}

// MarshalJSON -
func (cc CommitmentsAndCiphertexts) MarshalJSON() ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.WriteByte('[')

	buf.WriteByte('"')
	buf.WriteString(cc.Commitments)
	buf.WriteByte('"')

	buf.WriteByte(',')

	ciphertexts, err := json.Marshal(cc.Ciphertexts)
	if err != nil {
		return nil, err
	}
	buf.Write(ciphertexts)

	buf.WriteByte(']')
	return buf.Bytes(), nil
}

// SaplingTransactionCiphertext -
type SaplingTransactionCiphertext struct {
	CV         string `json:"cv"`
	EPK        string `json:"epk"`
	PayloadEnc string `json:"payload_enc"`
	NonceEnc   string `json:"nonce_enc"`
	PayloadOut string `json:"payload_out"`
	NonceOut   string `json:"nonce_out"`
}

// InlinedEndorsement -
type InlinedEndorsement struct {
	Branch     string                        `json:"branch"`
	Operations *InlinedEndorsementOperations `json:"operations,omitempty"`
	Signature  string                        `json:"signature"`
}

// InlinedEndorsementOperations -
type InlinedEndorsementOperations struct {
	Kind  string `json:"kind"`
	Level int    `json:"level"`
}

// InlinedPreendorsement -
type InlinedPreendorsement struct {
	Branch     string                           `json:"branch"`
	Operations *InlinedPreendorsementOperations `json:"operations,omitempty"`
	Signature  string                           `json:"signature"`
}

// InlinedPreendorsementOperations -
type InlinedPreendorsementOperations struct {
	Kind             string `json:"kind"`
	Slot             uint64 `json:"slot"`
	Level            uint64 `json:"level"`
	Round            int64  `json:"round"`
	BlockPayloadHash string `json:"block_payload_hash"`
}

// EndorsementWithSlotEntity -
type EndorsementWithSlotEntity struct {
	Branch    string               `json:"branch"`
	Operation EndorsementOperation `json:"operations"`
	Signature string               `json:"signature"`
}

// EndorsementOperation -
type EndorsementOperation struct {
	Level uint64 `json:"level"`
}

// TxRollupOrigination -
type TxRollupOrigination struct {
	Kind                string                    `json:"kind"`
	Source              string                    `json:"source"`
	Fee                 string                    `json:"fee"`
	Counter             string                    `json:"counter"`
	GasLimit            string                    `json:"gas_limit"`
	StorageLimit        string                    `json:"storage_limit"`
	TxRollupOrigination any                       `json:"tx_rollup_origination"`
	Metadata            *ManagerOperationMetadata `json:"metadata,omitempty"`
}

// TxRollupCommitment -
type TxRollupCommitment struct {
	Level           uint64   `json:"level"`
	Messages        []string `json:"messages"`
	Predecessor     string   `json:"predecessor"`
	InboxMerkleRoot string   `json:"inbox_merkle_root"`
}

// TxRollupCommit -
type TxRollupCommit struct {
	Kind         string                    `json:"kind"`
	Source       string                    `json:"source"`
	Fee          string                    `json:"fee"`
	Counter      string                    `json:"counter"`
	GasLimit     string                    `json:"gas_limit"`
	StorageLimit string                    `json:"storage_limit"`
	Rollup       string                    `json:"rollup"`
	Commitment   TxRollupCommitment        `json:"commitment"`
	Metadata     *ManagerOperationMetadata `json:"metadata,omitempty"`
}

// TxRollupDispatchTickets -
type TxRollupDispatchTickets struct {
	Kind              string                    `json:"kind"`
	Source            string                    `json:"source"`
	Fee               string                    `json:"fee"`
	Counter           string                    `json:"counter"`
	GasLimit          string                    `json:"gas_limit"`
	StorageLimit      string                    `json:"storage_limit"`
	TxRollup          string                    `json:"tx_rollup"`
	Level             int                       `json:"level"`
	ContextHash       string                    `json:"context_hash"`
	MessageIndex      int                       `json:"message_index"`
	MessageResultPath []string                  `json:"message_result_path"`
	TicketsInfo       []TicketsInfo             `json:"tickets_info"`
	Metadata          *ManagerOperationMetadata `json:"metadata,omitempty"`
}

// TicketsInfo -
type TicketsInfo struct {
	Contents stdJSON.RawMessage `json:"contents"`
	Ty       stdJSON.RawMessage `json:"ty"`
	Ticketer string             `json:"ticketer"`
	Amount   string             `json:"amount"`
	Claimer  string             `json:"claimer"`
}

// TxRollupFinalizeCommitment -
type TxRollupFinalizeCommitment struct {
	Kind         string                    `json:"kind"`
	Source       string                    `json:"source"`
	Fee          string                    `json:"fee"`
	Counter      string                    `json:"counter"`
	GasLimit     string                    `json:"gas_limit"`
	StorageLimit string                    `json:"storage_limit"`
	Rollup       string                    `json:"rollup"`
	Metadata     *ManagerOperationMetadata `json:"metadata,omitempty"`
}

// TxRollupRejection -
type TxRollupRejection struct {
	Kind         string `json:"kind"`
	Source       string `json:"source"`
	Fee          string `json:"fee"`
	Counter      string `json:"counter"`
	GasLimit     string `json:"gas_limit"`
	StorageLimit string `json:"storage_limit"`
	Rollup       string `json:"rollup"`
	Level        int    `json:"level"`
	Message      struct {
		Batch string `json:"batch"`
	} `json:"message"`
	MessagePosition       string   `json:"message_position"`
	MessagePath           []string `json:"message_path"`
	MessageResultHash     string   `json:"message_result_hash"`
	MessageResultPath     []string `json:"message_result_path"`
	PreviousMessageResult struct {
		ContextHash      string `json:"context_hash"`
		WithdrawListHash string `json:"withdraw_list_hash"`
	} `json:"previous_message_result"`
	PreviousMessageResultPath []string `json:"previous_message_result_path"`
	Proof                     struct {
		Version uint64             `json:"version"`
		Before  stdJSON.RawMessage `json:"before"`
		After   stdJSON.RawMessage `json:"after"`
		State   stdJSON.RawMessage `json:"proof"`
	} `json:"proof"`
	Metadata struct {
		BalanceUpdates []struct {
			Kind     string `json:"kind"`
			Contract string `json:"contract,omitempty"`
			Change   string `json:"change"`
			Origin   string `json:"origin"`
			Category string `json:"category,omitempty"`
		} `json:"balance_updates"`
		OperationResult struct {
			Status         string `json:"status"`
			BalanceUpdates []struct {
				Kind     string `json:"kind"`
				Category string `json:"category,omitempty"`
				Contract string `json:"contract,omitempty"`
				BondID   struct {
					TxRollup string `json:"tx_rollup"`
				} `json:"bond_id,omitempty"`
				Change string `json:"change"`
				Origin string `json:"origin"`
			} `json:"balance_updates"`
			ConsumedGas      string `json:"consumed_gas"`
			ConsumedMilligas string `json:"consumed_milligas"`
		} `json:"operation_result"`
	} `json:"metadata"`
}

// TxRollupRemoveCommitment -
type TxRollupRemoveCommitment struct {
	Kind         string                    `json:"kind"`
	Source       string                    `json:"source"`
	Fee          string                    `json:"fee"`
	Counter      string                    `json:"counter"`
	GasLimit     string                    `json:"gas_limit"`
	StorageLimit string                    `json:"storage_limit"`
	Rollup       string                    `json:"rollup"`
	Metadata     *ManagerOperationMetadata `json:"metadata,omitempty"`
}

// TxRollupSubmitBatch -
type TxRollupSubmitBatch struct {
	Kind         string                    `json:"kind"`
	Source       string                    `json:"source"`
	Fee          string                    `json:"fee"`
	Counter      string                    `json:"counter"`
	GasLimit     string                    `json:"gas_limit"`
	StorageLimit string                    `json:"storage_limit"`
	Rollup       string                    `json:"rollup"`
	Content      string                    `json:"content"`
	Metadata     *ManagerOperationMetadata `json:"metadata,omitempty"`
}

// UpdateConsensusKey -
type UpdateConsensusKey struct {
	Kind         string                    `json:"kind"`
	Source       string                    `json:"source"`
	Fee          string                    `json:"fee"`
	Counter      string                    `json:"counter"`
	GasLimit     string                    `json:"gas_limit"`
	StorageLimit string                    `json:"storage_limit"`
	Pk           string                    `json:"pk"`
	Metadata     *ManagerOperationMetadata `json:"metadata,omitempty"`
}

// DrainDelegate -
type DrainDelegate struct {
	Kind         string                    `json:"kind"`
	ConsensusKey string                    `json:"consensus_key"`
	Delegate     string                    `json:"delegate"`
	Destination  string                    `json:"destination"`
	Metadata     *ManagerOperationMetadata `json:"metadata,omitempty"`
}
