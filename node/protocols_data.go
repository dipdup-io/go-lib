package node

// ProtocolInfo -
type ProtocolInfo struct {
	ExpectedEnvVersion int                 `json:"expected_env_version"`
	Components         []ProtocolComponent `json:"components"`
}

// ProtocolComponent -
type ProtocolComponent struct {
	Name           string `json:"name"`
	Interface      string `json:"interface,omitempty"`
	Implementation string `json:"implementation"`
}
