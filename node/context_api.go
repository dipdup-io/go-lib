package node

import (
	"context"
	stdJSON "encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// ContextAPI -
type ContextAPI interface {
	BigMap(ctx context.Context, blockID string, bigMapID uint64) ([]byte, error)
	BigMapKey(ctx context.Context, blockID string, bigMapID uint64, key string) (interface{}, error)
	CacheContracts(ctx context.Context, blockID string) (interface{}, error)
	CacheContractsSize(ctx context.Context, blockID string) (uint64, error)
	CacheContractsSizeLimit(ctx context.Context, blockID string) (uint64, error)
	Constants(ctx context.Context, blockID string) (Constants, error)
	Contracts(ctx context.Context, blockID string) ([]string, error)
	Contract(ctx context.Context, blockID, contract string) (ContractInfo, error)
	ContractBalance(ctx context.Context, blockID, contract string) (string, error)
	ContractCounter(ctx context.Context, blockID, contract string) (string, error)
	ContractDelegate(ctx context.Context, blockID, contract string) (string, error)
	ContractEntrypoints(ctx context.Context, blockID, contract string) (Entrypoints, error)
	ContractEntrypoint(ctx context.Context, blockID, contract, entrypoint string) (stdJSON.RawMessage, error)
	ContractScript(ctx context.Context, blockID, contract string) (Script, error)
	ContractStorage(ctx context.Context, blockID, contract string) (stdJSON.RawMessage, error)
	Delegates(ctx context.Context, blockID string, active DelegateType) ([]string, error)
	Delegate(ctx context.Context, blockID, pkh string) (Delegate, error)
	DelegateDeactivated(ctx context.Context, blockID, pkh string) (bool, error)
	DelegateBalance(ctx context.Context, blockID, pkh string) (string, error)
	DelegateContracts(ctx context.Context, blockID, pkh string) ([]string, error)
	DelegateGracePeriod(ctx context.Context, blockID, pkh string) (int, error)
	DelegateStakingBalance(ctx context.Context, blockID, pkh string) (string, error)
	DelegateVotingPower(ctx context.Context, blockID, pkh string) (int, error)
	ActiveDelegatesWithRolls(ctx context.Context, blockID string) ([]string, error)
	LiquidityBakingCPMMAddress(ctx context.Context, blockID string) (string, error)
}

// Context -
type Context struct {
	baseURL string
	chainID string
	client  *client
}

// NewContext -
func NewContext(baseURL, chainID string) *Context {
	return &Context{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		chainID: chainID,
		client:  newClient(),
	}
}

// NewMainContext -
func NewMainContext(baseURL string) *Context {
	return NewContext(baseURL, "main")
}

// BigMap -
func (api *Context) BigMap(ctx context.Context, blockID string, bigMapID uint64) ([]byte, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s/context/big_maps/%d", api.chainID, blockID, bigMapID), nil)
	if err != nil {
		return nil, err
	}
	return req.doWithBytesResponse(ctx, api.client)
}

// BigMapKey -
func (api *Context) BigMapKey(ctx context.Context, blockID string, bigMapID uint64, key string) (interface{}, error) {
	return nil, errors.New("not implemented")
}

// CacheContracts -
func (api *Context) CacheContracts(ctx context.Context, blockID string) (interface{}, error) {
	return nil, errors.New("not implemented")
}

// CacheContractsSize -
func (api *Context) CacheContractsSize(ctx context.Context, blockID string) (uint64, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s/context/cache/contracts/size", api.chainID, blockID), nil)
	if err != nil {
		return 0, err
	}
	var result uint64
	err = req.doWithJSONResponse(ctx, api.client, &result)
	return result, err
}

// CacheContractsSizeLimit -
func (api *Context) CacheContractsSizeLimit(ctx context.Context, blockID string) (uint64, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s/context/cache/contracts/size_limit", api.chainID, blockID), nil)
	if err != nil {
		return 0, err
	}
	var result uint64
	err = req.doWithJSONResponse(ctx, api.client, &result)
	return result, err
}

// Constants -
func (api *Context) Constants(ctx context.Context, blockID string) (Constants, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s/context/constants", api.chainID, blockID), nil)
	if err != nil {
		return Constants{}, err
	}
	var result Constants
	err = req.doWithJSONResponse(ctx, api.client, &result)
	return result, err
}

// Contracts -
func (api *Context) Contracts(ctx context.Context, blockID string) ([]string, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s/context/contracts", api.chainID, blockID), nil)
	if err != nil {
		return nil, err
	}
	var result []string
	err = req.doWithJSONResponse(ctx, api.client, &result)
	return result, err
}

// Contract -
func (api *Context) Contract(ctx context.Context, blockID, contract string) (ContractInfo, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s/context/contracts/%s", api.chainID, blockID, contract), nil)
	if err != nil {
		return ContractInfo{}, err
	}
	var result ContractInfo
	err = req.doWithJSONResponse(ctx, api.client, &result)
	return result, err
}

// ContractBalance -
func (api *Context) ContractBalance(ctx context.Context, blockID, contract string) (string, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s/context/contracts/%s/balance", api.chainID, blockID, contract), nil)
	if err != nil {
		return "", err
	}
	var result string
	err = req.doWithJSONResponse(ctx, api.client, &result)
	return result, err
}

// ContractCounter -
func (api *Context) ContractCounter(ctx context.Context, blockID, contract string) (string, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s/context/contracts/%s/counter", api.chainID, blockID, contract), nil)
	if err != nil {
		return "", err
	}
	var result string
	err = req.doWithJSONResponse(ctx, api.client, &result)
	return result, err
}

// ContractDelegate -
func (api *Context) ContractDelegate(ctx context.Context, blockID, contract string) (string, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s/context/contracts/%s/delegate", api.chainID, blockID, contract), nil)
	if err != nil {
		return "", err
	}
	var result string
	err = req.doWithJSONResponse(ctx, api.client, &result)
	return result, err
}

// ContractEntrypoints -
func (api *Context) ContractEntrypoints(ctx context.Context, blockID, contract string) (Entrypoints, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s/context/contracts/%s/entrypoints", api.chainID, blockID, contract), nil)
	if err != nil {
		return Entrypoints{}, err
	}
	var result Entrypoints
	err = req.doWithJSONResponse(ctx, api.client, &result)
	return result, err
}

// ContractEntrypoint -
func (api *Context) ContractEntrypoint(ctx context.Context, blockID, contract, entrypoint string) (stdJSON.RawMessage, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s/context/contracts/%s/entrypoint", api.chainID, blockID, contract), nil)
	if err != nil {
		return nil, err
	}
	return req.doWithBytesResponse(ctx, api.client)
}

// ContractScript -
func (api *Context) ContractScript(ctx context.Context, blockID, contract string) (Script, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s/context/contracts/%s/script", api.chainID, blockID, contract), nil)
	if err != nil {
		return Script{}, err
	}
	var result Script
	err = req.doWithJSONResponse(ctx, api.client, &result)
	return result, err
}

// ContractStorage -
func (api *Context) ContractStorage(ctx context.Context, blockID, contract string) (stdJSON.RawMessage, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s/context/contracts/%s/storage", api.chainID, blockID, contract), nil)
	if err != nil {
		return nil, err
	}
	return req.doWithBytesResponse(ctx, api.client)
}

// Delegates - Lists all registered delegates. `active` can get `active` or `inactive` values, othewise it skip.
func (api *Context) Delegates(ctx context.Context, blockID string, active DelegateType) ([]string, error) {
	url := fmt.Sprintf("chains/%s/blocks/%s/context/delegates", api.chainID, blockID)
	if active == ActiveDelegateType || active == InactiveDelegateType {
		url += fmt.Sprintf("?%s", active)
	}
	req, err := newGetRequest(api.baseURL, url, nil)
	if err != nil {
		return nil, err
	}
	var result []string
	err = req.doWithJSONResponse(ctx, api.client, &result)
	return result, err
}

// Delegate -
func (api *Context) Delegate(ctx context.Context, blockID, pkh string) (Delegate, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s/context/delegates/%s", api.chainID, blockID, pkh), nil)
	if err != nil {
		return Delegate{}, err
	}
	var result Delegate
	err = req.doWithJSONResponse(ctx, api.client, &result)
	return result, err
}

// DelegateDeactivated -
func (api *Context) DelegateDeactivated(ctx context.Context, blockID, pkh string) (bool, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s/context/delegates/%s/deactivated", api.chainID, blockID, pkh), nil)
	if err != nil {
		return false, err
	}
	var result bool
	err = req.doWithJSONResponse(ctx, api.client, &result)
	return result, err
}

// DelegateBalance -
func (api *Context) DelegateBalance(ctx context.Context, blockID, pkh string) (string, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s/context/delegates/%s/delegated_balance", api.chainID, blockID, pkh), nil)
	if err != nil {
		return "", err
	}
	var result string
	err = req.doWithJSONResponse(ctx, api.client, &result)
	return result, err
}

// DelegateContracts -
func (api *Context) DelegateContracts(ctx context.Context, blockID, pkh string) ([]string, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s/context/delegates/%s/delegated_contracts", api.chainID, blockID, pkh), nil)
	if err != nil {
		return nil, err
	}
	var result []string
	err = req.doWithJSONResponse(ctx, api.client, &result)
	return result, err
}

// DelegateGracePeriod -
func (api *Context) DelegateGracePeriod(ctx context.Context, blockID, pkh string) (int, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s/context/delegates/%s/grace_period", api.chainID, blockID, pkh), nil)
	if err != nil {
		return 0, err
	}
	var result int
	err = req.doWithJSONResponse(ctx, api.client, &result)
	return result, err
}

// DelegateStakingBalance -
func (api *Context) DelegateStakingBalance(ctx context.Context, blockID, pkh string) (string, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s/context/delegates/%s/staking_balance", api.chainID, blockID, pkh), nil)
	if err != nil {
		return "", err
	}
	var result string
	err = req.doWithJSONResponse(ctx, api.client, &result)
	return result, err
}

// DelegateVotingPower -
func (api *Context) DelegateVotingPower(ctx context.Context, blockID, pkh string) (int, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s/context/delegates/%s/voting_power", api.chainID, blockID, pkh), nil)
	if err != nil {
		return 0, err
	}
	var result int
	err = req.doWithJSONResponse(ctx, api.client, &result)
	return result, err
}

// ActiveDelegatesWithRolls -
func (api *Context) ActiveDelegatesWithRolls(ctx context.Context, blockID string) ([]string, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s/context/raw/json/active_delegates_with_rolls", api.chainID, blockID), nil)
	if err != nil {
		return nil, err
	}
	var result []string
	err = req.doWithJSONResponse(ctx, api.client, &result)
	return result, err
}

// LiquidityBakingCPMMAddress -
func (api *Context) LiquidityBakingCPMMAddress(ctx context.Context, blockID string) (string, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/blocks/%s/context/liquidity_baking/cpmm_address", api.chainID, blockID), nil)
	if err != nil {
		return "", err
	}
	var result string
	err = req.doWithJSONResponse(ctx, api.client, &result)
	return result, err
}
