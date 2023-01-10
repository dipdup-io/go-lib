package data

import "time"

// Cycle -
type Cycle struct {
	Index           uint64    `json:"index"`
	FirstLevel      uint64    `json:"firstLevel"`
	StartTime       time.Time `json:"startTime"`
	LastLevel       uint64    `json:"lastLevel"`
	EndTime         time.Time `json:"endTime"`
	SnapshotIndex   uint64    `json:"snapshotIndex"`
	SnapshotLevel   uint64    `json:"snapshotLevel"`
	RandomSeed      string    `json:"randomSeed,omitempty"`
	TotalBakers     uint64    `json:"totalBakers"`
	TotalStaking    uint64    `json:"totalStaking"`
	TotalDelegators uint64    `json:"totalDelegators"`
	TotalDelegated  uint64    `json:"totalDelegated"`
	SelectedBakers  uint64    `json:"selectedBakers"`
	SelectedStake   uint64    `json:"selectedStake"`
	Quote           *Quote    `json:"quote,omitempty"`
	TotalRolls      uint64    `json:"totalRolls"`
}
