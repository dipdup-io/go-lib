package api

import (
	"context"
	"fmt"
)

// AccountCounter - Returns account counter
func (tzkt *API) AccountCounter(ctx context.Context, address string) (uint64, error) {
	return tzkt.count(ctx, fmt.Sprintf("/v1/accounts/%s/counter", address), nil)
}

// AccountCounter - Returns a number of accounts.
func (tzkt *API) AccountsCount(ctx context.Context, filters map[string]string) (uint64, error) {
	return tzkt.count(ctx, "/v1/accounts/count", filters)
}
