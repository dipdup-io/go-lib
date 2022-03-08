package node

// InvalidBlock -
type InvalidBlock struct {
	Hash   string        `json:"block"`
	Level  uint64        `json:"level"`
	Errors []interface{} `json:"errors"`
}

// Bootstrapped -
type Bootstrapped struct {
	Bootstrapped bool   `json:"bootstrapped"`
	SyncState    string `json:"sync_state"`
}

// Caboose -
type Caboose struct {
	Hash  string `json:"block_hash"`
	Level uint64 `json:"level"`
}

// Checkpoint -
type Checkpoint struct {
	Hash  string `json:"block_hash"`
	Level uint64 `json:"level"`
}

// Savepoint -
type Savepoint struct {
	Hash  string `json:"block_hash"`
	Level uint64 `json:"level"`
}
