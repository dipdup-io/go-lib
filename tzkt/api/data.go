package api

import (
	stdJSON "encoding/json"
	"time"
)

const (
	KindTransaction     = "transaction"
	KindOrigination     = "origination"
	KindEndorsement     = "endorsement"
	KindBallot          = "ballot"
	KindProposal        = "proposal"
	KindDoubleBaking    = "double_baking"
	KindDoubleEndorsing = "double_endorsing"
	KindActivation      = "activation"
	KindNonceRevelation = "nonce_revelation"
	KindDelegation      = "delegation"
	KindReveal          = "reveal"
)

// Operation -
type Operation struct {
	ID    uint64 `json:"id"`
	Level uint64 `json:"level"`
	Hash  string `json:"hash"`
	Kind  string `json:"type"`
	Block string `json:"block"`
}

// Address -
type Address struct {
	Alias   string `json:"alias"`
	Address string `json:"address"`
}

// Block -
type Block struct {
	Level         uint64    `json:"level"`
	Hash          string    `json:"hash"`
	Timestamp     time.Time `json:"timestamp"`
	Proto         int64     `json:"proto"`
	Priority      int64     `json:"priority"`
	Validations   int64     `json:"validations"`
	Deposit       int64     `json:"deposit"`
	Reward        int64     `json:"reward"`
	Fees          int64     `json:"fees"`
	NonceRevealed bool      `json:"nonceRevealed"`
	Baker         Address   `json:"baker"`
}

// Head -
type Head struct {
	Level        uint64    `json:"level"`
	Hash         string    `json:"hash"`
	Protocol     string    `json:"protocol"`
	Timestamp    time.Time `json:"timestamp"`
	VotingEpoch  int64     `json:"votingEpoch"`
	VotingPeriod int64     `json:"votingPeriod"`
	KnownLevel   int64     `json:"knownLevel"`
	LastSync     time.Time `json:"lastSync"`
	Synced       bool      `json:"synced"`
	QuoteLevel   int64     `json:"quoteLevel"`
	QuoteBtc     float64   `json:"quoteBtc"`
	QuoteEur     float64   `json:"quoteEur"`
	QuoteUsd     float64   `json:"quoteUsd"`
}

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
