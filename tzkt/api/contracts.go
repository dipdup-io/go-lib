package api

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/dipdup-net/go-lib/tzkt/data"
	"github.com/pkg/errors"
)

// GetContractJSONSchema -
func (tzkt *API) GetContractJSONSchema(ctx context.Context, address string) (response data.ContractJSONSchema, err error) {
	err = tzkt.json(ctx, fmt.Sprintf("/v1/contracts/%s/interface", address), nil, false, &response)
	return
}

// GetContractStorage -
func (tzkt *API) GetContractStorage(ctx context.Context, address string, output interface{}) error {
	return tzkt.json(ctx, fmt.Sprintf("/v1/contracts/%s/storage", address), nil, false, &output)
}

// BuildContractParameters -
func (tzkt *API) BuildContractParameters(ctx context.Context, contract, entrypoint string, parameters interface{}) ([]byte, error) {
	response, err := tzkt.post(ctx, fmt.Sprintf("/v1/contracts/%s/entrypoints/%s/build", contract, entrypoint), nil, parameters)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusOK:
		return io.ReadAll(response.Body)
	case http.StatusNoContent:
		return nil, nil
	default:
		return nil, errors.New(fmt.Sprintf("%s: %s %s", response.Status, entrypoint, contract))
	}
}

// GetContractByAddress -
func (tzkt *API) GetContractByAddress(ctx context.Context, address string) (response data.Contract, err error) {
	err = tzkt.json(ctx, fmt.Sprintf("/v1/contracts/%s", address), nil, false, &response)
	return
}

// ListContracts -
func (tzkt *API) ListContracts(ctx context.Context, filters map[string]string) (response []data.Contract, err error) {
	err = tzkt.json(ctx, "/v1/contracts", filters, false, &response)
	return
}
