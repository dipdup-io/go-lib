package node

import (
	stdJSON "encoding/json"

	"github.com/pkg/errors"
)

// InvalidBlock -
type InvalidBlock struct {
	Hash   string        `json:"block"`
	Level  uint64        `json:"level"`
	Errors []interface{} `json:"errors"`
}

// Bootstrapped -
type Bootstrapped struct {
	Bootstrapped bool   `json:"bootstrapped"`
	SyncState    string `json:"sync_state"`
}

// Caboose -
type Caboose struct {
	Hash  string `json:"block_hash"`
	Level uint64 `json:"level"`
}

// Checkpoint -
type Checkpoint struct {
	Hash  string `json:"block_hash"`
	Level uint64 `json:"level"`
}

// Savepoint -
type Savepoint struct {
	Hash  string `json:"block_hash"`
	Level uint64 `json:"level"`
}

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

// IsManager -
func IsManager(kind string) bool {
	return kind == KindDelegation || kind == KindOrigination || kind == KindReveal || kind == KindTransaction || kind == KindSetDepositsLimit
}
