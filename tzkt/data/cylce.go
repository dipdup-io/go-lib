package data

import "time"

// Cycle -
type Cycle struct {
	Index                        uint64    `json:"index"`
	FirstLevel                   uint64    `json:"firstLevel"`
	StartTime                    time.Time `json:"startTime"`
	LastLevel                    uint64    `json:"lastLevel"`
	EndTime                      time.Time `json:"endTime"`
	SnapshotLevel                uint64    `json:"snapshotLevel"`
	RandomSeed                   string    `json:"randomSeed,omitempty"`
	TotalBakers                  uint64    `json:"totalBakers"`
	TotalBakingPower             uint64    `json:"totalBakingPower,omitempty"`
	BlockReward                  uint64    `json:"blockReward,omitempty"`
	BlockBonusPerSlot            uint64    `json:"blockBonusPerSlot,omitempty"`
	AttestationRewardPerSlot     uint64    `json:"attestationRewardPerSlot,omitempty"`
	NonceRevelationReward        uint64    `json:"nonceRevelationReward,omitempty"`
	VdfRevelationReward          uint64    `json:"vdfRevelationReward,omitempty"`
	DalAttestationRewardPerShard uint64    `json:"dalAttestationRewardPerShard,omitempty"`
	Quote                        *Quote    `json:"quote,omitempty"`
}
