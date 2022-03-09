package node

import (
	stdJSON "encoding/json"
	"time"
)

// Header -
type Header struct {
	Protocol         string    `json:"protocol"`
	ChainID          string    `json:"chain_id"`
	Hash             string    `json:"hash"`
	Level            uint64    `json:"level"`
	Proto            int       `json:"proto"`
	Predecessor      string    `json:"predecessor"`
	Timestamp        time.Time `json:"timestamp"`
	ValidationPass   int       `json:"validation_pass"`
	OperationsHash   string    `json:"operations_hash"`
	Fitness          []string  `json:"fitness"`
	Context          string    `json:"context"`
	Priority         int       `json:"priority"`
	ProofOfWorkNonce string    `json:"proof_of_work_nonce"`
	Signature        string    `json:"signature"`
}

// BlockMetadata -
type BlockMetadata struct {
	Protocol        string `json:"protocol"`
	NextProtocol    string `json:"next_protocol"`
	TestChainStatus struct {
		Status string `json:"status"`
	} `json:"test_chain_status"`
	MaxOperationsTTL       int `json:"max_operations_ttl"`
	MaxOperationDataLength int `json:"max_operation_data_length"`
	MaxBlockHeaderLength   int `json:"max_block_header_length"`
	MaxOperationListLength []struct {
		MaxSize int `json:"max_size"`
		MaxOp   int `json:"max_op,omitempty"`
	} `json:"max_operation_list_length"`
	Baker            string    `json:"baker"`
	LevelInfo        LevelInfo `json:"level_info"`
	VotingPeriodInfo struct {
		VotingPeriod struct {
			Index         int    `json:"index"`
			Kind          string `json:"kind"`
			StartPosition int    `json:"start_position"`
		} `json:"voting_period"`
		Position  int `json:"position"`
		Remaining int `json:"remaining"`
	} `json:"voting_period_info"`
	NonceHash                 string                     `json:"nonce_hash"`
	ConsumedGas               string                     `json:"consumed_gas"`
	Deactivated               []interface{}              `json:"deactivated"`
	BalanceUpdates            []BalanceUpdate            `json:"balance_updates"`
	LiquidityBakingEscapeEma  int                        `json:"liquidity_baking_escape_ema"`
	ImplicitOperationsResults []ImplicitOperationsResult `json:"implicit_operations_results"`
}

// Block -
type Block struct {
	Protocol   string        `json:"protocol"`
	ChainID    string        `json:"chain_id"`
	Hash       string        `json:"hash"`
	Header     Header        `json:"header"`
	Metadata   BlockMetadata `json:"metadata"`
	Operations [][]OperationGroup
}

// BalanceUpdate -
type BalanceUpdate struct {
	Kind     string `json:"kind"`
	Contract string `json:"contract,omitempty"`
	Change   string `json:"change"`
	Category string `json:"category,omitempty"`
	Origin   string `json:"origin,omitempty"`
	Delegate string `json:"delegate,omitempty"`
	Cycle    uint64 `json:"cycle,omitempty"`
	Level    uint64 `json:"level,omitempty"`
}

// ImplicitOperationsResult -
type ImplicitOperationsResult struct {
	Kind                string             `json:"kind"`
	BalanceUpdates      []BalanceUpdate    `json:"balance_updates"`
	OriginatedContracts []string           `json:"originated_contracts,omitempty"`
	StorageSize         int64              `json:"storage_size,string"`
	PaidStorageSizeDiff int64              `json:"paid_storage_size_diff,string"`
	Storage             stdJSON.RawMessage `json:"storage,omitempty"`
	ConsumedGas         int64              `json:"consumed_gas,string,omitempty"`
	ConsumedMilligas    int64              `json:"consumed_milligas,string,omitempty"`
}

// LevelInfo -
type LevelInfo struct {
	Level              int64 `json:"level"`
	LevelPosition      int64 `json:"level_position"`
	Cycle              int64 `json:"cycle"`
	CyclePosition      int64 `json:"cycle_position"`
	ExpectedCommitment bool  `json:"expected_commitment"`
}

// ProtocolData -
type ProtocolData struct {
	Protocol                  string `json:"protocol"`
	Priority                  int    `json:"priority"`
	ProofOfWorkNonce          string `json:"proof_of_work_nonce"`
	LiquidityBakingEscapeVote bool   `json:"liquidity_baking_escape_vote"`
	Signature                 string `json:"signature"`
}

// HeaderShell -
type HeaderShell struct {
	Level          int64     `json:"level"`
	Proto          int       `json:"proto"`
	Predecessor    string    `json:"predecessor"`
	Timestamp      time.Time `json:"timestamp"`
	ValidationPass int       `json:"validation_pass"`
	OperationsHash string    `json:"operations_hash"`
	Fitness        []string  `json:"fitness"`
	Context        string    `json:"context"`
}

// BlockProtocols -
type BlockProtocols struct {
	Protocol     string `json:"protocol"`
	NextProtocol string `json:"next_protocol"`
}

// BlockBallot -
type BlockBallot struct {
	Pkh    string `json:"pkh"`
	Ballot string `json:"ballot"`
}

// BlockBallots -
type BlockBallots struct {
	Yay  int `json:"yay"`
	Nay  int `json:"nay"`
	Pass int `json:"pass"`
}

// VotingPeriod -
type VotingPeriod struct {
	VotingPeriod struct {
		Index         int    `json:"index"`
		Kind          string `json:"kind"`
		StartPosition int    `json:"start_position"`
	} `json:"voting_period"`
	Position  int `json:"position"`
	Remaining int `json:"remaining"`
}

// Rolls -
type Rolls struct {
	Pkh   string `json:"pkh"`
	Rolls int    `json:"rolls"`
}

// BlocksArgs -
type BlocksArgs struct {
	Length   uint
	HeadHash string
}
