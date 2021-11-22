package api

import "context"

// GetEndorsements -
func (tzkt *API) GetEndorsements(ctx context.Context, filters map[string]string) (operations []Operation, err error) {
	err = tzkt.json(ctx, "/v1/operations/endorsements", filters, &operations)
	return
}

// GetBallots -
func (tzkt *API) GetBallots(ctx context.Context, filters map[string]string) (operations []Operation, err error) {
	err = tzkt.json(ctx, "/v1/operations/ballots", filters, &operations)
	return
}

// GetProposals -
func (tzkt *API) GetProposals(ctx context.Context, filters map[string]string) (operations []Operation, err error) {
	err = tzkt.json(ctx, "/v1/operations/proposals", filters, &operations)
	return
}

// GetActivations -
func (tzkt *API) GetActivations(ctx context.Context, filters map[string]string) (operations []Operation, err error) {
	err = tzkt.json(ctx, "/v1/operations/activations", filters, &operations)
	return
}

// GetDoubleBakings -
func (tzkt *API) GetDoubleBakings(ctx context.Context, filters map[string]string) (operations []Operation, err error) {
	err = tzkt.json(ctx, "/v1/operations/double_baking", filters, &operations)
	return
}

// GetDoubleEndorsings -
func (tzkt *API) GetDoubleEndorsings(ctx context.Context, filters map[string]string) (operations []Operation, err error) {
	err = tzkt.json(ctx, "/v1/operations/double_endorsing", filters, &operations)
	return
}

// GetNonceRevelations -
func (tzkt *API) GetNonceRevelations(ctx context.Context, filters map[string]string) (operations []Operation, err error) {
	err = tzkt.json(ctx, "/v1/operations/nonce_revelations", filters, &operations)
	return
}

// GetDelegations -
func (tzkt *API) GetDelegations(ctx context.Context, filters map[string]string) (operations []Operation, err error) {
	err = tzkt.json(ctx, "/v1/operations/delegations", filters, &operations)
	return
}

// GetOriginations -
func (tzkt *API) GetOriginations(ctx context.Context, filters map[string]string) (operations []Operation, err error) {
	err = tzkt.json(ctx, "/v1/operations/originations", filters, &operations)
	return
}

// GetTransactions -
func (tzkt *API) GetTransactions(ctx context.Context, filters map[string]string) (operations []Operation, err error) {
	err = tzkt.json(ctx, "/v1/operations/transactions", filters, &operations)
	return
}

// GetReveals -
func (tzkt *API) GetReveals(ctx context.Context, filters map[string]string) (operations []Operation, err error) {
	err = tzkt.json(ctx, "/v1/operations/reveals", filters, &operations)
	return
}

// GetRegisterConstants -
func (tzkt *API) GetRegisterConstants(ctx context.Context, filters map[string]string) (operations []Operation, err error) {
	err = tzkt.json(ctx, "/v1/operations/register_constants", filters, &operations)
	return
}
