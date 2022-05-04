package api

import (
	stdJSON "encoding/json"
	"time"

	"github.com/shopspring/decimal"
)

const (
	KindTransaction            = "transaction"
	KindOrigination            = "origination"
	KindEndorsement            = "endorsement"
	KindBallot                 = "ballot"
	KindProposal               = "proposal"
	KindDoubleBaking           = "double_baking"
	KindDoubleEndorsing        = "double_endorsing"
	KindActivation             = "activation"
	KindNonceRevelation        = "nonce_revelation"
	KindDelegation             = "delegation"
	KindReveal                 = "reveal"
	KindRegisterGlobalConstant = "register_constant"
)

// urls
const (
	BaseURL = "https://api.tzkt.io"
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

// Origination -
type Origination struct {
	Type             string    `json:"type"`
	ID               uint64    `json:"id"`
	Level            uint64    `json:"level"`
	Timestamp        time.Time `json:"timestamp"`
	Block            string    `json:"block"`
	Hash             string    `json:"hash"`
	Counter          uint64    `json:"counter"`
	Sender           Address   `json:"sender"`
	Initiator        Address   `json:"initiator"`
	Nonce            *uint64   `json:"nonce,omitempty"`
	GasLimit         uint64    `json:"gasLimit"`
	GasUsed          uint64    `json:"gasUsed"`
	StorageLimit     uint64    `json:"storageLimit"`
	StorageUsed      uint64    `json:"storageUsed"`
	BakerFee         uint64    `json:"bakerFee"`
	StorageFee       uint64    `json:"storageFee"`
	AllocationFee    uint64    `json:"allocationFee"`
	ContractBalance  uint64    `json:"contractBalance"`
	ContractManager  Address   `json:"contractManager"`
	ContractDelegate Address   `json:"contractDelegate"`
	Status           string    `json:"status"`
	Originated       Contract  `json:"originatedContract"`
	Errors           []Error   `json:"errors,omitempty"`
}

// Transaction -
type Transaction struct {
	Type          string          `json:"type"`
	Sender        Address         `json:"sender"`
	Target        Address         `json:"target"`
	Initiator     Address         `json:"initiator"`
	Amount        decimal.Decimal `json:"amount"`
	Parameter     *Parameters     `json:"parameter"`
	Timestamp     time.Time       `json:"timestamp"`
	ID            uint64          `json:"id"`
	Level         uint64          `json:"level"`
	Counter       uint64          `json:"counter"`
	GasLimit      uint64          `json:"gasLimit"`
	GasUsed       uint64          `json:"gasUsed"`
	StorageLimit  uint64          `json:"storageLimit"`
	StorageUsed   uint64          `json:"storageUsed"`
	BakerFee      uint64          `json:"bakerFee"`
	StorageFee    uint64          `json:"storageFee"`
	AllocationFee uint64          `json:"allocationFee"`
	Status        string          `json:"status"`
	Parameters    string          `json:"parameters"`
	Block         string          `json:"block"`
	Hash          string          `json:"hash"`
	HasInternals  bool            `json:"hasInternals"`
	Nonce         *uint64         `json:"nonce,omitempty"`
}

// Delegation -
type Delegation struct {
	Block       string          `json:"block"`
	Hash        string          `json:"hash"`
	Type        string          `json:"type"`
	Status      string          `json:"status"`
	Sender      Address         `json:"sender"`
	NewDelegate Address         `json:"newDelegate"`
	Timestamp   time.Time       `json:"timestamp"`
	Amount      decimal.Decimal `json:"amount"`
	ID          uint64          `json:"id"`
	Level       uint64          `json:"level"`
	Counter     uint64          `json:"counter"`
	GasLimit    uint64          `json:"gasLimit"`
	GasUsed     uint64          `json:"gasUsed"`
	BakerFee    uint64          `json:"bakerFee"`
	Nonce       *uint64         `json:"nonce,omitempty"`
}

// Reveal -
type Reveal struct {
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	Block     string    `json:"block"`
	Hash      string    `json:"hash"`
	Timestamp time.Time `json:"timestamp"`
	Sender    Address   `json:"sender"`
	ID        uint64    `json:"id"`
	Level     uint64    `json:"level"`
	Counter   uint64    `json:"counter"`
	GasLimit  uint64    `json:"gasLimit"`
	GasUsed   uint64    `json:"gasUsed"`
	BakerFee  uint64    `json:"bakerFee"`
	Nonce     *uint64   `json:"nonce,omitempty"`
}

// Error -
type Error struct {
	Type string `json:"type"`
}

// Address -
type Address struct {
	Alias   string `json:"alias"`
	Address string `json:"address"`
}

// Contract -
type Contract struct {
	Kind     string `json:"kind"`
	Address  string `json:"address"`
	TypeHash int    `json:"typeHash"`
	CodeHash int    `json:"codeHash"`
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

// ContractJSONSchema -
type ContractJSONSchema struct {
	Storage     JSONSchema             `json:"storageSchema"`
	Entrypoints []EntrypointJSONSchema `json:"entrypoints"`
	BigMaps     []BigMapJSONSchema     `json:"bigMaps"`
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

// EntrypointJSONSchema -
type EntrypointJSONSchema struct {
	Name      string     `json:"name"`
	Parameter JSONSchema `json:"parameterSchema"`
}

// BigMapJSONSchema -
type BigMapJSONSchema struct {
	Name  string     `json:"name"`
	Path  string     `json:"path"`
	Key   JSONSchema `json:"keySchema"`
	Value JSONSchema `json:"valueSchema"`
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

// Migration -
type Migration struct {
	Type          string    `json:"type"`
	ID            uint64    `json:"id"`
	Level         uint64    `json:"level"`
	Timestamp     time.Time `json:"timestamp"`
	Block         string    `json:"block"`
	Kind          string    `json:"kind"`
	Account       Address   `json:"account"`
	BalanceChange int64     `json:"balanceChange"`
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
	Value T      `json:"metadata"`
}
