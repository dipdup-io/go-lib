package api

import "context"

// GetTokenTransfers -
func (tzkt *API) GetTokenTransfers(ctx context.Context, filters map[string]string) (transfers []Transfer, err error) {
	err = tzkt.json(ctx, "/v1/tokens/transfers", filters, false, &transfers)
	return
}
