package api

import (
	"context"

	"github.com/dipdup-net/go-lib/tzkt/data"
)

// GetSmartRollupsCount - Returns a total number of smart rollups.
func (tzkt *API) GetSmartRollupsCount(ctx context.Context) (uint64, error) {
	return tzkt.count(ctx, "/v1/smart_rollups/count", nil)
}

// GetSmartRollups -
func (tzkt *API) GetSmartRollups(ctx context.Context, filters map[string]string) (response []data.SmartRollup, err error) {
	err = tzkt.json(ctx, "/v1/smart_rollups", filters, false, &response)
	return
}
