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
	Level        int64           `json:"level"`
	Hash         string          `json:"hash"`
	Protocol     string          `json:"protocol"`
	Timestamp    time.Time       `json:"timestamp"`
	VotingEpoch  int64           `json:"votingEpoch"`
	VotingPeriod int64           `json:"votingPeriod"`
	KnownLevel   int64           `json:"knownLevel"`
	LastSync     time.Time       `json:"lastSync"`
	Synced       bool            `json:"synced"`
	QuoteLevel   int64           `json:"quoteLevel"`
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
	Cycle       int64     `json:"cycle"`
	Level       int64     `json:"level"`
	Hash        string    `json:"hash"`
	Timestamp   time.Time `json:"timestamp"`
	Proto       int64     `json:"proto"`
	Priority    int64     `json:"priority"`
	Validations int64     `json:"validations"`
	Deposit     int64     `json:"deposit"`
	Reward      int64     `json:"reward"`
	Fees        int64     `json:"fees"`
	LbEscapeEma int64     `json:"lbEscapeEma"`
	Baker       Address   `json:"baker"`
	Software    struct {
		Version string    `json:"version"`
		Date    time.Time `json:"date"`
	} `json:"software"`
	NonceRevealed bool `json:"nonceRevealed"`
	LbEscapeVote  bool `json:"lbEscapeVote"`
}

// Address -
type Address struct {
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

// BigMapUpdate -
type BigMapUpdate struct {
	ID        int64     `json:"id"`
	Level     int64     `json:"level"`
	Timestamp time.Time `json:"timestamp"`
	Bigmap    int64     `json:"bigmap"`
	Contract  Address   `json:"contract"`
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

	Sender        Address         `json:"sender"`
	Target        Address         `json:"target"`
	Amount        decimal.Decimal `json:"amount"`
	Parameter     Parameter       `json:"parameter"`
	Timestamp     time.Time       `json:"timestamp"`
	ID            int64           `json:"id"`
	Level         int64           `json:"level"`
	Counter       int64           `json:"counter"`
	GasLimit      int64           `json:"gasLimit"`
	GasUsed       int64           `json:"gasUsed"`
	StorageLimit  int64           `json:"storageLimit"`
	StorageUsed   int64           `json:"storageUsed"`
	BakerFee      int64           `json:"bakerFee"`
	StorageFee    int64           `json:"storageFee"`
	AllocationFee int64           `json:"allocationFee"`
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
