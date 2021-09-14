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
	ID         uint64      `json:"id" mapstructure:"id"`
	Level      uint64      `json:"level" mapstructure:"level"`
	Hash       string      `json:"hash" mapstructure:"hash"`
	Kind       string      `json:"type" mapstructure:"type"`
	Block      string      `json:"block" mapstructure:"block"`
	Delegate   *Address    `json:"delegate,omitempty" mapstructure:"delegate,omitempty"`
	GasUsed    *uint64     `json:"gasUsed,omitempty" mapstructure:"gasUsed,omitempty"`
	BakerFee   *uint64     `json:"bakerFee,omitempty" mapstructure:"bakerFee,omitempty"`
	Parameters *Parameters `json:"parameter,omitempty" mapstructure:"parameter,omitempty"`
}

// Parameters -
type Parameters struct {
	Entrypoint string             `json:"entrypoint"`
	Value      stdJSON.RawMessage `json:"value"`
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

// Delegate -
type Delegate struct {
	Type                   string    `json:"type"`
	Alias                  string    `json:"alias"`
	Address                string    `json:"address"`
	PublicKey              string    `json:"publicKey"`
	Balance                int64     `json:"balance"`
	FrozenDeposits         int64     `json:"frozenDeposits"`
	FrozenRewards          int64     `json:"frozenRewards"`
	FrozenFees             int64     `json:"frozenFees"`
	Counter                int64     `json:"counter"`
	ActivationLevel        int64     `json:"activationLevel"`
	StakingBalance         int64     `json:"stakingBalance"`
	NumContracts           int64     `json:"numContracts"`
	NumDelegators          int64     `json:"numDelegators"`
	NumBlocks              int64     `json:"numBlocks"`
	NumEndorsements        int64     `json:"numEndorsements"`
	NumBallots             int64     `json:"numBallots"`
	NumProposals           int64     `json:"numProposals"`
	NumActivations         int64     `json:"numActivations"`
	NumDoubleBaking        int64     `json:"numDoubleBaking"`
	NumDoubleEndorsing     int64     `json:"numDoubleEndorsing"`
	NumNonceRevelations    int64     `json:"numNonceRevelations"`
	NumRevelationPenalties int64     `json:"numRevelationPenalties"`
	NumDelegations         int64     `json:"numDelegations"`
	NumOriginations        int64     `json:"numOriginations"`
	NumTransactions        int64     `json:"numTransactions"`
	NumReveals             int64     `json:"numReveals"`
	NumMigrations          int64     `json:"numMigrations"`
	FirstActivity          int64     `json:"firstActivity"`
	LastActivity           int64     `json:"lastActivity"`
	FirstActivityTime      time.Time `json:"firstActivityTime"`
	LastActivityTime       time.Time `json:"lastActivityTime"`
	ActivationTime         time.Time `json:"activationTime"`
	Software               Software  `json:"software"`
	Active                 bool      `json:"active"`
	Revealed               bool      `json:"revealed"`
}

// Software -
type Software struct {
	Version string    `json:"version"`
	Date    time.Time `json:"date"`
}

// Right -
type Right struct {
	Type      string    `json:"type"`
	Cycle     uint64    `json:"cycle"`
	Level     uint64    `json:"level"`
	Timestamp time.Time `json:"timestamp"`
	Slots     uint64    `json:"slots"`
	Baker     Address   `json:"baker"`
	Status    string    `json:"status"`
}
