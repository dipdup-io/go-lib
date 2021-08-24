package signalr

import "errors"

var (
	ErrUnknownMessageType = errors.New("unknown message type")
	ErrMessageParsing     = errors.New("can't parse message")
	ErrEmptyResponse      = errors.New("empty response from server")
	ErrHandshake          = errors.New("handshake error")
	ErrInvalidStatusCode  = errors.New("invalid status code")
	ErrNegotiate          = errors.New("negotiate error")
	ErrInvalidScheme      = errors.New("invalid URL scheme. Expected https or http. Got")
	ErrConnectionClose    = errors.New("connection is closed")
	ErrTimeout            = errors.New("connection timeout")
)
