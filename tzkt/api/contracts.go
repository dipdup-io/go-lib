package api

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
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

// BuildContractParameters -
func (tzkt *API) BuildContractParameters(ctx context.Context, contract, entrypoint string, parameters interface{}) ([]byte, error) {
	response, err := tzkt.post(ctx, fmt.Sprintf("/v1/contracts/%s/entrypoints/%s/build", contract, entrypoint), nil, parameters)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusOK:
		return ioutil.ReadAll(response.Body)
	case http.StatusNoContent:
		return nil, nil
	default:
		return nil, errors.New(fmt.Sprintf("%s: %s %s", response.Status, entrypoint, contract))
	}
}
