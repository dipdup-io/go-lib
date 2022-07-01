package api

import (
	"context"
	"fmt"

	"github.com/dipdup-net/go-lib/tzkt/data"
)

// GetBlock -
func (tzkt *API) GetBlock(ctx context.Context, level uint64) (b data.Block, err error) {
	err = tzkt.json(ctx, fmt.Sprintf("/v1/blocks/%d", level), nil, false, &b)
	return
}

// GetBlocks -
func (tzkt *API) GetBlocks(ctx context.Context, filters map[string]string) (b []data.Block, err error) {
	err = tzkt.json(ctx, "/v1/blocks", filters, false, &b)
	return
}
