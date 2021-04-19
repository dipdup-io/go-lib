package node

import "time"

// NodeOption -
type NodeOption func(*NodeRPC)

// WithTimeout -
func WithTimeout(seconds uint64) NodeOption {
	return func(m *NodeRPC) {
		if seconds > 0 {
			m.timeout = time.Duration(seconds) * time.Second
		} else {
			m.timeout = 10 * time.Second
		}
	}
}
