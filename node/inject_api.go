package node

import (
	"context"
	"net/url"
)

// InjectAPI -
type InjectAPI interface {
	InjectOperation(ctx context.Context, request InjectOperationRequest) (string, error)
}

type Inject struct {
	baseURL string
	client  *client
}

// NewInject -
func NewInject(baseURL string) *Inject {
	return &Inject{
		baseURL: baseURL,
		client:  newClient(),
	}
}

// InjectOperation -
func (api *Inject) InjectOperation(ctx context.Context, request InjectOperationRequest) (string, error) {
	queryArgs := make(url.Values)
	if request.Async {
		queryArgs.Add("async", "true")
	}
	if request.ChainID != "" {
		queryArgs.Add("chain", request.ChainID)
	}
	req, err := newPostRequest(api.baseURL, "injection/operation", queryArgs, request.Operation)
	if err != nil {
		return "", err
	}
	var hash string
	err = req.doWithJSONResponse(ctx, api.client, &hash)
	return hash, err
}
