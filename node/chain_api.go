package node

import (
	"context"
	"fmt"
	"strings"
)

// ChainAPI -
type ChainAPI interface {
	ChainID(ctx context.Context) (string, error)
	InvalidBlocks(ctx context.Context) ([]InvalidBlock, error)
	InvalidBlock(ctx context.Context, blockHash string) (InvalidBlock, error)
	IsBootstrapped(ctx context.Context) (Bootstrapped, error)
	LevelsCaboose(ctx context.Context) (Caboose, error)
	LevelsCheckpoint(ctx context.Context) (Checkpoint, error)
	LevelsSavepoint(ctx context.Context) (Savepoint, error)
	PendingOperations(ctx context.Context) (MempoolResponse, error)
}

// Chain -
type Chain struct {
	baseURL string
	chainID string
	client  *client
}

// NewChain -
func NewChain(baseURL, chainID string) *Chain {
	return &Chain{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		chainID: chainID,
		client:  newClient(),
	}
}

// NewMainChain -
func NewMainChain(baseURL string) *Chain {
	return NewChain(baseURL, "main")
}

// ChainID -
func (api *Chain) ChainID(ctx context.Context) (string, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/chain_id", api.chainID), nil)
	if err != nil {
		return "", err
	}
	var chainID string
	err = req.doWithJSONResponse(ctx, api.client, &chainID)
	return chainID, err
}

// InvalidBlocks -
func (api *Chain) InvalidBlocks(ctx context.Context) ([]InvalidBlock, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/invalid_blocks", api.chainID), nil)
	if err != nil {
		return nil, err
	}
	var result []InvalidBlock
	err = req.doWithJSONResponse(ctx, api.client, &result)
	return result, err
}

// InvalidBlock -
func (api *Chain) InvalidBlock(ctx context.Context, blockHash string) (InvalidBlock, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/invalid_blocks/%s", api.chainID, blockHash), nil)
	if err != nil {
		return InvalidBlock{}, err
	}
	var result InvalidBlock
	err = req.doWithJSONResponse(ctx, api.client, &result)
	return result, err
}

// IsBootstrapped -
func (api *Chain) IsBootstrapped(ctx context.Context) (Bootstrapped, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/is_bootstrapped", api.chainID), nil)
	if err != nil {
		return Bootstrapped{}, err
	}
	var result Bootstrapped
	err = req.doWithJSONResponse(ctx, api.client, &result)
	return result, err
}

// LevelsCaboose -
func (api *Chain) LevelsCaboose(ctx context.Context) (Caboose, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/levels/caboose", api.chainID), nil)
	if err != nil {
		return Caboose{}, err
	}
	var result Caboose
	err = req.doWithJSONResponse(ctx, api.client, &result)
	return result, err
}

// LevelsCheckpoint -
func (api *Chain) LevelsCheckpoint(ctx context.Context) (Checkpoint, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/levels/checkpoint", api.chainID), nil)
	if err != nil {
		return Checkpoint{}, err
	}
	var result Checkpoint
	err = req.doWithJSONResponse(ctx, api.client, &result)
	return result, err
}

// LevelsSavepoint -
func (api *Chain) LevelsSavepoint(ctx context.Context) (Savepoint, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/levels/savepoint", api.chainID), nil)
	if err != nil {
		return Savepoint{}, err
	}
	var result Savepoint
	err = req.doWithJSONResponse(ctx, api.client, &result)
	return result, err
}

// PendingOperations -
func (api *Chain) PendingOperations(ctx context.Context) (MempoolResponse, error) {
	req, err := newGetRequest(api.baseURL, fmt.Sprintf("chains/%s/mempool/pending_operations?version=1", api.chainID), nil)
	if err != nil {
		return MempoolResponse{}, err
	}
	var result MempoolResponse
	err = req.doWithJSONResponse(ctx, api.client, &result)
	return result, err
}
