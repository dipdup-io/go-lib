package api

// GetRights -
func (tzkt *API) GetRights(filters map[string]string) (rights []Right, err error) {
	err = tzkt.json("/v1/rights", filters, &rights)
	return
}
