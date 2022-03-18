package node

import "context"

// GeneralAPI -
type GeneralAPI interface {
	Version(ctx context.Context) (Version, error)
	StatsGC(ctx context.Context) (StatsGC, error)
	StatsMemory(ctx context.Context) (StatsMemory, error)

	URL() string
}

// General -
type General struct {
	baseURL string
	client  *client
}

// NewGeneral -
func NewGeneral(baseURL string) *General {
	return &General{
		baseURL: baseURL,
		client:  newClient(),
	}
}

// URL -
func (api *General) URL() string {
	return api.baseURL
}

// Version -
func (api *General) Version(ctx context.Context) (Version, error) {
	req, err := newGetRequest(api.baseURL, "version", nil)
	if err != nil {
		return Version{}, err
	}
	var result Version
	err = req.doWithJSONResponse(ctx, api.client, &result)
	return result, err
}

// StatsGC -
func (api *General) StatsGC(ctx context.Context) (StatsGC, error) {
	req, err := newGetRequest(api.baseURL, "stats/gc", nil)
	if err != nil {
		return StatsGC{}, err
	}
	var result StatsGC
	err = req.doWithJSONResponse(ctx, api.client, &result)
	return result, err
}

// StatsMemory -
func (api *General) StatsMemory(ctx context.Context) (StatsMemory, error) {
	req, err := newGetRequest(api.baseURL, "stats/memory", nil)
	if err != nil {
		return StatsMemory{}, err
	}
	var result StatsMemory
	err = req.doWithJSONResponse(ctx, api.client, &result)
	return result, err
}
