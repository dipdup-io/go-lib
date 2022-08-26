package data

// kinds
const (
	KindTransaction              = "transaction"
	KindOrigination              = "origination"
	KindEndorsement              = "endorsement"
	KindPreendorsement           = "preendorsement"
	KindBallot                   = "ballot"
	KindProposal                 = "proposal"
	KindDoubleBaking             = "double_baking"
	KindDoubleEndorsing          = "double_endorsing"
	KindDoublePreendorsing       = "double_preendorsing"
	KindActivation               = "activation"
	KindMigration                = "migration"
	KindNonceRevelation          = "nonce_revelation"
	KindDelegation               = "delegation"
	KindReveal                   = "reveal"
	KindRegisterGlobalConstant   = "register_constant"
	KindTransferTicket           = "transfer_ticket"
	KindTxRollupCommit           = "tx_rollup_commit"
	KindRollupDispatchTickets    = "tx_rollup_dispatch_tickets"
	KindRollupFinalizeCommitment = "tx_rollup_finalize_commitment"
	KindTxRollupOrigination      = "tx_rollup_origination"
	KindTxRollupRejection        = "tx_rollup_rejection"
	KindTxRollupRemoveCommitment = "tx_rollup_remove_commitment"
	KindRollupReturnBond         = "tx_rollup_return_bond"
	KindRollupSubmitBatch        = "tx_rollup_submit_batch"
	KindSetDepositsLimit         = "set_deposits_limit"
	KindRevelationPenalty        = "revelation_penalty"
	KindBaking                   = "baking"
	KindEndorsingReward          = "endorsing_reward"
	KindVdfRevelation            = "vdf_revelation"
	KindIncreasePaidStorage      = "increase_paid_storage"
)

// urls
const (
	BaseURL       = "https://api.tzkt.io"
	BaseEventsURL = "https://api.tzkt.io/v1/events"
)
