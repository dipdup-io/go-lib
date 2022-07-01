package data

import (
	"encoding/json"
	stdJSON "encoding/json"
	"time"
)

// Address -
type Address struct {
	Alias   string `json:"alias"`
	Address string `json:"address"`
}

// Contract -
type Contract struct {
	Kind     string `json:"kind"`
	Alias    string `json:"alias,omitempty"`
	Address  string `json:"address,omitempty"`
	TypeHash int    `json:"typeHash"`
	CodeHash int    `json:"codeHash"`
}

// JSONSchema -
type JSONSchema struct {
	Schema               string                `json:"$schema,omitempty"`
	Type                 string                `json:"type,omitempty"`
	Comment              string                `json:"$comment,omitempty"`
	Required             []string              `json:"required,omitempty"`
	Properties           map[string]JSONSchema `json:"properties,omitempty"`
	OneOf                []JSONSchema          `json:"oneOf"`
	AdditionalProperties AdditionalProperties  `json:"additionalProperties,omitempty"`
	PropertyNames        *JSONSchema           `json:"propertyNames,omitempty"`
	Items                *JSONSchema           `json:"items,omitempty"`
}

// AdditionalProperties -
type AdditionalProperties struct {
	Value *JSONSchema `json:"-"`
}

// UnmarshalJSON -
func (props *AdditionalProperties) UnmarshalJSON(data []byte) error {
	var flag bool
	if err := json.Unmarshal(data, &flag); err == nil {
		props.Value = nil
		return nil
	}

	props.Value = &JSONSchema{}
	return json.Unmarshal(data, props.Value)
}

// Transfer -
type Transfer struct {
	ID            uint64    `json:"id"`
	Level         uint64    `json:"level"`
	Timestamp     time.Time `json:"timestamp"`
	Token         Token     `json:"token"`
	From          *Address  `json:"from,omitempty"`
	To            *Address  `json:"to,omitempty"`
	Amount        string    `json:"amount"`
	TransactionID *uint64   `json:"transactionId,omitempty"`
	OriginationID *uint64   `json:"originationId,omitempty"`
	MigrationID   *uint64   `json:"migrationId,omitempty"`
}

// Token -
type Token struct {
	ID       uint64             `json:"id"`
	Contract Address            `json:"contract"`
	TokenID  string             `json:"tokenId"`
	Standard string             `json:"standard"`
	Metadata stdJSON.RawMessage `json:"metadata,omitempty"`
}
