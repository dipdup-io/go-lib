package signalr

import "errors"

var (
	ErrUnknownMessageType = errors.New("Unknown message type")
	ErrMessageParsing     = errors.New("Can't parse message")
	ErrEmptyResponse      = errors.New("Empty response from server")
	ErrHandshake          = errors.New("Handshake error")
	ErrInvalidStatusCode  = errors.New("Invalid status code")
	ErrNegotiate          = errors.New("Negotiate error")
	ErrInvalidScheme      = errors.New("Invalid URL scheme. Expected https or http. Got")
	ErrConnectionClose    = errors.New("Connection is closed")
	ErrTimeout            = errors.New("Connection timeout")
)
