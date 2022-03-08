package node

import (
	stdJSON "encoding/json"

	"github.com/pkg/errors"
)

// Errors
var (
	ErrUnknownKind = errors.New("Unknown operation kind")
)

// MempoolResponse -
type MempoolResponse struct {
	Applied       []Applied `json:"applied"`
	Refused       []Failed  `json:"refused"`
	BranchRefused []Failed  `json:"branch_refused"`
	BranchDelayed []Failed  `json:"branch_delayed"`
}

// Applied -
type Applied struct {
	Hash      string             `json:"hash"`
	Branch    string             `json:"branch"`
	Signature string             `json:"signature"`
	Contents  []Content          `json:"contents"`
	Raw       stdJSON.RawMessage `json:"raw"`
}

// UnmarshalJSON -
func (a *Applied) UnmarshalJSON(data []byte) error {
	type buf Applied
	if err := json.Unmarshal(data, (*buf)(a)); err != nil {
		return err
	}
	a.Raw = data
	return nil
}

// Failed -
type Failed struct {
	Hash      string             `json:"-"`
	Protocol  string             `json:"protocol"`
	Branch    string             `json:"branch"`
	Contents  []Content          `json:"contents"`
	Signature string             `json:"signature,omitempty"`
	Error     stdJSON.RawMessage `json:"error,omitempty"`
	Raw       stdJSON.RawMessage `json:"raw"`
}

// UnmarshalJSON -
func (f *Failed) UnmarshalJSON(data []byte) error {
	var body []stdJSON.RawMessage
	if err := json.Unmarshal(data, &body); err != nil {
		return err
	}
	if len(body) != 2 {
		return errors.Errorf("Invalid failed operation body %s", string(data))
	}
	if err := json.Unmarshal(body[0], &f.Hash); err != nil {
		return err
	}
	type buf Failed
	if err := json.Unmarshal(body[1], (*buf)(f)); err != nil {
		return err
	}
	f.Raw = data
	return nil
}

// FailedMonitor -
type FailedMonitor struct {
	Hash      string             `json:"hash"`
	Protocol  string             `json:"protocol"`
	Branch    string             `json:"branch"`
	Contents  []Content          `json:"contents"`
	Signature string             `json:"signature,omitempty"`
	Error     stdJSON.RawMessage `json:"error,omitempty"`
	Raw       stdJSON.RawMessage `json:"raw"`
}

// UnmarshalJSON -
func (f *FailedMonitor) UnmarshalJSON(data []byte) error {
	type buf FailedMonitor
	if err := json.Unmarshal(data, (*buf)(f)); err != nil {
		return err
	}
	f.Raw = data
	return nil
}

// Contents -
type Content struct {
	Kind string             `json:"kind"`
	Body stdJSON.RawMessage `json:"-"`
}

// UnmarshalJSON -
func (c *Content) UnmarshalJSON(data []byte) error {
	type buf Content
	if err := json.Unmarshal(data, (*buf)(c)); err != nil {
		return err
	}
	c.Body = data
	return nil
}

// HeadMetadata -
type HeadMetadata struct {
	Protocol        string `json:"protocol"`
	NextProtocol    string `json:"next_protocol"`
	TestChainStatus struct {
		Status string `json:"status"`
	} `json:"test_chain_status"`
	MaxOperationsTTL       uint64 `json:"max_operations_ttl"`
	MaxOperationDataLength uint64 `json:"max_operation_data_length"`
	MaxBlockHeaderLength   uint64 `json:"max_block_header_length"`
	MaxOperationListLength []struct {
		MaxSize uint64 `json:"max_size"`
		MaxOp   uint64 `json:"max_op,omitempty"`
	} `json:"max_operation_list_length"`
	Baker string `json:"baker"`
	Level struct {
		Level                uint64 `json:"level"`
		LevelPosition        uint64 `json:"level_position"`
		Cycle                uint64 `json:"cycle"`
		CyclePosition        uint64 `json:"cycle_position"`
		VotingPeriod         uint64 `json:"voting_period"`
		VotingPeriodPosition uint64 `json:"voting_period_position"`
		ExpectedCommitment   bool   `json:"expected_commitment"`
	} `json:"level"`
	LevelInfo struct {
		Level              uint64 `json:"level"`
		LevelPosition      uint64 `json:"level_position"`
		Cycle              uint64 `json:"cycle"`
		CyclePosition      uint64 `json:"cycle_position"`
		ExpectedCommitment bool   `json:"expected_commitment"`
	} `json:"level_info"`
	VotingPeriodKind string `json:"voting_period_kind"`
	VotingPeriodInfo struct {
		VotingPeriod struct {
			Index         uint64 `json:"index"`
			Kind          string `json:"kind"`
			StartPosition uint64 `json:"start_position"`
		} `json:"voting_period"`
		Position  int `json:"position"`
		Remaining int `json:"remaining"`
	} `json:"voting_period_info"`
	NonceHash      interface{}     `json:"nonce_hash"`
	ConsumedGas    string          `json:"consumed_gas"`
	Deactivated    []interface{}   `json:"deactivated"`
	BalanceUpdates []BalanceUpdate `json:"balance_updates"`
}

// IsManager -
func IsManager(kind string) bool {
	return kind == KindDelegation || kind == KindOrigination || kind == KindReveal || kind == KindTransaction
}

// InjectOperationRequest -
type InjectOperationRequest struct {
	Operation string
	ChainID   string
	Async     bool
}
