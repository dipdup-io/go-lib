package signalr

import (
	stdJSON "encoding/json"
)

// MessageType -
type MessageType int

const (
	MessageTypeInvocation MessageType = iota + 1
	MessageTypeStreamItem
	MessageTypeCompletion
	MessageTypeStreamInvocation
	MessageTypeCancelInvocation
	MessageTypePing
	MessageTypeCloseMessage
)

// HandshakeRequest -
type HandshakeRequest struct {
	Protocol string `json:"protocol"`
	Version  int    `json:"version"`
}

func newHandshakeRequest() HandshakeRequest {
	return HandshakeRequest{
		Protocol: "json",
		Version:  1,
	}
}

// Error -
type Error struct {
	Error string `json:"string,omitempty"`
}

// Type -
type Type struct {
	Type MessageType `json:"type"`
}

// Message -
type Message struct {
	Type
	ID      string            `json:"invocationId"`
	Headers map[string]string `json:"headers,omitempty"`
}

// Invocation -  a `Invocation` message
type Invocation struct {
	Message
	Target    string               `json:"target"`
	Arguments []stdJSON.RawMessage `json:"arguments"`
	StreamsID []string             `json:"streamIds,omitempty"`
}

// StreamInvocation - a `StreamInvocation` message
type StreamInvocation Invocation

// StreamItem -  a `StreamItem` message
type StreamItem struct {
	Message
	Item string `json:"item"`
}

// Completion -  a `Completion` message
type Completion struct {
	Message
	Result uint64 `json:"result,omitempty"`
	Error  string `json:"string,omitempty"`
}

// CancelInvocation -  a `CancelInvocation` message
type CancelInvocation Message

// PingMessage -  a `PingMessage` message
type PingMessage Type

// CloseMessage -  a `CloseMessage` message
type CloseMessage struct {
	Type
	Error          string `json:"string,omitempty"`
	AllowReconnect bool   `json:"allowReconnect,omitempty"`
}

func newCloseMessage() CloseMessage {
	return CloseMessage{
		Type: Type{
			Type: MessageTypeCloseMessage,
		},
	}
}

// NegotiateResponse -
type NegotiateResponse struct {
	ConnectionToken     string               `json:"connectionToken"`
	ConnectionID        string               `json:"connectionId"`
	NegotiateVersion    int                  `json:"negotiateVersion"`
	AvailableTransports []AvailableTransport `json:"availableTransports"`
}

// AvailableTransport -
type AvailableTransport struct {
	Transport       string   `json:"transport"`
	TransferFormats []string `json:"transferFormats"`
}

// RedirectResponse -
type RedirectResponse struct {
	URL         string `json:"url"`
	AccessToken string `json:"accessToken"`
}

// NewInvocation -
func NewInvocation(id, target string, args ...interface{}) Invocation {
	if args == nil {
		args = make([]interface{}, 0)
	}

	arguments := make([]stdJSON.RawMessage, len(args))
	for i := range args {
		data, err := json.Marshal(args[i])
		if err == nil {
			arguments[i] = data
		}
	}

	return Invocation{
		Message: Message{
			Type: Type{MessageTypeInvocation},
			ID:   id,
		},
		Target:    target,
		Arguments: arguments,
	}
}
