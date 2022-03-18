package node

import (
	"context"
	"fmt"
)

// ProtocolsAPI -
type ProtocolsAPI interface {
	GetProtocols(ctx context.Context) ([]string, error)
	Protocol(ctx context.Context, hash string) (ProtocolInfo, error)
	Environment(ctx context.Context, hash string) (int, error)
}

// Protocols -
type Protocols struct {
	baseURL string
	client  *client
}

// NewProtocols -
func NewProtocols(baseURL string) *Protocols {
	return &Protocols{
		baseURL: baseURL,
		client:  newClient(),
	}
}

// GetProtocols -
func (api *Protocols) GetProtocols(ctx context.Context) ([]string, error) {
	req, err := newGetRequest(api.baseURL, "protocols", nil)
	if err != nil {
		return nil, err
	}
	var result []string
	err = req.doWithJSONResponse(ctx, api.client, &result)
	return result, err
}

// Protocol -
func (api *Protocols) Protocol(ctx context.Context, hash string) (ProtocolInfo, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("protocols/%s", hash), nil)
	if err != nil {
		return ProtocolInfo{}, err
	}
	var result ProtocolInfo
	err = req.doWithJSONResponse(ctx, api.client, &result)
	return result, err
}

// Environment -
func (api *Protocols) Environment(ctx context.Context, hash string) (int, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("protocols/%s/environment", hash), nil)
	if err != nil {
		return 0, err
	}
	var result int
	err = req.doWithJSONResponse(ctx, api.client, &result)
	return result, err
}
