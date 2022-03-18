package node

import (
	"bytes"
	stdJSON "encoding/json"

	"github.com/pkg/errors"
)

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

	switch op.Kind {
	case KindActivation:
		var activation AccountActivation
		if err := json.Unmarshal(data, &activation); err != nil {
			return err
		}
		op.Body = activation
	case KindBallot:
		var ballot Ballot
		if err := json.Unmarshal(data, &ballot); err != nil {
			return err
		}
		op.Body = ballot
	case KindDelegation:
		var delegation Delegation
		if err := json.Unmarshal(data, &delegation); err != nil {
			return err
		}
		op.Body = delegation
	case KindDoubleBaking:
		var evidence DoubleBakingEvidence
		if err := json.Unmarshal(data, &evidence); err != nil {
			return err
		}
		op.Body = evidence
	case KindDoubleEndorsing:
		var evidence DoubleEndorsementEvidence
		if err := json.Unmarshal(data, &evidence); err != nil {
			return err
		}
		op.Body = evidence
	case KindEndorsement:
		var endorsement Endorsement
		if err := json.Unmarshal(data, &endorsement); err != nil {
			return err
		}
		op.Body = endorsement
	case KindEndorsementWithSlot:
		var endorsement EndorsementWithSlot
		if err := json.Unmarshal(data, &endorsement); err != nil {
			return err
		}
		op.Body = endorsement
	case KindOrigination:
		var origination Origination
		if err := json.Unmarshal(data, &origination); err != nil {
			return err
		}
		op.Body = origination
	case KindProposal:
		var proposal Proposal
		if err := json.Unmarshal(data, &proposal); err != nil {
			return err
		}
		op.Body = proposal
	case KindReveal:
		var reveal Reveal
		if err := json.Unmarshal(data, &reveal); err != nil {
			return err
		}
		op.Body = reveal
	case KindNonceRevelation:
		var seed SeedNonceRevelation
		if err := json.Unmarshal(data, &seed); err != nil {
			return err
		}
		op.Body = seed
	case KindTransaction:
		var transaction Transaction
		if err := json.Unmarshal(data, &transaction); err != nil {
			return err
		}
		op.Body = transaction
	case KindRegisterGlobalConstant:
		var register RegisterGlobalConstant
		if err := json.Unmarshal(data, &register); err != nil {
			return err
		}
		op.Body = register
	case KindDoublePreendorsement:
		var doublePreendorsement DoublePreendorsementEvidence
		if err := json.Unmarshal(data, &doublePreendorsement); err != nil {
			return err
		}
		op.Body = doublePreendorsement
	case KindSetDepositsLimit:
		var setDepositsLimit SetDepositsLimit
		if err := json.Unmarshal(data, &setDepositsLimit); err != nil {
			return err
		}
		op.Body = setDepositsLimit
	case KindPreendorsement:
		var preendorsement Preendorsement
		if err := json.Unmarshal(data, &preendorsement); err != nil {
			return err
		}
		op.Body = preendorsement
	}
	return nil
}

// AccountActivation -
func (op Operation) AccountActivation() (AccountActivation, error) {
	if op.Kind != KindActivation {
		return AccountActivation{}, errors.Errorf("invalid kind of operation: %s", op.Kind)
	}
	if op.Body == nil {
		return AccountActivation{}, errors.New("nil operation body")
	}
	activation, ok := op.Body.(AccountActivation)
	if !ok {
		return AccountActivation{}, errors.Errorf("invalid body type: %T", op.Body)
	}
	return activation, nil
}

// Ballot -
func (op Operation) Ballot() (Ballot, error) {
	if op.Kind != KindBallot {
		return Ballot{}, errors.Errorf("invalid kind of operation: %s", op.Kind)
	}
	if op.Body == nil {
		return Ballot{}, errors.New("nil operation body")
	}
	ballot, ok := op.Body.(Ballot)
	if !ok {
		return Ballot{}, errors.Errorf("invalid body type: %T", op.Body)
	}
	return ballot, nil
}

// Endorsement -
func (op Operation) Endorsement() (Endorsement, error) {
	if op.Kind != KindEndorsement {
		return Endorsement{}, errors.Errorf("invalid kind of operation: %s", op.Kind)
	}
	if op.Body == nil {
		return Endorsement{}, errors.New("nil operation body")
	}
	endorsement, ok := op.Body.(Endorsement)
	if !ok {
		return Endorsement{}, errors.Errorf("invalid body type: %T", op.Body)
	}
	return endorsement, nil
}

// Preendorsement -
func (op Operation) Preendorsement() (Preendorsement, error) {
	if op.Kind != KindPreendorsement {
		return Preendorsement{}, errors.Errorf("invalid kind of operation: %s", op.Kind)
	}
	if op.Body == nil {
		return Preendorsement{}, errors.New("nil operation body")
	}
	preendorsement, ok := op.Body.(Preendorsement)
	if !ok {
		return Preendorsement{}, errors.Errorf("invalid body type: %T", op.Body)
	}
	return preendorsement, nil
}

// EndorsementWithSlot -
func (op Operation) EndorsementWithSlot() (EndorsementWithSlot, error) {
	if op.Kind != KindEndorsementWithSlot {
		return EndorsementWithSlot{}, errors.Errorf("invalid kind of operation: %s", op.Kind)
	}
	if op.Body == nil {
		return EndorsementWithSlot{}, errors.New("nil operation body")
	}
	endorsement, ok := op.Body.(EndorsementWithSlot)
	if !ok {
		return EndorsementWithSlot{}, errors.Errorf("invalid body type: %T", op.Body)
	}
	return endorsement, nil
}

// Delegation -
func (op Operation) Delegation() (Delegation, error) {
	if op.Kind != KindDelegation {
		return Delegation{}, errors.Errorf("invalid kind of operation: %s", op.Kind)
	}
	if op.Body == nil {
		return Delegation{}, errors.New("nil operation body")
	}
	delegation, ok := op.Body.(Delegation)
	if !ok {
		return Delegation{}, errors.Errorf("invalid body type: %T", op.Body)
	}
	return delegation, nil
}

// DoubleBakingEvidence -
func (op Operation) DoubleBakingEvidence() (DoubleBakingEvidence, error) {
	if op.Kind != KindDelegation {
		return DoubleBakingEvidence{}, errors.Errorf("invalid kind of operation: %s", op.Kind)
	}
	if op.Body == nil {
		return DoubleBakingEvidence{}, errors.New("nil operation body")
	}
	evidence, ok := op.Body.(DoubleBakingEvidence)
	if !ok {
		return DoubleBakingEvidence{}, errors.Errorf("invalid body type: %T", op.Body)
	}
	return evidence, nil
}

// DoubleEndorsementEvidence -
func (op Operation) DoubleEndorsementEvidence() (DoubleEndorsementEvidence, error) {
	if op.Kind != KindDelegation {
		return DoubleEndorsementEvidence{}, errors.Errorf("invalid kind of operation: %s", op.Kind)
	}
	if op.Body == nil {
		return DoubleEndorsementEvidence{}, errors.New("nil operation body")
	}
	evidence, ok := op.Body.(DoubleEndorsementEvidence)
	if !ok {
		return DoubleEndorsementEvidence{}, errors.Errorf("invalid body type: %T", op.Body)
	}
	return evidence, nil
}

// Origination -
func (op Operation) Origination() (Origination, error) {
	if op.Kind != KindOrigination {
		return Origination{}, errors.Errorf("invalid kind of operation: %s", op.Kind)
	}
	if op.Body == nil {
		return Origination{}, errors.New("nil operation body")
	}
	origination, ok := op.Body.(Origination)
	if !ok {
		return Origination{}, errors.Errorf("invalid body type: %T", op.Body)
	}
	return origination, nil
}

// Proposal -
func (op Operation) Proposal() (Proposal, error) {
	if op.Kind != KindProposal {
		return Proposal{}, errors.Errorf("invalid kind of operation: %s", op.Kind)
	}
	if op.Body == nil {
		return Proposal{}, errors.New("nil operation body")
	}
	proposal, ok := op.Body.(Proposal)
	if !ok {
		return Proposal{}, errors.Errorf("invalid body type: %T", op.Body)
	}
	return proposal, nil
}

// Reveal -
func (op Operation) Reveal() (Reveal, error) {
	if op.Kind != KindReveal {
		return Reveal{}, errors.Errorf("invalid kind of operation: %s", op.Kind)
	}
	if op.Body == nil {
		return Reveal{}, errors.New("nil operation body")
	}
	reveal, ok := op.Body.(Reveal)
	if !ok {
		return Reveal{}, errors.Errorf("invalid body type: %T", op.Body)
	}
	return reveal, nil
}

// Reveal -
func (op Operation) RegisterGlobalConstant() (RegisterGlobalConstant, error) {
	if op.Kind != KindRegisterGlobalConstant {
		return RegisterGlobalConstant{}, errors.Errorf("invalid kind of operation: %s", op.Kind)
	}
	if op.Body == nil {
		return RegisterGlobalConstant{}, errors.New("nil operation body")
	}
	register, ok := op.Body.(RegisterGlobalConstant)
	if !ok {
		return RegisterGlobalConstant{}, errors.Errorf("invalid body type: %T", op.Body)
	}
	return register, nil
}

// SeedNonceRevelation -
func (op Operation) SeedNonceRevelation() (SeedNonceRevelation, error) {
	if op.Kind != KindNonceRevelation {
		return SeedNonceRevelation{}, errors.Errorf("invalid kind of operation: %s", op.Kind)
	}
	if op.Body == nil {
		return SeedNonceRevelation{}, errors.New("nil operation body")
	}
	seed, ok := op.Body.(SeedNonceRevelation)
	if !ok {
		return SeedNonceRevelation{}, errors.Errorf("invalid body type: %T", op.Body)
	}
	return seed, nil
}

// SetDepositsLimit -
func (op Operation) SetDepositsLimit() (SetDepositsLimit, error) {
	if op.Kind != KindSetDepositsLimit {
		return SetDepositsLimit{}, errors.Errorf("invalid kind of operation: %s", op.Kind)
	}
	if op.Body == nil {
		return SetDepositsLimit{}, errors.New("nil operation body")
	}
	tx, ok := op.Body.(SetDepositsLimit)
	if !ok {
		return SetDepositsLimit{}, errors.Errorf("invalid body type: %T", op.Body)
	}
	return tx, nil
}

// Transaction -
func (op Operation) Transaction() (Transaction, error) {
	if op.Kind != KindTransaction {
		return Transaction{}, errors.Errorf("invalid kind of operation: %s", op.Kind)
	}
	if op.Body == nil {
		return Transaction{}, errors.New("nil operation body")
	}
	tx, ok := op.Body.(Transaction)
	if !ok {
		return Transaction{}, errors.Errorf("invalid body type: %T", op.Body)
	}
	return tx, nil
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
	BalanceUpdates           []BalanceUpdate            `json:"balance_updates"`
	OperationResult          OperationResult            `json:"operation_result"`
	InternalOperationResults []InternalOperationResults `json:"internal_operation_results,omitempty"`
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
	StorageSize                  string              `json:"storage_size,omitempty"`
	PaidStorageSizeDiff          string              `json:"paid_storage_size_diff,omitempty"`
	AllocatedDestinationContract bool                `json:"allocated_destination_contract,omitempty"`
	Errors                       []ResultError       `json:"errors,omitempty"`
}

// InternalOperationResults -
type InternalOperationResults struct {
	Kind        string              `json:"kind"`
	Source      string              `json:"source"`
	Nonce       uint64              `json:"nonce"`
	Amount      string              `json:"amount,omitempty"`
	PublicKey   string              `json:"public_key,omitempty"`
	Destination string              `json:"destination,omitempty"`
	Balance     string              `json:"balance,omitempty"`
	Delegate    string              `json:"delegate,omitempty"`
	Script      *stdJSON.RawMessage `json:"script,omitempty"`
	Parameters  *Parameters         `json:"paramaters,omitempty"`
	Result      OperationResult     `json:"result"`
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
