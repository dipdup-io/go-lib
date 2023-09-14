package data

import "time"

type SmartRollup struct {
	Type                        string    `json:"type"`
	ID                          uint64    `json:"id"`
	Address                     string    `json:"address"`
	Alias                       string    `json:"alias"`
	Creator                     Address   `json:"creator"`
	PvmKind                     string    `json:"pvmKind"`
	GenesisCommitment           string    `json:"genesisCommitment"`
	LastCommitment              string    `json:"lastCommitment"`
	InboxLevel                  uint64    `json:"inboxLevel"`
	TotalStakers                uint64    `json:"totalStakers"`
	ActiveStakers               uint64    `json:"activeStakers"`
	ExecutedCommitments         uint64    `json:"executedCommitments"`
	CementedCommitments         uint64    `json:"cementedCommitments"`
	PendingCommitments          uint64    `json:"pendingCommitments"`
	RefutedCommitments          uint64    `json:"refutedCommitments"`
	OrphanCommitments           uint64    `json:"orphanCommitments"`
	SmartRollupBonds            uint64    `json:"smartRollupBonds"`
	ActiveTokensCount           uint64    `json:"activeTokensCount"`
	TokenBalancesCount          uint64    `json:"tokenBalancesCount"`
	TokenTransfersCount         uint64    `json:"tokenTransfersCount"`
	NumTransactions             uint64    `json:"numTransactions"`
	TransferTicketCount         uint64    `json:"transferTicketCount"`
	SmartRollupCementCount      uint64    `json:"smartRollupCementCount"`
	SmartRollupExecuteCount     uint64    `json:"smartRollupExecuteCount"`
	SmartRollupOriginateCount   uint64    `json:"smartRollupOriginateCount"`
	SmartRollupPublishCount     uint64    `json:"smartRollupPublishCount"`
	SmartRollupRecoverBondCount uint64    `json:"smartRollupRecoverBondCount"`
	SmartRollupRefuteCount      uint64    `json:"smartRollupRefuteCount"`
	RefutationGamesCount        uint64    `json:"refutationGamesCount"`
	ActiveRefutationGamesCount  uint64    `json:"activeRefutationGamesCount"`
	FirstActivity               uint64    `json:"firstActivity"`
	FirstActivityTime           time.Time `json:"firstActivityTime"`
	LastActivity                uint64    `json:"lastActivity"`
	LastActivityTime            time.Time `json:"lastActivityTime"`
}
