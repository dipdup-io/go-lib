package data

import "time"

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

// AccountMetadata -
type AccountMetadata struct {
	Address     string `json:"address"`
	Kind        string `json:"kind"`
	Alias       string `json:"alias"`
	Description string `json:"description"`
	Site        string `json:"site"`
	Support     string `json:"support"`
	Email       string `json:"email"`
	Twitter     string `json:"twitter"`
	Telegram    string `json:"telegram"`
	Discord     string `json:"discord"`
	Reddit      string `json:"reddit"`
	Slack       string `json:"slack"`
	Github      string `json:"github"`
	Gitlab      string `json:"gitlab"`
	Instagram   string `json:"instagram"`
	Facebook    string `json:"facebook"`
	Medium      string `json:"medium"`
}

// MetadataConstraint -
type MetadataConstraint interface {
	AccountMetadata
}

// Metadata -
type Metadata[T MetadataConstraint] struct {
	Key   string `json:"key"`
	Value T      `json:"extras"`
}

// Account -
type Account struct {
	Type              string    `json:"type"`
	Address           string    `json:"address"`
	Kind              string    `json:"kind"`
	Tzips             []string  `json:"tzips"`
	Alias             string    `json:"alias"`
	Balance           int64     `json:"balance"`
	Creator           Address   `json:"creator"`
	NumContracts      int64     `json:"numContracts"`
	NumDelegations    int64     `json:"numDelegations"`
	NumOriginations   int64     `json:"numOriginations"`
	NumTransactions   int64     `json:"numTransactions"`
	NumReveals        int64     `json:"numReveals"`
	NumMigrations     int64     `json:"numMigrations"`
	FirstActivity     int64     `json:"firstActivity"`
	FirstActivityTime time.Time `json:"firstActivityTime"`
	LastActivity      int64     `json:"lastActivity"`
	LastActivityTime  time.Time `json:"lastActivityTime"`
	TypeHash          int64     `json:"typeHash"`
	CodeHash          int64     `json:"codeHash"`
}
