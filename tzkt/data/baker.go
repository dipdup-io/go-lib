package data

import "time"

// Right -
type Right struct {
	Type      string    `json:"type"`
	Cycle     uint64    `json:"cycle"`
	Level     uint64    `json:"level"`
	Timestamp time.Time `json:"timestamp"`
	Slots     uint64    `json:"slots"`
	Baker     Address   `json:"baker"`
	Status    string    `json:"status"`
}
