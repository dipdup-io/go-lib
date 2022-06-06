package node

import (
	stdJSON "encoding/json"
	"strconv"
)

// Constants -
type Constants struct {
	ProofOfWorkNonceSize         int64            `json:"proof_of_work_nonce_size"`
	NonceLength                  int64            `json:"nonce_length"`
	MaxAnonOpsPerBlock           int64            `json:"max_anon_ops_per_block"`
	MaxOperationDataLength       int64            `json:"max_operation_data_length"`
	MaxProposalsPerDelegate      int64            `json:"max_proposals_per_delegate"`
	PreservedCycles              uint64           `json:"preserved_cycles"`
	BlocksPerCycle               uint64           `json:"blocks_per_cycle"`
	BlocksPerCommitment          int64            `json:"blocks_per_commitment"`
	BlocksPerRollSnapshot        int64            `json:"blocks_per_roll_snapshot"`
	BlocksPerVotingPeriod        int64            `json:"blocks_per_voting_period"`
	TimeBetweenBlocks            Int64StringSlice `json:"time_between_blocks"`
	EndorsersPerBlock            int64            `json:"endorsers_per_block"`
	HardGasLimitPerOperation     int64            `json:"hard_gas_limit_per_operation,string"`
	HardGasLimitPerBlock         int64            `json:"hard_gas_limit_per_block,string"`
	ProofOfWorkThreshold         int64            `json:"proof_of_work_threshold,string"`
	TokensPerRoll                int64            `json:"tokens_per_roll,string"`
	MichelsonMaximumTypeSize     int64            `json:"michelson_maximum_type_size"`
	SeedNonceRevelationTip       int64            `json:"seed_nonce_revelation_tip,string"`
	OriginationSize              int64            `json:"origination_size"`
	BlockSecurityDeposit         int64            `json:"block_security_deposit,string"`
	EndorsementSecurityDeposit   int64            `json:"endorsement_security_deposit,string"`
	BakingRewardPerEndorsement   Int64StringSlice `json:"baking_reward_per_endorsement"`
	EndorsementReward            Int64StringSlice `json:"endorsement_reward"`
	CostPerByte                  int64            `json:"cost_per_byte,string"`
	HardStorageLimitPerOperation int64            `json:"hard_storage_limit_per_operation,string"`
	TestChainDuration            int64            `json:"test_chain_duration,string"`
	QuorumMin                    int64            `json:"quorum_min"`
	QuorumMax                    int64            `json:"quorum_max"`
	MinProposalQuorum            int64            `json:"min_proposal_quorum"`
	InitialEndorsers             int64            `json:"initial_endorsers"`
	DelayPerMissingEndorsement   int64            `json:"delay_per_missing_endorsement,string"`
	MinimalBlockDelay            int64            `json:"minimal_block_delay,string,omitempty"`
}

// Int64StringSlice -
type Int64StringSlice []int64

// UnmarshalJSON -
func (slice *Int64StringSlice) UnmarshalJSON(data []byte) error {
	s := make([]string, 0)
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*slice = make([]int64, len(s))
	for i := range s {
		value, err := strconv.ParseInt(s[i], 10, 64)
		if err != nil {
			return err
		}
		(*slice)[i] = value
	}
	return nil
}

// ContractInfo -
type ContractInfo struct {
	Balance   string             `json:"balance"`
	Delegate  string             `json:"delegate,omitempty"`
	Script    stdJSON.RawMessage `json:"script,omitempty"`
	Counter   string             `json:"counter,omitempty"`
	Spendable bool               `json:"spendable,omitempty"`
	Manager   string             `json:"manager,omitempty"`
}

// Entrypoints -
type Entrypoints struct {
	Entrypoints map[string]stdJSON.RawMessage `json:"entrypoints"`
}

// Script -
type Script struct {
	Code   stdJSON.RawMessage `json:"code"`
	Strage stdJSON.RawMessage `json:"storage"`
}

// Delegate -
type Delegate struct {
	Balance              string               `json:"balance"`
	FrozenBalance        string               `json:"frozen_balance"`
	FrozenBalanceByCycle FrozenBalanceByCycle `json:"frozen_balance_by_cycle"`
	StakingBalance       string               `json:"staking_balance"`
	DelegatedContracts   []string             `json:"delegated_contracts"`
	DelegatedBalance     string               `json:"delegated_balance"`
	Deactivated          bool                 `json:"deactivated"`
	GracePeriod          int                  `json:"grace_period"`
	VotingPower          int                  `json:"voting_power"`
}

// FrozenBalanceByCycle -
type FrozenBalanceByCycle struct {
	Cycle    int    `json:"cycle"`
	Deposits string `json:"deposits"`
	Fees     string `json:"fees"`
	Rewards  string `json:"rewards"`
}

// DelegateType -
type DelegateType string

// delegate types
const (
	ActiveDelegateType   DelegateType = "active"
	InactiveDelegateType DelegateType = "inactive"
	AllDelegateType      DelegateType = "all"
)

// TxRollupState -
type TxRollupState struct {
	LastRemovedCommitmentHashes *LastRemovedCommitmentHashes `json:"last_removed_commitment_hashes,omitempty"`
	FinalizedCommitments        RollupStateCommitment        `json:"finalized_commitments"`
	UnfinalizedCommitments      RollupStateCommitment        `json:"unfinalized_commitments"`
	UncommittedInboxes          RollupStateCommitment        `json:"uncommitted_inboxes"`
	CommitmentNewestHash        *string                      `json:"commitment_newest_hash,omitempty"`
	TezosHeadLevel              *uint64                      `json:"tezos_head_level,omitempty"`
	BurnPerByte                 stdJSON.Number               `json:"burn_per_byte,omitempty"`
	AllocatedStorage            stdJSON.Number               `json:"allocated_storage"`
	OccupiedStorage             stdJSON.Number               `json:"occupied_storage"`
	InboxEma                    uint64                       `json:"inbox_ema"`
	CommitmentsWatermark        *uint64                      `json:"commitments_watermark,omitempty"`
}

// LastRemovedCommitmentHashes -
type LastRemovedCommitmentHashes struct {
	LastMessageHash string `json:"last_message_hash"`
	CommitmentHash  string `json:"commitment_hash"`
}

// RollupStateCommitment -
type RollupStateCommitment struct {
	Next   *uint64 `json:"next,omitempty"`
	Newest *uint64 `json:"newest,omitempty"`
	Oldest *uint64 `json:"oldest,omitempty"`
}

// RollupCommitmentForBlock -
type RollupCommitmentForBlock struct {
	Commitment     RollupCommitment `json:"commitment"`
	CommitmentHash string           `json:"commitment_hash"`
	Committer      string           `json:"committer"`
	SubmittedAt    uint64           `json:"submitted_at"`
	FinalizedAt    uint64           `json:"finalized_at"`
}

// RollupCommitment -
type RollupCommitment struct {
	Level           uint64                   `json:"level"`
	Messages        RollupCommitmentMessages `json:"messages"`
	Predecessor     *string                  `json:"predecessor,omitempty"`
	InboxMerkleRoot string                   `json:"inbox_merkle_root"`
}

// RollupCommitmentMessages -
type RollupCommitmentMessages struct {
	Count                 uint64 `json:"count"`
	Root                  string `json:"root"`
	LastMessageResultHash string `json:"last_message_result_hash"`
}

// TxRollupInbox -
type TxRollupInbox struct {
	InboxLength   uint64 `json:"inbox_length"`
	CumulatedSize uint64 `json:"cumulated_size"`
	MerkleRoot    string `json:"merkle_root"`
}
