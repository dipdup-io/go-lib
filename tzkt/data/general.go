package data

import (
	"encoding/json"
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
	RampUpCycles                      int64   `json:"rampUpCycles"`
	NoRewardCycles                    int64   `json:"noRewardCycles"`
	ConsensusRightsDelay              int64   `json:"consensusRightsDelay"`
	DelegateParametersActivationDelay int64   `json:"delegateParametersActivationDelay"`
	BlocksPerCycle                    int64   `json:"blocksPerCycle"`
	BlocksPerCommitment               int64   `json:"blocksPerCommitment"`
	BlocksPerSnapshot                 int64   `json:"blocksPerSnapshot"`
	BlocksPerVoting                   int64   `json:"blocksPerVoting"`
	TimeBetweenBlocks                 int64   `json:"timeBetweenBlocks"`
	AttestersPerBlock                 int64   `json:"attestersPerBlock"`
	HardOperationGasLimit             int64   `json:"hardOperationGasLimit"`
	HardOperationStorageLimit         int64   `json:"hardOperationStorageLimit"`
	HardBlockGasLimit                 int64   `json:"hardBlockGasLimit"`
	MinimalStake                      int64   `json:"minimalStake"`
	MinimalFrozenStake                int64   `json:"minimalFrozenStake"`
	BlockDeposit                      int64   `json:"blockDeposit"`
	BlockReward                       []int64 `json:"blockReward"`
	AttestationDeposit                int64   `json:"attestationDeposit"`
	AttestationReward                 []int64 `json:"attestationReward"`
	OriginationSize                   int64   `json:"originationSize"`
	ByteCost                          int64   `json:"byteCost"`
	ProposalQuorum                    int64   `json:"proposalQuorum"`
	BallotQuorumMin                   int64   `json:"ballotQuorumMin"`
	BallotQuorumMax                   int64   `json:"ballotQuorumMax"`
	LbToggleThreshold                 int64   `json:"lbToggleThreshold"`
	ConsensusThreshold                int64   `json:"consensusThreshold"`
	MinParticipationNumerator         int64   `json:"minParticipationNumerator"`
	MinParticipationDenominator       int64   `json:"minParticipationDenominator"`
	DenunciationPeriod                int64   `json:"denunciationPeriod"`
	SlashingDelay                     int64   `json:"slashingDelay"`
	MaxDelegatedOverFrozenRatio       int64   `json:"maxDelegatedOverFrozenRatio"`
	MaxExternalOverOwnStakeRatio      int64   `json:"maxExternalOverOwnStakeRatio"`
	SmartRollupOriginationSize        int64   `json:"smartRollupOriginationSize"`
	SmartRollupStakeAmount            int64   `json:"smartRollupStakeAmount"`
	SmartRollupChallengeWindow        int64   `json:"smartRollupChallengeWindow"`
	SmartRollupCommitmentPeriod       int64   `json:"smartRollupCommitmentPeriod"`
	SmartRollupTimeoutPeriod          int64   `json:"smartRollupTimeoutPeriod"`
	DalNumberOfShards                 int64   `json:"dalNumberOfShards"`
	Dictator                          string  `json:"dictator"`
}

// ProtocolMetadata -
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
	Quote             *Quote    `json:"quote,omitempty"`
}
