package data

import (
	"time"

	"github.com/shopspring/decimal"
)

// Block -
type Block struct {
	Cycle              int64     `json:"cycle"`
	Level              uint64    `json:"level"`
	Hash               string    `json:"hash"`
	Timestamp          time.Time `json:"timestamp"`
	Proto              int64     `json:"proto"`
	Validations        int64     `json:"validations"`
	BlockRound         uint64    `json:"blockRound"`
	PayloadRound       uint64    `json:"payloadRound"`
	Deposit            int64     `json:"deposit"`
	RewardDelegated    int64     `json:"rewardDelegated"`
	RewardStakedOwn    int64     `json:"rewardStakedOwn"`
	RewardStakedEdge   int64     `json:"rewardStakedEdge"`
	RewardStakedShared int64     `json:"rewardStakedShared"`
	BonusDelegated     int64     `json:"bonusDelegated"`
	BonusStakedOwn     int64     `json:"bonusStakedOwn"`
	BonusStakedEdge    int64     `json:"bonusStakedEdge"`
	BonusStakedShared  int64     `json:"bonusStakedShared"`
	Fees               int64     `json:"fees"`
	NonceRevealed      bool      `json:"nonceRevealed"`
	LbToggleEma        uint64    `json:"lbToggleEma"`
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
