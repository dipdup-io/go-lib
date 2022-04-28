package api

import (
	"context"
	"fmt"
)

// GetBlock -
func (tzkt *API) GetBlock(ctx context.Context, level uint64) (b Block, err error) {
	err = tzkt.json(ctx, fmt.Sprintf("/v1/blocks/%d", level), nil, false, &b)
	return
}

// GetBlocks -
func (tzkt *API) GetBlocks(ctx context.Context, filters map[string]string) (b []Block, err error) {
	err = tzkt.json(ctx, "/v1/blocks", filters, false, &b)
	return
}
