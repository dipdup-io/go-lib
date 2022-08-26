package data

import (
	"encoding/json"
	stdJSON "encoding/json"
	"time"
)

// Address -
type Address struct {
	Alias   string `json:"alias"`
	Address string `json:"address"`
}

// OriginatedContract -
type OriginatedContract struct {
	Kind     string `json:"kind"`
	Alias    string `json:"alias,omitempty"`
	Address  string `json:"address,omitempty"`
	TypeHash int    `json:"typeHash"`
	CodeHash int    `json:"codeHash"`
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
	ID             uint64             `json:"id"`
	Contract       Address            `json:"contract"`
	TokenID        string             `json:"tokenId"`
	Standard       string             `json:"standard"`
	Metadata       stdJSON.RawMessage `json:"metadata,omitempty"`
	FirstLevel     uint64             `json:"firstLevel"`
	FirstTime      time.Time          `json:"firstTime"`
	LastLevel      uint64             `json:"lastLevel"`
	LastTime       time.Time          `json:"lastTime"`
	TransfersCount uint64             `json:"transfersCount"`
	BalancesCount  uint64             `json:"balancesCount"`
	HoldersCount   uint64             `json:"holdersCount"`
	TotalMinted    string             `json:"totalMinted"`
	TotalBurned    string             `json:"totalBurned"`
	TotalSupply    string             `json:"totalSupply"`
}

// Protocol -
type Protocol struct {
	Code            int64              `json:"code"`
	Hash            string             `json:"hash"`
	FirstLevel      uint64             `json:"firstLevel"`
	FirstCycle      uint64             `json:"firstCycle"`
	FirstCycleLevel uint64             `json:"firstCycleLevel"`
	LastLevel       uint64             `json:"lastLevel,omitempty"`
	Constants       *ProtocolConstants `json:"constants,omitempty"`
	Metadata        *ProtocolMetadata  `json:"metadata,omitempty"`
}

// ProtocolConstants -
type ProtocolConstants struct {
	RampUpCycles                         int64   `json:"rampUpCycles"`
	NoRewardCycles                       int64   `json:"noRewardCycles"`
	PreservedCycles                      int64   `json:"preservedCycles"`
	BlocksPerCycle                       int64   `json:"blocksPerCycle"`
	BlocksPerCommitment                  int64   `json:"blocksPerCommitment"`
	BlocksPerSnapshot                    int64   `json:"blocksPerSnapshot"`
	BlocksPerVoting                      int64   `json:"blocksPerVoting"`
	TimeBetweenBlocks                    int64   `json:"timeBetweenBlocks"`
	EndorsersPerBlock                    int64   `json:"endorsersPerBlock"`
	HardOperationGasLimit                int64   `json:"hardOperationGasLimit"`
	HardOperationStorageLimit            int64   `json:"hardOperationStorageLimit"`
	HardBlockGasLimit                    int64   `json:"hardBlockGasLimit"`
	TokensPerRoll                        int64   `json:"tokensPerRoll"`
	RevelationReward                     int64   `json:"revelationReward"`
	BlockDeposit                         int64   `json:"blockDeposit"`
	BlockReward                          []int64 `json:"blockReward"`
	EndorsementDeposit                   int64   `json:"endorsementDeposit"`
	EndorsementReward                    []int64 `json:"endorsementReward"`
	OriginationSize                      int64   `json:"originationSize"`
	ByteCost                             int64   `json:"byteCost"`
	ProposalQuorum                       int64   `json:"proposalQuorum"`
	BallotQuorumMin                      int64   `json:"ballotQuorumMin"`
	BallotQuorumMax                      int64   `json:"ballotQuorumMax"`
	LbSubsidy                            int64   `json:"lbSubsidy"`
	LbSunsetLevel                        int64   `json:"lbSunsetLevel"`
	LbToggleThreshold                    int64   `json:"lbToggleThreshold"`
	ConsensusThreshold                   int64   `json:"consensusThreshold"`
	MinParticipationNumerator            int64   `json:"minParticipationNumerator"`
	MinParticipationDenominator          int64   `json:"minParticipationDenominator"`
	MaxSlashingPeriod                    int64   `json:"maxSlashingPeriod"`
	FrozenDepositsPercentage             int64   `json:"frozenDepositsPercentage"`
	DoubleBakingPunishment               int64   `json:"doubleBakingPunishment"`
	DoubleEndorsingPunishmentNumerator   int64   `json:"doubleEndorsingPunishmentNumerator"`
	DoubleEndorsingPunishmentDenominator int64   `json:"doubleEndorsingPunishmentDenominator"`
	TxRollupOriginationSize              int64   `json:"txRollupOriginationSize"`
	TxRollupCommitmentBond               int64   `json:"txRollupCommitmentBond"`
	LbEscapeThreshold                    int64   `json:"lbEscapeThreshold"`
}

// ProtocolConstants -
type ProtocolMetadata struct {
	Docs  string `json:"docs"`
	Alias string `json:"alias"`
}

// Statistics -
type Statistics struct {
	Level             uint64    `json:"level"`
	Timestamp         time.Time `json:"timestamp"`
	TotalSupply       uint64    `json:"totalSupply"`
	CirculatingSupply uint64    `json:"circulatingSupply"`
	TotalBootstrapped uint64    `json:"totalBootstrapped"`
	TotalCommitments  uint64    `json:"totalCommitments"`
	TotalActivated    uint64    `json:"totalActivated"`
	TotalCreated      uint64    `json:"totalCreated"`
	TotalBurned       uint64    `json:"totalBurned"`
	TotalBanished     uint64    `json:"totalBanished"`
	TotalFrozen       uint64    `json:"totalFrozen"`
	TotalRollupBonds  uint64    `json:"totalRollupBonds"`
	TotalVested       uint64    `json:"totalVested"`
	Quote             *Quote    `json:"quote,omitempty"`
}
