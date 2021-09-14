package api

import "fmt"

// GetBigmapsCount -
func (tzkt *API) GetBigmapsCount() (uint64, error) {
	return tzkt.count("/v1/bigmaps/count")
}

// GetBigmaps -
func (tzkt *API) GetBigmaps(filters map[string]string) (response []BigMap, err error) {
	err = tzkt.json("/v1/bigmaps", filters, &response)
	return
}

// GetBigmapUpdates -
func (tzkt *API) GetBigmapUpdates(filters map[string]string) (response []BigMapUpdate, err error) {
	err = tzkt.json("/v1/bigmaps/updates", filters, &response)
	return
}

// GetBigmapByID -
func (tzkt *API) GetBigmapByID(id uint64, filters map[string]string) (response BigMap, err error) {
	err = tzkt.json(fmt.Sprintf("/v1/bigmaps/%d", id), filters, &response)
	return
}

// GetBigmapKeys -
func (tzkt *API) GetBigmapKeys(id uint64, filters map[string]string) (response []BigMapKey, err error) {
	err = tzkt.json(fmt.Sprintf("/v1/bigmaps/%d/keys", id), filters, &response)
	return
}

// GetBigmapKey -
func (tzkt *API) GetBigmapKey(id uint64, key string, filters map[string]string) (response BigMapKey, err error) {
	err = tzkt.json(fmt.Sprintf("/v1/bigmaps/%d/keys/%s", id, key), filters, &response)
	return
}

// GetBigmapKeyUpdates -
func (tzkt *API) GetBigmapKeyUpdates(id uint64, key string, filters map[string]string) (response []BigMapKeyUpdate, err error) {
	err = tzkt.json(fmt.Sprintf("/v1/bigmaps/%d/keys/%s/updates", id, key), filters, &response)
	return
}
