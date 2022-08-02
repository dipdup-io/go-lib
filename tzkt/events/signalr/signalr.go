package signalr

import (
	"context"

	"github.com/rs/zerolog"
)

// SignalR -
type SignalR struct {
	hub *Hub
	t   *Transport
	log zerolog.Logger
	url string
}

// NewSignalR -
func NewSignalR(url string) *SignalR {
	return &SignalR{
		t:   NewTransport(url),
		url: url,
	}
}

// SetLogger -
func (s *SignalR) SetLogger(log zerolog.Logger) {
	s.t.log = log
	s.log = log
}

// Connect - connect to server
func (s *SignalR) Connect(ctx context.Context, version Version) error {
	resp, err := s.t.Negotiate(version)
	if err != nil {
		return err
	}
	var id string
	switch version {
	case Version0:
		id = resp.ConnectionID
	case Version1:
		id = resp.ConnectionToken
	}

	hub, err := NewHub(s.url, id)
	if err != nil {
		return err
	}
	s.hub = hub
	s.hub.log = s.log

	return s.hub.Connect(ctx)
}

// Messages - listens message channel
func (s *SignalR) Messages() <-chan interface{} {
	return s.hub.msgs
}

// Close - close connection
func (s *SignalR) Close() error {
	return s.hub.Close()
}

// Send - send message to server
func (s *SignalR) Send(msg interface{}) error {
	return s.hub.Send(msg)
}

// SetOnReconnect -
func (s *SignalR) SetOnReconnect(onReconnect func() error) {
	s.hub.onReconnect = onReconnect
}

// IsConnected -
func (s *SignalR) IsConnected() bool {
	return s.hub != nil && s.hub.conn != nil
}
