package node

const (
	LazyStorageDiffKindBigMap  = "big_map"
	LazyStorageDiffKindSapling = "sapling_state"
)

const (
	KindActivation             = "activate_account"
	KindBallot                 = "ballot"
	KindDelegation             = "delegation"
	KindDoubleBaking           = "double_baking_evidence"
	KindDoubleEndorsing        = "double_endorsement_evidence"
	KindEndorsement            = "endorsement"
	KindEndorsementWithSlot    = "endorsement_with_slot"
	KindOrigination            = "origination"
	KindProposal               = "proposals"
	KindReveal                 = "reveal"
	KindNonceRevelation        = "seed_nonce_revelation"
	KindTransaction            = "transaction"
	KindRegisterGlobalConstant = "register_global_constant"
	KindPreendorsement         = "preendorsement"
	KindSetDepositsLimit       = "set_deposits_limit"
	KindDoublePreendorsement   = "double_preendorsement_evidence"
)

const (
	BigMapActionUpdate = "update"
	BigMapActionRemove = "remove"
	BigMapActionAlloc  = "alloc"
	BigMapActionCopy   = "copy"
)

const (
	BalanceUpdatesKindContract = "contract"
	BalanceUpdatesKindFreezer  = "freezer"
)

const (
	BalanceUpdatesCategoryReward   = "rewards"
	BalanceUpdatesCategoryFees     = "fees"
	BalanceUpdatesCategoryDeposits = "deposits"
)

const (
	BalanceUpdatesOriginBlock     = "block"
	BalanceUpdatesOriginMigration = "migration"
	BalanceUpdatesOriginSubsidy   = "subsidy"
)

const (
	HeadBlock = "head"
)
