package api

import (
	"context"

	"github.com/dipdup-net/go-lib/tzkt/data"
)

// GetTokenTransfers -
func (tzkt *API) GetTokenTransfers(ctx context.Context, filters map[string]string) (transfers []data.Transfer, err error) {
	err = tzkt.json(ctx, "/v1/tokens/transfers", filters, false, &transfers)
	return
}
