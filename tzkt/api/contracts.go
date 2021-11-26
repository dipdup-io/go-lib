package api

import (
	"context"
	"fmt"
)

// GetContractJSONSchema -
func (tzkt *API) GetContractJSONSchema(ctx context.Context, address string) (response ContractJSONSchema, err error) {
	err = tzkt.json(ctx, fmt.Sprintf("/v1/contracts/%s/interface", address), nil, &response)
	return
}

// GetContractStorage -
func (tzkt *API) GetContractStorage(ctx context.Context, address string, output interface{}) error {
	return tzkt.json(ctx, fmt.Sprintf("/v1/contracts/%s/storage", address), nil, &output)
}
