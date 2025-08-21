package data

import (
	"time"

	stdJSON "encoding/json"

	"github.com/shopspring/decimal"
)

// OperationConstraint -
type OperationConstraint interface {
	Transaction | Origination | Delegation | Reveal | RegisterConstant | Endorsement | Preendorsement |
		Ballot | Proposal | Activation | TransferTicket | TxRollupCommit | TxRollupDispatchTicket |
		TxRollupFinalizeCommitment | TxRollupOrigination | TxRollupRejection | TxRollupRemoveCommitment |
		TxRollupReturnBond | TxRollupSubmitBatch | NonceRevelation | DoubleBaking | DoubleEndorsing | SetDepositsLimit |
		DoublePreendorsing | Baking | RevelationPenalty | EndorsingReward | VdfRevelation | IncreasePaidStorage |
		DrainDelegate | UpdateConsensusKey | SmartRollupAddMessage | SmartRollupCement | SmartRollupExecute |
		SmartRollupOriginate | SmartRollupPublish | SmartRollupRefute | SmartRollupRecoverBond | DalPublishCommitment
}

// Operation -
type Operation struct {
	ID         uint64      `json:"id"`
	Level      uint64      `json:"level"`
	Hash       string      `json:"hash"`
	Type       string      `json:"type"`
	Block      string      `json:"block"`
	Delegate   *Address    `json:"delegate,omitempty"`
	GasUsed    *uint64     `json:"gasUsed,omitempty"`
	BakerFee   *uint64     `json:"bakerFee,omitempty"`
	Parameters *Parameters `json:"parameter,omitempty"`
}

// Parameters -
type Parameters struct {
	Entrypoint string             `json:"entrypoint"`
	Value      stdJSON.RawMessage `json:"value"`
}

// Origination -
type Origination struct {
	Type             string              `json:"type"`
	ID               uint64              `json:"id"`
	Level            uint64              `json:"level"`
	Timestamp        time.Time           `json:"timestamp"`
	Block            string              `json:"block"`
	Hash             string              `json:"hash"`
	Counter          uint64              `json:"counter"`
	Sender           *Address            `json:"sender"`
	Initiator        *Address            `json:"initiator"`
	Nonce            *uint64             `json:"nonce,omitempty"`
	GasLimit         uint64              `json:"gasLimit"`
	GasUsed          uint64              `json:"gasUsed"`
	StorageLimit     uint64              `json:"storageLimit"`
	StorageUsed      uint64              `json:"storageUsed"`
	BakerFee         uint64              `json:"bakerFee"`
	StorageFee       uint64              `json:"storageFee"`
	AllocationFee    uint64              `json:"allocationFee"`
	ContractBalance  uint64              `json:"contractBalance"`
	ContractManager  *Address            `json:"contractManager"`
	ContractDelegate *Address            `json:"contractDelegate"`
	Code             stdJSON.RawMessage  `json:"code"`
	Storage          stdJSON.RawMessage  `json:"storage"`
	Diffs            stdJSON.RawMessage  `json:"diffs"`
	Status           string              `json:"status"`
	Originated       *OriginatedContract `json:"originatedContract,omitempty"`
	Errors           []Error             `json:"errors,omitempty"`
	Quote            *Quote              `json:"quote,omitempty"`
}

// Transaction -
type Transaction struct {
	Type          string          `json:"type"`
	Sender        Address         `json:"sender"`
	Target        Address         `json:"target"`
	Initiator     Address         `json:"initiator"`
	Amount        decimal.Decimal `json:"amount"`
	Parameter     *Parameters     `json:"parameter,omitempty"`
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
	Quote         *Quote          `json:"quote,omitempty"`
}

// Delegation -
type Delegation struct {
	Block        string          `json:"block"`
	Hash         string          `json:"hash"`
	Type         string          `json:"type"`
	Status       string          `json:"status"`
	Sender       *Address        `json:"sender,omitempty"`
	NewDelegate  *Address        `json:"newDelegate,omitempty"`
	Initiator    *Address        `json:"initiator,omitempty"`
	PrevDelegate *Address        `json:"prevDelegate,omitempty"`
	Timestamp    time.Time       `json:"timestamp"`
	Amount       decimal.Decimal `json:"amount"`
	ID           uint64          `json:"id"`
	Level        uint64          `json:"level"`
	Counter      uint64          `json:"counter"`
	GasLimit     uint64          `json:"gasLimit"`
	GasUsed      uint64          `json:"gasUsed"`
	BakerFee     uint64          `json:"bakerFee"`
	Nonce        *uint64         `json:"nonce,omitempty"`
	Quote        *Quote          `json:"quote,omitempty"`
}

// Reveal -
type Reveal struct {
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	Block     string    `json:"block"`
	Hash      string    `json:"hash"`
	Timestamp time.Time `json:"timestamp"`
	Sender    *Address  `json:"sender,omitempty"`
	ID        uint64    `json:"id"`
	Level     uint64    `json:"level"`
	Counter   uint64    `json:"counter"`
	GasLimit  uint64    `json:"gasLimit"`
	GasUsed   uint64    `json:"gasUsed"`
	BakerFee  uint64    `json:"bakerFee"`
	Nonce     *uint64   `json:"nonce,omitempty"`
	Quote     *Quote    `json:"quote,omitempty"`
}

// Endorsement -
type Endorsement struct {
	Type      string    `json:"type"`
	Block     string    `json:"block"`
	Hash      string    `json:"hash"`
	ID        uint64    `json:"id"`
	Level     uint64    `json:"level"`
	Deposit   int64     `json:"deposit"`
	Rewards   int64     `json:"rewards"`
	Slots     int       `json:"slots"`
	Timestamp time.Time `json:"timestamp"`
	Delegate  Address   `json:"delegate"`
	Quote     *Quote    `json:"quote,omitempty"`
}

// Preendorsement -
type Preendorsement struct {
	Type      string    `json:"type"`
	Block     string    `json:"block"`
	Hash      string    `json:"hash"`
	ID        uint64    `json:"id"`
	Level     uint64    `json:"level"`
	Timestamp time.Time `json:"timestamp"`
	Delegate  Address   `json:"delegate"`
	Slots     int       `json:"slots"`
	Quote     *Quote    `json:"quote,omitempty"`
}

// Ballot -
type Ballot struct {
	Type        string        `json:"type"`
	Block       string        `json:"block"`
	Hash        string        `json:"hash"`
	Vote        string        `json:"vote"`
	ID          uint64        `json:"id"`
	Level       uint64        `json:"level"`
	Timestamp   time.Time     `json:"timestamp"`
	Period      PeriodInfo    `json:"period"`
	Proposal    ProposalAlias `json:"proposal"`
	Delegate    Address       `json:"delegate"`
	VotingPower int64         `json:"votingPower"`
	Quote       *Quote        `json:"quote,omitempty"`
}

// Proposal -
type Proposal struct {
	Type        string        `json:"type"`
	ID          uint64        `json:"id"`
	Level       uint64        `json:"level"`
	Timestamp   time.Time     `json:"timestamp"`
	Block       string        `json:"block"`
	Hash        string        `json:"hash"`
	Period      PeriodInfo    `json:"period"`
	Proposal    ProposalAlias `json:"proposal"`
	Delegate    Address       `json:"delegate"`
	VotingPower int64         `json:"votingPower"`
	Duplicated  bool          `json:"duplicated"`
	Quote       *Quote        `json:"quote,omitempty"`
}

// Activation -
type Activation struct {
	Type      string    `json:"type"`
	ID        uint64    `json:"id"`
	Level     uint64    `json:"level"`
	Timestamp time.Time `json:"timestamp"`
	Block     string    `json:"block"`
	Hash      string    `json:"hash"`
	Account   Address   `json:"account"`
	Balance   int64     `json:"balance"`
	Quote     *Quote    `json:"quote,omitempty"`
}

// RegisterConstant -
type RegisterConstant struct {
	Type         string             `json:"type"`
	Block        string             `json:"block"`
	Hash         string             `json:"hash"`
	Status       string             `json:"status"`
	Address      string             `json:"address"`
	ID           uint64             `json:"id"`
	Level        uint64             `json:"level"`
	Counter      uint64             `json:"counter"`
	GasLimit     uint64             `json:"gasLimit"`
	GasUsed      uint64             `json:"gasUsed"`
	StorageLimit uint64             `json:"storageLimit"`
	StorageUsed  uint64             `json:"storageUsed"`
	BakerFee     uint64             `json:"bakerFee"`
	StorageFee   uint64             `json:"storageFee"`
	Value        stdJSON.RawMessage `json:"value"`
	Timestamp    time.Time          `json:"timestamp"`
	Sender       Address            `json:"sender"`
	Errors       []Error            `json:"errors,omitempty"`
	Quote        *Quote             `json:"quote,omitempty"`
}

// NonceRevelation -
type NonceRevelation struct {
	Type          string    `json:"type"`
	ID            uint64    `json:"id"`
	Level         uint64    `json:"level"`
	Timestamp     time.Time `json:"timestamp"`
	Block         string    `json:"block"`
	Hash          string    `json:"hash"`
	Baker         Address   `json:"baker"`
	Sender        Address   `json:"sender"`
	RevealedLevel int       `json:"revealedLevel"`
	RevealedCycle int       `json:"revealedCycle"`
	Nonce         string    `json:"nonce"`
	Reward        int64     `json:"reward"`
	Quote         *Quote    `json:"quote,omitempty"`
	BakerRewards  int64     `json:"bakerRewards"`
}

// ProposalAlias -
type ProposalAlias struct {
	Alias string `json:"alias"`
	Hash  string `json:"hash"`
}

// PeriodInfo -
type PeriodInfo struct {
	Index      int    `json:"index"`
	Epoch      int    `json:"epoch"`
	Kind       string `json:"kind"`
	FirstLevel int    `json:"firstLevel"`
	LastLevel  int    `json:"lastLevel"`
}

// TransferTicket -
type TransferTicket struct {
	Type         string             `json:"type"`
	Block        string             `json:"block"`
	Hash         string             `json:"hash"`
	Entrypoint   string             `json:"entrypoint"`
	Status       string             `json:"status"`
	Timestamp    time.Time          `json:"timestamp"`
	Sender       Address            `json:"sender"`
	ID           uint64             `json:"id"`
	Level        uint64             `json:"level"`
	Counter      uint64             `json:"counter"`
	GasLimit     uint64             `json:"gasLimit"`
	GasUsed      uint64             `json:"gasUsed"`
	StorageLimit uint64             `json:"storageLimit"`
	StorageUsed  uint64             `json:"storageUsed"`
	BakerFee     uint64             `json:"bakerFee"`
	StorageFee   uint64             `json:"storageFee"`
	Target       Address            `json:"target"`
	Ticketer     Address            `json:"ticketer"`
	Amount       decimal.Decimal    `json:"amount"`
	ContentType  stdJSON.RawMessage `json:"contentType"`
	Content      stdJSON.RawMessage `json:"content"`
	Errors       []Error            `json:"errors,omitempty"`
	Quote        *Quote             `json:"quote,omitempty"`
}

// TxRollupCommit -
type TxRollupCommit struct {
	Timestamp    time.Time `json:"timestamp"`
	Type         string    `json:"type"`
	Block        string    `json:"block"`
	Hash         string    `json:"hash"`
	Status       string    `json:"status"`
	ID           uint64    `json:"id"`
	Level        uint64    `json:"level"`
	Counter      uint64    `json:"counter"`
	GasLimit     uint64    `json:"gasLimit"`
	GasUsed      uint64    `json:"gasUsed"`
	StorageLimit uint64    `json:"storageLimit"`
	StorageUsed  uint64    `json:"storageUsed"`
	BakerFee     uint64    `json:"bakerFee"`
	Bond         uint64    `json:"bond"`
	Sender       Address   `json:"sender"`
	Rollup       Address   `json:"rollup"`
	Errors       []Error   `json:"errors,omitempty"`
	Quote        *Quote    `json:"quote,omitempty"`
}

// TxRollupDispatchTicket -
type TxRollupDispatchTicket struct {
	Timestamp    time.Time `json:"timestamp"`
	Type         string    `json:"type"`
	Block        string    `json:"block"`
	Hash         string    `json:"hash"`
	Status       string    `json:"status"`
	Counter      uint64    `json:"counter"`
	GasLimit     uint64    `json:"gasLimit"`
	GasUsed      uint64    `json:"gasUsed"`
	StorageLimit uint64    `json:"storageLimit"`
	StorageUsed  uint64    `json:"storageUsed"`
	BakerFee     uint64    `json:"bakerFee"`
	StorageFee   uint64    `json:"storageFee"`
	ID           uint64    `json:"id"`
	Level        uint64    `json:"level"`
	Sender       Address   `json:"sender"`
	Rollup       Address   `json:"rollup"`
	Errors       []Error   `json:"errors,omitempty"`
	Quote        *Quote    `json:"quote,omitempty"`
}

// TxRollupFinalizeCommitment -
type TxRollupFinalizeCommitment struct {
	Type         string    `json:"type"`
	Block        string    `json:"block"`
	Hash         string    `json:"hash"`
	Status       string    `json:"status"`
	ID           uint64    `json:"id"`
	Level        uint64    `json:"level"`
	Counter      uint64    `json:"counter"`
	GasLimit     uint64    `json:"gasLimit"`
	GasUsed      uint64    `json:"gasUsed"`
	StorageLimit uint64    `json:"storageLimit"`
	StorageUsed  uint64    `json:"storageUsed"`
	BakerFee     uint64    `json:"bakerFee"`
	Timestamp    time.Time `json:"timestamp"`
	Sender       Address   `json:"sender"`
	Rollup       Address   `json:"rollup"`
	Errors       []Error   `json:"errors,omitempty"`
	Quote        *Quote    `json:"quote,omitempty"`
}

// TxRollupOrigination -
type TxRollupOrigination struct {
	ID            uint64    `json:"id"`
	Level         uint64    `json:"level"`
	Counter       uint64    `json:"counter"`
	GasLimit      uint64    `json:"gasLimit"`
	GasUsed       uint64    `json:"gasUsed"`
	StorageLimit  uint64    `json:"storageLimit"`
	StorageUsed   uint64    `json:"storageUsed"`
	BakerFee      uint64    `json:"bakerFee"`
	AllocationFee uint64    `json:"allocationFee"`
	Rollup        Address   `json:"rollup"`
	Sender        Address   `json:"sender"`
	Timestamp     time.Time `json:"timestamp"`
	Block         string    `json:"block"`
	Hash          string    `json:"hash"`
	Type          string    `json:"type"`
	Status        string    `json:"status"`
	Errors        []Error   `json:"errors,omitempty"`
	Quote         *Quote    `json:"quote,omitempty"`
}

// TxRollupRejection-
type TxRollupRejection struct {
	Type         string    `json:"type"`
	Block        string    `json:"block"`
	Hash         string    `json:"hash"`
	Status       string    `json:"status"`
	ID           uint64    `json:"id"`
	Level        uint64    `json:"level"`
	Counter      uint64    `json:"counter"`
	GasLimit     uint64    `json:"gasLimit"`
	GasUsed      uint64    `json:"gasUsed"`
	StorageLimit uint64    `json:"storageLimit"`
	StorageUsed  uint64    `json:"storageUsed"`
	BakerFee     uint64    `json:"bakerFee"`
	Reward       uint64    `json:"reward"`
	Loss         uint64    `json:"loss"`
	Rollup       Address   `json:"rollup"`
	Sender       Address   `json:"sender"`
	Committer    Address   `json:"committer"`
	Timestamp    time.Time `json:"timestamp"`
	Errors       []Error   `json:"errors,omitempty"`
	Quote        *Quote    `json:"quote,omitempty"`
}

// TxRollupRemoveCommitment -
type TxRollupRemoveCommitment struct {
	Type         string    `json:"type"`
	Status       string    `json:"status"`
	Block        string    `json:"block"`
	Hash         string    `json:"hash"`
	ID           uint64    `json:"id"`
	Level        uint64    `json:"level"`
	Counter      uint64    `json:"counter"`
	GasLimit     uint64    `json:"gasLimit"`
	GasUsed      uint64    `json:"gasUsed"`
	StorageLimit uint64    `json:"storageLimit"`
	StorageUsed  uint64    `json:"storageUsed"`
	BakerFee     uint64    `json:"bakerFee"`
	Sender       Address   `json:"sender"`
	Rollup       Address   `json:"rollup"`
	Timestamp    time.Time `json:"timestamp"`
	Errors       []Error   `json:"errors,omitempty"`
	Quote        *Quote    `json:"quote,omitempty"`
}

// TxRollupReturnBond -
type TxRollupReturnBond struct {
	Type         string    `json:"type"`
	Block        string    `json:"block"`
	Hash         string    `json:"hash"`
	Status       string    `json:"status"`
	ID           uint64    `json:"id"`
	Level        uint64    `json:"level"`
	Counter      uint64    `json:"counter"`
	GasLimit     uint64    `json:"gasLimit"`
	GasUsed      uint64    `json:"gasUsed"`
	StorageLimit uint64    `json:"storageLimit"`
	StorageUsed  uint64    `json:"storageUsed"`
	BakerFee     uint64    `json:"bakerFee"`
	Bond         uint64    `json:"bond"`
	Timestamp    time.Time `json:"timestamp"`
	Rollup       Address   `json:"rollup"`
	Sender       Address   `json:"sender"`
	Errors       []Error   `json:"errors,omitempty"`
	Quote        *Quote    `json:"quote,omitempty"`
}

// TxRollupSubmitBatch -
type TxRollupSubmitBatch struct {
	Type         string    `json:"type"`
	Block        string    `json:"block"`
	Hash         string    `json:"hash"`
	Status       string    `json:"status"`
	ID           uint64    `json:"id"`
	Level        uint64    `json:"level"`
	Counter      uint64    `json:"counter"`
	GasLimit     uint64    `json:"gasLimit"`
	GasUsed      uint64    `json:"gasUsed"`
	StorageLimit uint64    `json:"storageLimit"`
	StorageUsed  uint64    `json:"storageUsed"`
	BakerFee     uint64    `json:"bakerFee"`
	StorageFee   uint64    `json:"storageFee"`
	Rollup       Address   `json:"rollup"`
	Sender       Address   `json:"sender"`
	Timestamp    time.Time `json:"timestamp"`
	Errors       []Error   `json:"errors,omitempty"`
	Quote        *Quote    `json:"quote,omitempty"`
}

// SetDepositsLimit -
type SetDepositsLimit struct {
	ID           uint64    `json:"id"`
	Level        uint64    `json:"level"`
	Timestamp    time.Time `json:"timestamp"`
	Type         string    `json:"type"`
	Block        string    `json:"block"`
	Hash         string    `json:"hash"`
	Sender       Address   `json:"sender"`
	Counter      uint64    `json:"counter"`
	GasLimit     uint64    `json:"gasLimit"`
	GasUsed      uint64    `json:"gasUsed"`
	StorageLimit uint64    `json:"storageLimit"`
	BakerFee     uint64    `json:"bakerFee"`
	Status       string    `json:"status"`
	Limit        string    `json:"limit"`
	Errors       []Error   `json:"errors,omitempty"`
	Quote        *Quote    `json:"quote,omitempty"`
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

// DoubleBaking -
type DoubleBaking struct {
	Type         string    `json:"type"`
	ID           uint64    `json:"id"`
	Level        uint64    `json:"level"`
	Timestamp    time.Time `json:"timestamp"`
	Block        string    `json:"block"`
	Hash         string    `json:"hash"`
	AccusedLevel uint64    `json:"accusedLevel"`
	Accuser      *Address  `json:"accuser"`
	Offender     *Address  `json:"offender"`
	Quote        *Quote    `json:"quote,omitempty"`
}

// DoubleEndorsing -
type DoubleEndorsing struct {
	Type         string    `json:"type"`
	ID           uint64    `json:"id"`
	Level        uint64    `json:"level"`
	Timestamp    time.Time `json:"timestamp"`
	Block        string    `json:"block"`
	Hash         string    `json:"hash"`
	AccusedLevel uint64    `json:"accusedLevel"`
	Accuser      *Address  `json:"accuser"`
	Offender     *Address  `json:"offender"`
	Quote        *Quote    `json:"quote,omitempty"`
}

// Quote -
type Quote struct {
	BTC decimal.Decimal `json:"btc,omitempty"`
	EUR decimal.Decimal `json:"eur,omitempty"`
	USD decimal.Decimal `json:"usd,omitempty"`
	CNY decimal.Decimal `json:"cny,omitempty"`
	JPY decimal.Decimal `json:"jpy,omitempty"`
	KRW decimal.Decimal `json:"krw,omitempty"`
	ETH decimal.Decimal `json:"eth,omitempty"`
	GBP decimal.Decimal `json:"gbp,omitempty"`
}

// Error -
type Error struct {
	Type string `json:"type"`
}

// Baking -
type Baking struct {
	Type               string    `json:"type"`
	ID                 uint64    `json:"id"`
	Level              uint64    `json:"level"`
	Timestamp          time.Time `json:"timestamp"`
	Block              string    `json:"block"`
	Proposer           *Address  `json:"proposer"`
	Producer           *Address  `json:"producer"`
	PayloadRound       int       `json:"payloadRound"`
	BlockRound         int       `json:"blockRound"`
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
	Quote              *Quote    `json:"quote,omitempty"`
}

// EndorsingReward -
type EndorsingReward struct {
	Type      string    `json:"type"`
	ID        uint64    `json:"id"`
	Level     uint64    `json:"level"`
	Timestamp time.Time `json:"timestamp"`
	Block     string    `json:"block"`
	Baker     *Address  `json:"baker"`
	Expected  int64     `json:"expected"`
	Received  int64     `json:"received"`
	Quote     *Quote    `json:"quote,omitempty"`
}

// RevelationPenalty -
type RevelationPenalty struct {
	Type        string    `json:"type"`
	ID          uint64    `json:"id"`
	Level       uint64    `json:"level"`
	Timestamp   time.Time `json:"timestamp"`
	Block       string    `json:"block"`
	Baker       *Address  `json:"baker"`
	MissedLevel int64     `json:"missedLevel"`
	Loss        int64     `json:"loss"`
	Quote       *Quote    `json:"quote,omitempty"`
}

// DoublePreendorsing -
type DoublePreendorsing struct {
	Type                 string    `json:"type"`
	ID                   uint64    `json:"id"`
	Level                uint64    `json:"level"`
	Timestamp            time.Time `json:"timestamp"`
	Block                string    `json:"block"`
	Hash                 string    `json:"hash"`
	AccusedLevel         uint64    `json:"accusedLevel"`
	Accuser              *Address  `json:"accuser"`
	Offender             *Address  `json:"offender"`
	Quote                *Quote    `json:"quote,omitempty"`
	AccuserRewards       int64     `json:"accuserRewards,omitempty"`
	OffenderLostDeposits int64     `json:"offenderLostDeposits,omitempty"`
	OffenderLostRewards  int64     `json:"offenderLostRewards,omitempty"`
	OffenderLostFees     int64     `json:"offenderLostFees,omitempty"`
}

// VdfRevelation -
type VdfRevelation struct {
	Type      string    `json:"type"`
	ID        uint64    `json:"id"`
	Level     uint64    `json:"level"`
	Timestamp time.Time `json:"timestamp"`
	Block     string    `json:"block"`
	Hash      string    `json:"hash"`
	Baker     *Address  `json:"baker"`
	Cycle     uint64    `json:"cycle"`
	Solution  string    `json:"solution"`
	Proof     string    `json:"proof"`
	Reward    uint64    `json:"reward"`
	Quote     *Quote    `json:"quote,omitempty"`
}

// IncreasePaidStorage -
type IncreasePaidStorage struct {
	Type         string          `json:"type"`
	ID           uint64          `json:"id"`
	Level        uint64          `json:"level"`
	Timestamp    time.Time       `json:"timestamp"`
	Block        string          `json:"block"`
	Hash         string          `json:"hash"`
	Sender       Address         `json:"sender"`
	Counter      uint64          `json:"counter"`
	GasLimit     uint64          `json:"gasLimit"`
	GasUsed      uint64          `json:"gasUsed"`
	StorageLimit uint64          `json:"storageLimit"`
	StorageUsed  uint64          `json:"storageUsed"`
	BakerFee     uint64          `json:"bakerFee"`
	StorageFee   uint64          `json:"storageFee"`
	Status       string          `json:"status"`
	Contract     Address         `json:"contract"`
	Amount       decimal.Decimal `json:"amount"`
}

// UpdateConsensusKey -
type UpdateConsensusKey struct {
	Type            string    `json:"type"`
	ID              uint64    `json:"id"`
	Level           uint64    `json:"level"`
	Timestamp       time.Time `json:"timestamp"`
	Block           string    `json:"block"`
	Hash            string    `json:"hash"`
	Sender          Address   `json:"sender"`
	Counter         uint64    `json:"counter"`
	GasLimit        uint64    `json:"gasLimit"`
	GasUsed         uint64    `json:"gasUsed"`
	StorageLimit    uint64    `json:"storageLimit"`
	BakerFee        uint64    `json:"bakerFee"`
	Status          string    `json:"status"`
	ActivationCycle uint64    `json:"activationCycle"`
	PublicKey       string    `json:"publicKey"`
	PublicKeyHash   string    `json:"publicKeyHash"`
	Errors          []Error   `json:"errors,omitempty"`
	Quote           *Quote    `json:"quote,omitempty"`
}

// DrainDelegate -
type DrainDelegate struct {
	Type      string    `json:"type"`
	ID        uint64    `json:"id"`
	Level     uint64    `json:"level"`
	Timestamp time.Time `json:"timestamp"`
	Block     string    `json:"block"`
	Hash      string    `json:"hash"`
	Delegate  Address   `json:"delegate"`
	Target    Address   `json:"target"`
	Amount    uint64    `json:"amount"`
	Fee       uint64    `json:"fee"`
	Quote     *Quote    `json:"quote,omitempty"`
}

// SmartRollupAddMessage -
type SmartRollupAddMessage struct {
	Type          string    `json:"type"`
	ID            uint64    `json:"id"`
	Level         uint64    `json:"level"`
	Timestamp     time.Time `json:"timestamp"`
	Hash          string    `json:"hash"`
	Sender        *Address  `json:"sender,omitempty"`
	Counter       uint64    `json:"counter"`
	GasLimit      uint64    `json:"gasLimit"`
	GasUsed       uint64    `json:"gasUsed"`
	StorageLimit  uint64    `json:"storageLimit"`
	BakerFee      uint64    `json:"bakerFee"`
	Status        string    `json:"status"`
	MessagesCount uint64    `json:"messagesCount"`
	Errors        []Error   `json:"errors,omitempty"`
	Quote         *Quote    `json:"quote,omitempty"`
}

// SmartRollupCement -
type SmartRollupCement struct {
	Type         string            `json:"type"`
	ID           uint64            `json:"id"`
	Level        uint64            `json:"level"`
	Timestamp    time.Time         `json:"timestamp"`
	Hash         string            `json:"hash"`
	Sender       *Address          `json:"sender,omitempty"`
	Counter      uint64            `json:"counter"`
	GasLimit     uint64            `json:"gasLimit"`
	GasUsed      uint64            `json:"gasUsed"`
	StorageLimit uint64            `json:"storageLimit"`
	BakerFee     uint64            `json:"bakerFee"`
	Status       string            `json:"status"`
	Rollup       *Address          `json:"rollup,omitempty"`
	Commitment   *SrCommitmentInfo `json:"commitment,omitempty"`
	Errors       []Error           `json:"errors,omitempty"`
	Quote        *Quote            `json:"quote,omitempty"`
}

// SrCommitmentInfo -
type SrCommitmentInfo struct {
	ID         uint64   `json:"id"`
	Initiator  *Address `json:"initiator,omitempty"`
	InboxLevel uint64   `json:"inboxLevel"`
	State      string   `json:"state"`
	Hash       string   `json:"hash"`
	Ticks      uint64   `json:"ticks"`
	FirstLevel uint64   `json:"firstLevel"`
	FirstTime  string   `json:"firstTime"`
}

// SmartRollupExecute -
type SmartRollupExecute struct {
	Type         string            `json:"type"`
	ID           uint64            `json:"id"`
	Level        uint64            `json:"level"`
	Timestamp    time.Time         `json:"timestamp"`
	Hash         string            `json:"hash"`
	Sender       *Address          `json:"sender,omitempty"`
	Counter      uint64            `json:"counter"`
	GasLimit     uint64            `json:"gasLimit"`
	GasUsed      uint64            `json:"gasUsed"`
	StorageLimit uint64            `json:"storageLimit"`
	StorageUsed  uint64            `json:"storageUsed"`
	BakerFee     uint64            `json:"bakerFee"`
	StorageFee   uint64            `json:"storageFee"`
	Status       string            `json:"status"`
	Rollup       *Address          `json:"rollup,omitempty"`
	Commitment   *SrCommitmentInfo `json:"commitment,omitempty"`
	Errors       []Error           `json:"errors,omitempty"`
	Quote        *Quote            `json:"quote,omitempty"`
}

// SmartRollupOriginate -
type SmartRollupOriginate struct {
	Type          string             `json:"type"`
	ID            uint64             `json:"id"`
	Level         uint64             `json:"level"`
	Timestamp     time.Time          `json:"timestamp"`
	Hash          string             `json:"hash"`
	Sender        *Address           `json:"sender,omitempty"`
	Counter       uint64             `json:"counter"`
	GasLimit      uint64             `json:"gasLimit"`
	GasUsed       uint64             `json:"gasUsed"`
	StorageLimit  uint64             `json:"storageLimit"`
	StorageUsed   uint64             `json:"storageUsed"`
	BakerFee      uint64             `json:"bakerFee"`
	StorageFee    uint64             `json:"storageFee"`
	Status        string             `json:"status"`
	Rollup        *Address           `json:"rollup,omitempty"`
	ParameterType stdJSON.RawMessage `json:"parameterType,omitempty"`
	Errors        []Error            `json:"errors,omitempty"`
	Quote         *Quote             `json:"quote,omitempty"`
}

// SmartRollupPublish -
type SmartRollupPublish struct {
	Type         string            `json:"type"`
	ID           uint64            `json:"id"`
	Level        uint64            `json:"level"`
	Timestamp    time.Time         `json:"timestamp"`
	Hash         string            `json:"hash"`
	Sender       *Address          `json:"sender,omitempty"`
	Counter      uint64            `json:"counter"`
	GasLimit     uint64            `json:"gasLimit"`
	GasUsed      uint64            `json:"gasUsed"`
	StorageLimit uint64            `json:"storageLimit"`
	BakerFee     uint64            `json:"bakerFee"`
	Status       string            `json:"status"`
	Rollup       *Address          `json:"rollup,omitempty"`
	Commitment   *SrCommitmentInfo `json:"commitment,omitempty"`
	Bond         uint64            `json:"bond"`
}

// SmartRollupRecoverBond -
type SmartRollupRecoverBond struct {
	Type         string    `json:"type"`
	ID           uint64    `json:"id"`
	Level        uint64    `json:"level"`
	Timestamp    time.Time `json:"timestamp"`
	Hash         string    `json:"hash"`
	Sender       *Address  `json:"sender,omitempty"`
	Counter      uint64    `json:"counter"`
	GasLimit     uint64    `json:"gasLimit"`
	GasUsed      uint64    `json:"gasUsed"`
	StorageLimit uint64    `json:"storageLimit"`
	BakerFee     uint64    `json:"bakerFee"`
	Status       string    `json:"status"`
	Rollup       *Address  `json:"rollup,omitempty"`
	Staker       *Address  `json:"staker,omitempty"`
	Bond         uint64    `json:"bond"`
}

// SmartRollupRefute
type SmartRollupRefute struct {
	Type         string      `json:"type"`
	ID           uint64      `json:"id"`
	Level        uint64      `json:"level"`
	Timestamp    time.Time   `json:"timestamp"`
	Hash         string      `json:"hash"`
	Sender       *Address    `json:"sender,omitempty"`
	Counter      uint64      `json:"counter"`
	GasLimit     uint64      `json:"gasLimit"`
	GasUsed      uint64      `json:"gasUsed"`
	StorageLimit uint64      `json:"storageLimit"`
	BakerFee     uint64      `json:"bakerFee"`
	Status       string      `json:"status"`
	Rollup       *Address    `json:"rollup,omitempty"`
	Game         *SrGameInfo `json:"game"`
	Move         string      `json:"move"`
	GameStatus   string      `json:"gameStatus"`
}

// SrGameInfo -
type SrGameInfo struct {
	ID                  uint64            `json:"id"`
	Initiator           *Address          `json:"initiator,omitempty"`
	InitiatorCommitment *SrCommitmentInfo `json:"initiatorCommitment,omitempty"`
	Opponent            *Address          `json:"opponent,omitempty"`
	OpponentCommitment  *SrCommitmentInfo `json:"opponentCommitment,omitempty"`
	InitiatorReward     uint64            `json:"initiatorReward"`
	InitiatorLoss       uint64            `json:"initiatorLoss"`
	OpponentReward      uint64            `json:"opponentReward"`
	OpponentLoss        uint64            `json:"opponentLoss"`
}

type DalPublishCommitment struct {
	Type         string    `json:"type"`
	ID           uint64    `json:"id"`
	Level        uint64    `json:"level"`
	Timestamp    time.Time `json:"timestamp"`
	Hash         string    `json:"hash"`
	Sender       *Address  `json:"sender"`
	Counter      uint64    `json:"counter"`
	GasLimit     uint64    `json:"gasLimit"`
	GasUsed      uint64    `json:"gasUsed"`
	StorageLimit uint64    `json:"storageLimit"`
	BakerFee     uint64    `json:"bakerFee"`
	Slot         int       `json:"slot"`
	Commitment   string    `json:"commitment"`
	Status       string    `json:"status"`
}
