package api

import "context"

// GetRights -
func (tzkt *API) GetRights(ctx context.Context, filters map[string]string) (rights []Right, err error) {
	err = tzkt.json(ctx, "/v1/rights", filters, false, &rights)
	return
}
