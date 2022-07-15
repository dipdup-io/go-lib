package api

import (
	"context"
	"fmt"

	"github.com/dipdup-net/go-lib/tools"
	"github.com/dipdup-net/go-lib/tzkt/data"
	"github.com/pkg/errors"
)

// GetEndorsements -
func (tzkt *API) GetEndorsements(ctx context.Context, filters map[string]string) (operations []data.Endorsement, err error) {
	err = tzkt.json(ctx, "/v1/operations/endorsements", filters, false, &operations)
	return
}

// GetBallots -
func (tzkt *API) GetBallots(ctx context.Context, filters map[string]string) (operations []data.Ballot, err error) {
	err = tzkt.json(ctx, "/v1/operations/ballots", filters, false, &operations)
	return
}

// GetProposals -
func (tzkt *API) GetProposals(ctx context.Context, filters map[string]string) (operations []data.Proposal, err error) {
	err = tzkt.json(ctx, "/v1/operations/proposals", filters, false, &operations)
	return
}

// GetActivations -
func (tzkt *API) GetActivations(ctx context.Context, filters map[string]string) (operations []data.Activation, err error) {
	err = tzkt.json(ctx, "/v1/operations/activations", filters, false, &operations)
	return
}

// GetDoubleBakings -
func (tzkt *API) GetDoubleBakings(ctx context.Context, filters map[string]string) (operations []data.DoubleBaking, err error) {
	err = tzkt.json(ctx, "/v1/operations/double_baking", filters, false, &operations)
	return
}

// GetDoubleEndorsings -
func (tzkt *API) GetDoubleEndorsings(ctx context.Context, filters map[string]string) (operations []data.DoubleEndorsing, err error) {
	err = tzkt.json(ctx, "/v1/operations/double_endorsing", filters, false, &operations)
	return
}

// GetNonceRevelations -
func (tzkt *API) GetNonceRevelations(ctx context.Context, filters map[string]string) (operations []data.NonceRevelation, err error) {
	err = tzkt.json(ctx, "/v1/operations/nonce_revelations", filters, false, &operations)
	return
}

// GetDelegations -
func (tzkt *API) GetDelegations(ctx context.Context, filters map[string]string) (operations []data.Delegation, err error) {
	err = tzkt.json(ctx, "/v1/operations/delegations", filters, false, &operations)
	return
}

// GetOriginations -
func (tzkt *API) GetOriginations(ctx context.Context, filters map[string]string) (operations []data.Origination, err error) {
	err = tzkt.json(ctx, "/v1/operations/originations", filters, false, &operations)
	return
}

// GetTransactions -
func (tzkt *API) GetTransactions(ctx context.Context, filters map[string]string) (operations []data.Transaction, err error) {
	err = tzkt.json(ctx, "/v1/operations/transactions", filters, false, &operations)
	return
}

// GetReveals -
func (tzkt *API) GetReveals(ctx context.Context, filters map[string]string) (operations []data.Reveal, err error) {
	err = tzkt.json(ctx, "/v1/operations/reveals", filters, false, &operations)
	return
}

// GetRegisterConstants -
func (tzkt *API) GetRegisterConstants(ctx context.Context, filters map[string]string) (operations []data.RegisterConstant, err error) {
	err = tzkt.json(ctx, "/v1/operations/register_constants", filters, false, &operations)
	return
}

// GetMigrations -
func (tzkt *API) GetMigrations(ctx context.Context, filters map[string]string) (operations []data.Migration, err error) {
	err = tzkt.json(ctx, "/v1/operations/migrations", filters, false, &operations)
	return
}

// GetOperationsByHash -
func (tzkt *API) GetOperationsByHash(ctx context.Context, hash string, filters map[string]string) (operations []data.Operation, err error) {
	err = tzkt.json(ctx, fmt.Sprintf("/v1/operations/%s", hash), filters, false, &operations)
	return
}

// GetTransactionsByHash -
func (tzkt *API) GetTransactionsByHash(ctx context.Context, hash string, filters map[string]string) (operations []data.Transaction, err error) {
	if !tools.IsOperationHash(hash) {
		return nil, errors.Errorf("invalid operation hash: %s", hash)
	}
	err = tzkt.json(ctx, fmt.Sprintf("/v1/operations/transactions/%s", hash), filters, false, &operations)
	return
}

// GetPreendorsement -
func (tzkt *API) GetPreendorsement(ctx context.Context, filters map[string]string) (operations []data.Preendorsement, err error) {
	err = tzkt.json(ctx, "/v1/operations/preendorsement", filters, false, &operations)
	return
}

// GetSetDepositsLimit -
func (tzkt *API) GetSetDepositsLimit(ctx context.Context, filters map[string]string) (operations []data.SetDepositsLimit, err error) {
	err = tzkt.json(ctx, "/v1/operations/set_deposits_limits", filters, false, &operations)
	return
}

// GetTxRollupCommit -
func (tzkt *API) GetTxRollupCommit(ctx context.Context, filters map[string]string) (operations []data.TxRollupCommit, err error) {
	err = tzkt.json(ctx, "/v1/operations/tx_rollup_commit", filters, false, &operations)
	return
}

// GetTxRollupDispatchTicket -
func (tzkt *API) GetTxRollupDispatchTicket(ctx context.Context, filters map[string]string) (operations []data.TxRollupDispatchTicket, err error) {
	err = tzkt.json(ctx, "/v1/operations/tx_rollup_dispatch_tickets", filters, false, &operations)
	return
}

// GetTxRollupFinalizeCommitment -
func (tzkt *API) GetTxRollupFinalizeCommitment(ctx context.Context, filters map[string]string) (operations []data.TxRollupFinalizeCommitment, err error) {
	err = tzkt.json(ctx, "/v1/operations/tx_rollup_finalize_commitment", filters, false, &operations)
	return
}

// GetTransferTicket -
func (tzkt *API) GetTxRollupOrigination(ctx context.Context, filters map[string]string) (operations []data.TxRollupOrigination, err error) {
	err = tzkt.json(ctx, "/v1/operations/tx_rollup_origination", filters, false, &operations)
	return
}

// GetTxRollupRejection -
func (tzkt *API) GetTxRollupRejection(ctx context.Context, filters map[string]string) (operations []data.TxRollupRejection, err error) {
	err = tzkt.json(ctx, "/v1/operations/tx_rollup_rejection", filters, false, &operations)
	return
}

// GetTxRollupRemoveCommitment -
func (tzkt *API) GetTxRollupRemoveCommitment(ctx context.Context, filters map[string]string) (operations []data.TxRollupRemoveCommitment, err error) {
	err = tzkt.json(ctx, "/v1/operations/tx_rollup_remove_commitment", filters, false, &operations)
	return
}

// GetTransferTicket -
func (tzkt *API) GetTxRollupReturnBond(ctx context.Context, filters map[string]string) (operations []data.TxRollupReturnBond, err error) {
	err = tzkt.json(ctx, "/v1/operations/tx_rollup_return_bond", filters, false, &operations)
	return
}

// GetTxRollupSubmitBatch -
func (tzkt *API) GetTxRollupSubmitBatch(ctx context.Context, filters map[string]string) (operations []data.TxRollupSubmitBatch, err error) {
	err = tzkt.json(ctx, "/v1/operations/tx_rollup_submit_batch", filters, false, &operations)
	return
}

// GetTransferTicket -
func (tzkt *API) GetTransferTicket(ctx context.Context, filters map[string]string) (operations []data.TransferTicket, err error) {
	err = tzkt.json(ctx, "/v1/operations/transfer_ticket", filters, false, &operations)
	return
}

// GetBakings -
func (tzkt *API) GetBakings(ctx context.Context, filters map[string]string) (operations []data.Baking, err error) {
	err = tzkt.json(ctx, "/v1/operations/baking", filters, false, &operations)
	return
}

// GetEndorsingRewards -
func (tzkt *API) GetEndorsingRewards(ctx context.Context, filters map[string]string) (operations []data.EndorsingReward, err error) {
	err = tzkt.json(ctx, "/v1/operations/endorsing_rewards", filters, false, &operations)
	return
}

// GetRevelationPenalties -
func (tzkt *API) GetRevelationPenalties(ctx context.Context, filters map[string]string) (operations []data.RevelationPenalty, err error) {
	err = tzkt.json(ctx, "/v1/operations/revelation_penalties", filters, false, &operations)
	return
}

// GetDoublePreendorsings -
func (tzkt *API) GetDoublePreendorsings(ctx context.Context, filters map[string]string) (operations []data.DoublePreendorsing, err error) {
	err = tzkt.json(ctx, "/v1/operations/double_preendorsing", filters, false, &operations)
	return
}
