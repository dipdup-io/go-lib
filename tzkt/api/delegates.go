package api

import "fmt"

// GetDelegates -
func (tzkt *API) GetDelegates(filters map[string]string) (delegates []Delegate, err error) {
	err = tzkt.json("/v1/delegates", filters, &delegates)
	return
}

// GetDelegatesCount -
func (tzkt *API) GetDelegatesCount() (uint64, error) {
	return tzkt.count("/v1/delegates/count")
}

// GetDelegateByAddress -
func (tzkt *API) GetDelegateByAddress(address string) (delegate Delegate, err error) {
	err = tzkt.json(fmt.Sprintf("/v1/delegates/%s", address), nil, &delegate)
	return
}
