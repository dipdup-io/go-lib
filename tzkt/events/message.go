package events

import (
	stdJSON "encoding/json"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

// MessageType - TzKT message type
type MessageType int

// message types
const (
	MessageTypeState MessageType = iota
	MessageTypeData
	MessageTypeReorg
	MessageTypeSubscribed
)

// Message - message struct
type Message struct {
	Channel string
	Type    MessageType `json:"type"`
	State   uint64      `json:"state"`
	Body    interface{} `json:"data"`
}

// String -
func (msg Message) String() string {
	s := fmt.Sprintf("channel=%s type=%d state=%d", msg.Channel, msg.Type, msg.State)
	if msg.Body != nil {
		s = fmt.Sprintf("%s data=%v", s, msg.Body)
	}
	return s
}

// Packet -
type Packet struct {
	Type  MessageType        `json:"type"`
	State uint64             `json:"state"`
	Data  stdJSON.RawMessage `json:"data,omitempty"`
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

// Block -
type Block struct {
	Cycle       uint64    `json:"cycle"`
	Level       uint64    `json:"level"`
	Hash        string    `json:"hash"`
	Timestamp   time.Time `json:"timestamp"`
	Proto       int64     `json:"proto"`
	Priority    int64     `json:"priority"`
	Validations int64     `json:"validations"`
	Deposit     int64     `json:"deposit"`
	Reward      int64     `json:"reward"`
	Fees        int64     `json:"fees"`
	LbEscapeEma int64     `json:"lbEscapeEma"`
	Baker       Alias     `json:"baker"`
	Software    struct {
		Version string    `json:"version"`
		Date    time.Time `json:"date"`
	} `json:"software"`
	NonceRevealed bool `json:"nonceRevealed"`
	LbEscapeVote  bool `json:"lbEscapeVote"`
}

// Alias -
type Alias struct {
	Alias   string `json:"alias"`
	Address string `json:"address"`
}

// Account -
type Account struct {
	Type              string    `json:"type"`
	Address           string    `json:"address"`
	Kind              string    `json:"kind"`
	Tzips             []string  `json:"tzips"`
	Alias             string    `json:"alias"`
	Balance           int64     `json:"balance"`
	Creator           Alias     `json:"creator"`
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

// BigMapUpdate -
type BigMapUpdate struct {
	ID        uint64    `json:"id"`
	Level     uint64    `json:"level"`
	Timestamp time.Time `json:"timestamp"`
	Bigmap    int64     `json:"bigmap"`
	Contract  Alias     `json:"contract"`
	Path      string    `json:"path"`
	Action    string    `json:"action"`
	Content   *Content  `json:"content,omitempty"`
}

// Content -
type Content struct {
	Hash  string             `json:"hash"`
	Key   string             `json:"key"`
	Value stdJSON.RawMessage `json:"value"`
}

// Operation -
type Operation struct {
	Type string `json:"type"`
}

// Transaction -
type Transaction struct {
	Operation

	Sender        Alias           `json:"sender"`
	Target        Alias           `json:"target"`
	Initiator     Alias           `json:"initiator"`
	Amount        decimal.Decimal `json:"amount"`
	Parameter     *Parameter      `json:"parameter"`
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
}

// Parameter -
type Parameter struct {
	Entrypoint string             `json:"entrypoint"`
	Value      stdJSON.RawMessage `json:"value"`
}

// Origination -
type Origination struct {
	Operation

	ID                 uint64              `json:"id"`
	Level              uint64              `json:"level"`
	Timestamp          time.Time           `json:"timestamp"`
	Block              string              `json:"block"`
	Hash               string              `json:"hash"`
	Counter            uint64              `json:"counter"`
	Initiator          *Alias              `json:"initiator"`
	Sender             *Alias              `json:"sender"`
	Nonce              *uint64             `json:"nonce"`
	GasLimit           uint64              `json:"gasLimit"`
	GasUsed            uint64              `json:"gasUsed"`
	StorageLimit       uint64              `json:"storageLimit"`
	StorageUsed        uint64              `json:"storageUsed"`
	BakerFee           uint64              `json:"bakerFee"`
	StorageFee         uint64              `json:"storageFee"`
	AllocationFee      uint64              `json:"allocationFee"`
	ContractBalance    uint64              `json:"contractBalance"`
	ContractManager    *Alias              `json:"contractManager"`
	ContractDelegate   *Alias              `json:"contractDelegate"`
	Code               stdJSON.RawMessage  `json:"code"`
	Storage            stdJSON.RawMessage  `json:"storage"`
	Diffs              stdJSON.RawMessage  `json:"diffs"`
	Status             string              `json:"status"`
	Errors             stdJSON.RawMessage  `json:"errors,omitempty"`
	OriginatedContract *OriginatedContract `json:"originatedContract,omitempty"`
	Quote              *QuoteShort         `json:"quote,omitempty"`
}

// Reveal -
type Reveal struct {
	Operation

	ID        uint64             `json:"id"`
	Level     uint64             `json:"level"`
	Timestamp time.Time          `json:"timestamp"`
	Block     string             `json:"block"`
	Hash      string             `json:"hash"`
	Sender    *Alias             `json:"sender"`
	Counter   uint64             `json:"counter"`
	GasLimit  uint64             `json:"gasLimit"`
	GasUsed   uint64             `json:"gasUsed"`
	BakerFee  uint64             `json:"bakerFee"`
	Status    string             `json:"status"`
	Errors    stdJSON.RawMessage `json:"errors,omitempty"`
	Quote     *QuoteShort        `json:"quote,omitempty"`
}

// Delegation -
type Delegation struct {
	Operation

	ID           uint64             `json:"id"`
	Level        uint64             `json:"level"`
	Timestamp    time.Time          `json:"timestamp"`
	Block        string             `json:"block"`
	Hash         string             `json:"hash"`
	Counter      uint64             `json:"counter"`
	Initiator    *Account           `json:"initiator"`
	Sender       *Account           `json:"sender"`
	Nonce        uint64             `json:"nonce"`
	GasLimit     uint64             `json:"gasLimit"`
	GasUsed      uint64             `json:"gasUsed"`
	BakerFee     uint64             `json:"bakerFee"`
	Amount       uint64             `json:"amount"`
	PrevDelegate *Account           `json:"prevDelegate"`
	NewDelegate  *Account           `json:"newDelegate"`
	Status       string             `json:"status"`
	Errors       stdJSON.RawMessage `json:"errors,omitempty"`
	Quote        *QuoteShort        `json:"quote,omitempty"`
}

// OriginatedContract -
type OriginatedContract struct {
	Kind     string `json:"kind"`
	Alias    string `json:"alias,omitempty"`
	Address  string `json:"address,omitempty"`
	TypeHash int64  `json:"typeHash"`
	CodeHash int64  `json:"codeHash"`
}

// QuoteShort -
type QuoteShort struct {
	BTC decimal.Decimal `json:"btc,omitempty"`
	EUR decimal.Decimal `json:"eur,omitempty"`
	USD decimal.Decimal `json:"usd,omitempty"`
	CNY decimal.Decimal `json:"cny,omitempty"`
	JPY decimal.Decimal `json:"jpy,omitempty"`
	KRW decimal.Decimal `json:"krw,omitempty"`
	ETH decimal.Decimal `json:"eth,omitempty"`
}
