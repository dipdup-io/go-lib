package api

import (
	"context"
	"fmt"

	"github.com/dipdup-net/go-lib/tzkt/data"
)

// GetBigmapsCount -
func (tzkt *API) GetBigmapsCount(ctx context.Context) (uint64, error) {
	return tzkt.count(ctx, "/v1/bigmaps/count", nil)
}

// GetBigmaps -
func (tzkt *API) GetBigmaps(ctx context.Context, filters map[string]string) (response []data.BigMap, err error) {
	err = tzkt.json(ctx, "/v1/bigmaps", filters, false, &response)
	return
}

// GetBigmapUpdates -
func (tzkt *API) GetBigmapUpdates(ctx context.Context, filters map[string]string) (response []data.BigMapUpdate, err error) {
	err = tzkt.json(ctx, "/v1/bigmaps/updates", filters, false, &response)
	return
}

// GetBigmapByID -
func (tzkt *API) GetBigmapByID(ctx context.Context, id uint64, filters map[string]string) (response data.BigMap, err error) {
	err = tzkt.json(ctx, fmt.Sprintf("/v1/bigmaps/%d", id), filters, false, &response)
	return
}

// GetBigmapKeys -
func (tzkt *API) GetBigmapKeys(ctx context.Context, id uint64, filters map[string]string) (response []data.BigMapKey, err error) {
	err = tzkt.json(ctx, fmt.Sprintf("/v1/bigmaps/%d/keys", id), filters, false, &response)
	return
}

// GetBigmapKey -
func (tzkt *API) GetBigmapKey(ctx context.Context, id uint64, key string, filters map[string]string) (response data.BigMapKey, err error) {
	err = tzkt.json(ctx, fmt.Sprintf("/v1/bigmaps/%d/keys/%s", id, key), filters, false, &response)
	return
}

// GetBigmapKeyUpdates -
func (tzkt *API) GetBigmapKeyUpdates(ctx context.Context, id uint64, key string, filters map[string]string) (response []data.BigMapKeyUpdate, err error) {
	err = tzkt.json(ctx, fmt.Sprintf("/v1/bigmaps/%d/keys/%s/updates", id, key), filters, false, &response)
	return
}
