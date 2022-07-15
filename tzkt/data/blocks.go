package data

import (
	"time"

	"github.com/shopspring/decimal"
)

// Block -
type Block struct {
	Cycle         int64     `json:"cycle"`
	Level         uint64    `json:"level"`
	Hash          string    `json:"hash"`
	Timestamp     time.Time `json:"timestamp"`
	Proto         int64     `json:"proto"`
	Priority      int64     `json:"priority"`
	Validations   int64     `json:"validations"`
	BlockRound    uint64    `json:"blockRound"`
	PayloadRound  uint64    `josn:"payloadRound"`
	Deposit       int64     `json:"deposit"`
	Reward        int64     `json:"reward"`
	Fees          int64     `json:"fees"`
	Bonus         int64     `json:"bonus"`
	LbEscapeEma   int64     `json:"lbEscapeEma"`
	NonceRevealed bool      `json:"nonceRevealed"`
	LbEscapeVote  bool      `json:"lbEscapeVote"`
	LbToggleEma   uint64    `json:"lbToggleEma"`
	Baker         Address   `json:"baker"`
}

// Head -
type Head struct {
	Chain        string          `json:"chain"`
	ChainID      string          `json:"chainId"`
	Cycle        int64           `json:"cycle"`
	Level        uint64          `json:"level"`
	Hash         string          `json:"hash"`
	Protocol     string          `json:"protocol"`
	Timestamp    time.Time       `json:"timestamp"`
	VotingEpoch  int64           `json:"votingEpoch"`
	VotingPeriod int64           `json:"votingPeriod"`
	KnownLevel   uint64          `json:"knownLevel"`
	LastSync     time.Time       `json:"lastSync"`
	Synced       bool            `json:"synced"`
	QuoteLevel   uint64          `json:"quoteLevel"`
	QuoteBtc     decimal.Decimal `json:"quoteBtc"`
	QuoteEur     decimal.Decimal `json:"quoteEur"`
	QuoteUsd     decimal.Decimal `json:"quoteUsd"`
	QuoteCny     decimal.Decimal `json:"quoteCny"`
	QuoteJpy     decimal.Decimal `json:"quoteJpy"`
	QuoteKrw     decimal.Decimal `json:"quoteKrw"`
	QuoteEth     decimal.Decimal `json:"quoteEth"`
}
