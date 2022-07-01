package data

import (
	stdJSON "encoding/json"
	"time"
)

// BigMap -
type BigMap struct {
	Ptr        int64              `json:"ptr"`
	Contract   Address            `json:"contract"`
	Path       string             `json:"path"`
	Active     bool               `json:"active"`
	FirstLevel uint64             `json:"firstLevel"`
	LastLevel  uint64             `json:"lastLevel"`
	TotalKeys  uint64             `json:"totalKeys"`
	ActiveKeys uint64             `json:"activeKeys"`
	Updates    uint64             `json:"updates"`
	KeyType    stdJSON.RawMessage `json:"keyType"`
	ValueType  stdJSON.RawMessage `json:"valueType"`
}

// BigMapUpdate -
type BigMapUpdate struct {
	ID        uint64               `json:"id"`
	Level     uint64               `json:"level"`
	Timestamp time.Time            `json:"timestamp"`
	Bigmap    int64                `json:"bigmap"`
	Contract  Address              `json:"contract"`
	Path      string               `json:"path"`
	Action    string               `json:"action"`
	Content   *BigMapUpdateContent `json:"content,omitempty"`
}

// BigMapKeyUpdate -
type BigMapKeyUpdate struct {
	ID        uint64             `json:"id"`
	Level     uint64             `json:"level"`
	Timestamp time.Time          `json:"timestamp"`
	Action    string             `json:"action"`
	Value     stdJSON.RawMessage `json:"value"`
}

// BigMapUpdateContent -
type BigMapUpdateContent struct {
	Hash  string             `json:"hash"`
	Key   stdJSON.RawMessage `json:"key"`
	Value stdJSON.RawMessage `json:"value"`
}

// BigMapKey -
type BigMapKey struct {
	ID         uint64             `json:"id"`
	Active     bool               `json:"active"`
	Hash       string             `json:"hash"`
	Key        string             `json:"key"`
	Value      stdJSON.RawMessage `json:"value"`
	FirstLevel uint64             `json:"firstLevel"`
	LastLevel  uint64             `json:"lastLevel"`
	Updates    uint64             `json:"updates"`
}

// BigMapJSONSchema -
type BigMapJSONSchema struct {
	Name  string     `json:"name"`
	Path  string     `json:"path"`
	Key   JSONSchema `json:"keySchema"`
	Value JSONSchema `json:"valueSchema"`
}
