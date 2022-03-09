package node

import (
	"context"
	"fmt"
)

// NetworkAPI -
type NetworkAPI interface {
	Connections(ctx context.Context) ([]Connection, error)
	Connection(ctx context.Context, peerID string) (Connection, error)
	Points(ctx context.Context) ([]NetworkPointWithURI, error)
	ConnectionVersion(ctx context.Context) (ConnectionVersion, error)
}

// Network -
type Network struct {
	baseURL string
	client  *client
}

// NewNetwork -
func NewNetwork(baseURL string) *Network {
	return &Network{
		baseURL: baseURL,
		client:  newClient(),
	}
}

// Connections -
func (api *Network) Connections(ctx context.Context) ([]Connection, error) {
	req, err := newGetRequest(api.baseURL, "network/connections", nil)
	if err != nil {
		return nil, err
	}
	var result []Connection
	err = req.doWithJSONResponse(ctx, api.client, &result)
	return result, err
}

// Connection -
func (api *Network) Connection(ctx context.Context, peerID string) (Connection, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("network/connections/%s", peerID), nil)
	if err != nil {
		return Connection{}, err
	}
	var result Connection
	err = req.doWithJSONResponse(ctx, api.client, &result)
	return result, err
}

// Points -
func (api *Network) Points(ctx context.Context) ([]NetworkPointWithURI, error) {
	req, err := newGetRequest(api.baseURL, "network/points", nil)
	if err != nil {
		return nil, err
	}
	var result []NetworkPointWithURI
	err = req.doWithJSONResponse(ctx, api.client, &result)
	return result, err
}

// ConnectionVersion -
func (api *Network) ConnectionVersion(ctx context.Context) (ConnectionVersion, error) {
	req, err := newGetRequest(api.baseURL, "network/version", nil)
	if err != nil {
		return ConnectionVersion{}, err
	}
	var result ConnectionVersion
	err = req.doWithJSONResponse(ctx, api.client, &result)
	return result, err
}
