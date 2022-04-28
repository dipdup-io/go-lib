package types

import "time"

// BigMapDiff -
type BigMapDiff struct {
	Ptr   int64
	Key   []byte
	Value []byte

	ID          int64
	KeyHash     string
	OperationID int64
	Level       int64
	Address     string
	IndexedTime int64
	Timestamp   time.Time
	Protocol    int64
}
