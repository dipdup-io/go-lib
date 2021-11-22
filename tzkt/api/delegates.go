package api

import (
	"context"
	"fmt"
)

// GetDelegates -
func (tzkt *API) GetDelegates(ctx context.Context, filters map[string]string) (delegates []Delegate, err error) {
	err = tzkt.json(ctx, "/v1/delegates", filters, &delegates)
	return
}

// GetDelegatesCount -
func (tzkt *API) GetDelegatesCount(ctx context.Context) (uint64, error) {
	return tzkt.count(ctx, "/v1/delegates/count")
}

// GetDelegateByAddress -
func (tzkt *API) GetDelegateByAddress(ctx context.Context, address string) (delegate Delegate, err error) {
	err = tzkt.json(ctx, fmt.Sprintf("/v1/delegates/%s", address), nil, &delegate)
	return
}
