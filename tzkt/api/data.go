package api

import (
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
	Baker         struct {
		Alias   string `json:"alias"`
		Address string `json:"address"`
	} `json:"baker"`
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
