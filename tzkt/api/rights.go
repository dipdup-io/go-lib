package api

import (
	"context"

	"github.com/dipdup-io/go-lib/tzkt/data"
)

// GetRights -
func (tzkt *API) GetRights(ctx context.Context, filters map[string]string) (rights []data.Right, err error) {
	err = tzkt.json(ctx, "/v1/rights", filters, false, &rights)
	return
}
