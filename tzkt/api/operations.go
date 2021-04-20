package api

// GetEndorsements -
func (tzkt *API) GetEndorsements(filters map[string]string) (operations []Operation, err error) {
	err = tzkt.json("/v1/operations/endorsements", filters, &operations)
	return
}

// GetBallots -
func (tzkt *API) GetBallots(filters map[string]string) (operations []Operation, err error) {
	err = tzkt.json("/v1/operations/ballots", filters, &operations)
	return
}

// GetProposals -
func (tzkt *API) GetProposals(filters map[string]string) (operations []Operation, err error) {
	err = tzkt.json("/v1/operations/proposals", filters, &operations)
	return
}

// GetActivations -
func (tzkt *API) GetActivations(filters map[string]string) (operations []Operation, err error) {
	err = tzkt.json("/v1/operations/activations", filters, &operations)
	return
}

// GetDoubleBakings -
func (tzkt *API) GetDoubleBakings(filters map[string]string) (operations []Operation, err error) {
	err = tzkt.json("/v1/operations/double_baking", filters, &operations)
	return
}

// GetDoubleEndorsings -
func (tzkt *API) GetDoubleEndorsings(filters map[string]string) (operations []Operation, err error) {
	err = tzkt.json("/v1/operations/double_endorsing", filters, &operations)
	return
}

// GetNonceRevelations -
func (tzkt *API) GetNonceRevelations(filters map[string]string) (operations []Operation, err error) {
	err = tzkt.json("/v1/operations/nonce_revelations", filters, &operations)
	return
}

// GetDelegations -
func (tzkt *API) GetDelegations(filters map[string]string) (operations []Operation, err error) {
	err = tzkt.json("/v1/operations/delegations", filters, &operations)
	return
}

// GetOriginations -
func (tzkt *API) GetOriginations(filters map[string]string) (operations []Operation, err error) {
	err = tzkt.json("/v1/operations/originations", filters, &operations)
	return
}

// GetTransactions -
func (tzkt *API) GetTransactions(filters map[string]string) (operations []Operation, err error) {
	err = tzkt.json("/v1/operations/transactions", filters, &operations)
	return
}

// GetReveals -
func (tzkt *API) GetReveals(filters map[string]string) (operations []Operation, err error) {
	err = tzkt.json("/v1/operations/reveals", filters, &operations)
	return
}
