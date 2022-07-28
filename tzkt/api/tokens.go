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

// GetTokens -
func (tzkt *API) GetTokens(ctx context.Context, filters map[string]string) (tokens []data.Token, err error) {
	err = tzkt.json(ctx, "/v1/tokens", filters, false, &tokens)
	return
}
