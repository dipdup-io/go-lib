package node

import "github.com/pkg/errors"

// Connection -
type Connection struct {
	PeerID           string             `json:"peer_id"`
	Point            Point              `json:"id_point"`
	RemoteSocketPort int                `json:"remote_socket_port"`
	AnnouncedVersion ConnectionVersion  `json:"announced_version"`
	LocalMetadata    ConnectionMetadata `json:"local_metadata"`
	RemoteMetadata   ConnectionMetadata `json:"remote_metadata"`
	Incoming         bool               `json:"incoming"`
	Private          bool               `json:"private"`
}

// Point -
type Point struct {
	Addr string `json:"addr"`
	Port int    `json:"port"`
}

// ConnectionMetadata -
type ConnectionMetadata struct {
	DisableMempool bool `json:"disable_mempool"`
	PrivateNode    bool `json:"private_node"`
}

// ConnectionVersion -
type ConnectionVersion struct {
	ChainName            string `json:"chain_name"`
	DistributedDbVersion int    `json:"distributed_db_version"`
	P2PVersion           int    `json:"p2p_version"`
}

// NetworkPointWithURI -
type NetworkPointWithURI struct {
	URI string
	NetworkPoint
}

// UnmarshalJSON -
func (n *NetworkPointWithURI) UnmarshalJSON(buf []byte) error {
	tmp := []interface{}{&n.URI, &n.NetworkPoint}
	wantLen := len(tmp)
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	if g, e := len(tmp), wantLen; g != e {
		return errors.Errorf("wrong number of fields in NetworkPointWithURI: %d != %d", g, e)
	}
	return nil
}

// NetworkPoint -
type NetworkPoint struct {
	Trusted         bool   `json:"trusted"`
	GreylistedUntil string `json:"greylisted_until"`
	State           struct {
		EventKind string `json:"event_kind"`
		P2PPeerID string `json:"p2p_peer_id"`
	} `json:"state"`
	P2PPeerID                 string   `json:"p2p_peer_id"`
	LastFailedConnection      string   `json:"last_failed_connection"`
	LastRejectedConnection    []string `json:"last_rejected_connection"`
	LastEstablishedConnection []string `json:"last_established_connection"`
	LastDisconnection         []string `json:"last_disconnection"`
	LastSeen                  []string `json:"last_seen"`
	LastMiss                  string   `json:"last_miss"`
}
