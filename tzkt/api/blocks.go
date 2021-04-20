package api

import "fmt"

// GetBlock -
func (tzkt *API) GetBlock(level uint64) (b Block, err error) {
	err = tzkt.json(fmt.Sprintf("/v1/blocks/%d", level), nil, &b)
	return
}

// GetBlocks -
func (tzkt *API) GetBlocks(filters map[string]string) (b []Block, err error) {
	err = tzkt.json("/v1/blocks", filters, &b)
	return
}
