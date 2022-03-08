package node

import "context"

// ConfigAPI -
type ConfigAPI interface {
	HistoryMode(ctx context.Context) (HistoryMode, error)
	UserActivatedProtocols(ctx context.Context) ([]ActivatedProtocol, error)
	UserActivatedUpgrades(ctx context.Context) ([]ActivatedUpgrades, error)
}

// Config -
type Config struct {
	baseURL string
	client  *client
}

// NewConfig -
func NewConfig(baseURL string) *Config {
	return &Config{
		baseURL: baseURL,
		client:  newClient(),
	}
}

// HistoryMode -
func (api *Config) HistoryMode(ctx context.Context) (HistoryMode, error) {
	req, err := newGetRequest(api.baseURL, "config/history_mode", nil)
	if err != nil {
		return HistoryMode{}, err
	}
	var result HistoryMode
	err = req.doWithJSONResponse(ctx, api.client, &result)
	return result, err
}

// UserActivatedProtocols -
func (api *Config) UserActivatedProtocols(ctx context.Context) ([]ActivatedProtocol, error) {
	req, err := newGetRequest(api.baseURL, "config/network/user_activated_protocol_overrides", nil)
	if err != nil {
		return nil, err
	}
	var result []ActivatedProtocol
	err = req.doWithJSONResponse(ctx, api.client, &result)
	return result, err
}

// UserActivatedUpgrades -
func (api *Config) UserActivatedUpgrades(ctx context.Context) ([]ActivatedUpgrades, error) {
	req, err := newGetRequest(api.baseURL, "config/network/user_activated_upgrades", nil)
	if err != nil {
		return nil, err
	}
	var result []ActivatedUpgrades
	err = req.doWithJSONResponse(ctx, api.client, &result)
	return result, err
}
