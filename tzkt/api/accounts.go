package api

import (
	"context"
	"fmt"
)

// AccountCounter - Returns account counter
func (tzkt *API) AccountCounter(ctx context.Context, address string) (uint64, error) {
	return tzkt.count(ctx, fmt.Sprintf("/v1/accounts/%s/counter", address))
}
