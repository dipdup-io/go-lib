package events

import "fmt"

// MessageType - TzKT message type
type MessageType int

// message types
const (
	MessageTypeState MessageType = iota
	MessageTypeData
	MessageTypeReorg
)

// Message - message struct
type Message struct {
	Channel string
	Type    MessageType
	State   uint64
	Body    interface{}
}

// String -
func (msg Message) String() string {
	s := fmt.Sprintf("channel=%s type=%d state=%d", msg.Channel, msg.Type, msg.State)
	if msg.Body != nil {
		s = fmt.Sprintf("%s data=%v", s, msg.Body)
	}
	return s
}
