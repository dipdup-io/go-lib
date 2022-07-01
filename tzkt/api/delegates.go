package api

import (
	"context"
	"fmt"

	"github.com/dipdup-net/go-lib/tzkt/data"
)

// GetDelegates -
func (tzkt *API) GetDelegates(ctx context.Context, filters map[string]string) (delegates []data.Delegate, err error) {
	err = tzkt.json(ctx, "/v1/delegates", filters, false, &delegates)
	return
}

// GetDelegatesCount -
func (tzkt *API) GetDelegatesCount(ctx context.Context) (uint64, error) {
	return tzkt.count(ctx, "/v1/delegates/count")
}

// GetDelegateByAddress -
func (tzkt *API) GetDelegateByAddress(ctx context.Context, address string) (delegate data.Delegate, err error) {
	err = tzkt.json(ctx, fmt.Sprintf("/v1/delegates/%s", address), nil, false, &delegate)
	return
}
