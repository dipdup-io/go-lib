package events

import (
	stdJSON "encoding/json"
	"fmt"
)

// MessageType - TzKT message type
type MessageType int

// message types
const (
	MessageTypeState MessageType = iota
	MessageTypeData
	MessageTypeReorg
	MessageTypeSubscribed
)

// Message - message struct
type Message struct {
	Channel string
	Type    MessageType `json:"type"`
	State   uint64      `json:"state"`
	Body    interface{} `json:"data"`
}

// String -
func (msg Message) String() string {
	s := fmt.Sprintf("channel=%s type=%d state=%d", msg.Channel, msg.Type, msg.State)
	if msg.Body != nil {
		s = fmt.Sprintf("%s data=%v", s, msg.Body)
	}
	return s
}

// Packet -
type Packet struct {
	Type  MessageType        `json:"type"`
	State uint64             `json:"state"`
	Data  stdJSON.RawMessage `json:"data,omitempty"`
}
