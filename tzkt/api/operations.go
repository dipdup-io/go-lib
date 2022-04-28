package api

import (
	"context"
)

// GetEndorsements -
func (tzkt *API) GetEndorsements(ctx context.Context, filters map[string]string) (operations []Operation, err error) {
	err = tzkt.json(ctx, "/v1/operations/endorsements", filters, false, &operations)
	return
}

// GetBallots -
func (tzkt *API) GetBallots(ctx context.Context, filters map[string]string) (operations []Operation, err error) {
	err = tzkt.json(ctx, "/v1/operations/ballots", filters, false, &operations)
	return
}

// GetProposals -
func (tzkt *API) GetProposals(ctx context.Context, filters map[string]string) (operations []Operation, err error) {
	err = tzkt.json(ctx, "/v1/operations/proposals", filters, false, &operations)
	return
}

// GetActivations -
func (tzkt *API) GetActivations(ctx context.Context, filters map[string]string) (operations []Operation, err error) {
	err = tzkt.json(ctx, "/v1/operations/activations", filters, false, &operations)
	return
}

// GetDoubleBakings -
func (tzkt *API) GetDoubleBakings(ctx context.Context, filters map[string]string) (operations []Operation, err error) {
	err = tzkt.json(ctx, "/v1/operations/double_baking", filters, false, &operations)
	return
}

// GetDoubleEndorsings -
func (tzkt *API) GetDoubleEndorsings(ctx context.Context, filters map[string]string) (operations []Operation, err error) {
	err = tzkt.json(ctx, "/v1/operations/double_endorsing", filters, false, &operations)
	return
}

// GetNonceRevelations -
func (tzkt *API) GetNonceRevelations(ctx context.Context, filters map[string]string) (operations []Operation, err error) {
	err = tzkt.json(ctx, "/v1/operations/nonce_revelations", filters, false, &operations)
	return
}

// GetDelegations -
func (tzkt *API) GetDelegations(ctx context.Context, filters map[string]string) (operations []Delegation, err error) {
	err = tzkt.json(ctx, "/v1/operations/delegations", filters, false, &operations)
	return
}

// GetOriginations -
func (tzkt *API) GetOriginations(ctx context.Context, filters map[string]string) (operations []Origination, err error) {
	err = tzkt.json(ctx, "/v1/operations/originations", filters, false, &operations)
	return
}

// GetTransactions -
func (tzkt *API) GetTransactions(ctx context.Context, filters map[string]string) (operations []Transaction, err error) {
	err = tzkt.json(ctx, "/v1/operations/transactions", filters, false, &operations)
	return
}

// GetReveals -
func (tzkt *API) GetReveals(ctx context.Context, filters map[string]string) (operations []Reveal, err error) {
	err = tzkt.json(ctx, "/v1/operations/reveals", filters, false, &operations)
	return
}

// GetRegisterConstants -
func (tzkt *API) GetRegisterConstants(ctx context.Context, filters map[string]string) (operations []Operation, err error) {
	err = tzkt.json(ctx, "/v1/operations/register_constants", filters, false, &operations)
	return
}

// GetMigrations -
func (tzkt *API) GetMigrations(ctx context.Context, filters map[string]string) (operations []Migration, err error) {
	err = tzkt.json(ctx, "/v1/operations/migrations", filters, false, &operations)
	return
}
