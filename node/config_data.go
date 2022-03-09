package node

// HistoryMode -
type HistoryMode struct {
	Mode string `json:"history_mode"`
}

// ActivatedProtocol -
type ActivatedProtocol struct {
	ReplacedProtocol    string `json:"replaced_protocol"`
	ReplacementProtocol string `json:"replacement_protocol"`
}

// ActivatedUpgrades -
type ActivatedUpgrades struct {
	Level               int    `json:"level"`
	ReplacementProtocol string `json:"replacement_protocol"`
}
