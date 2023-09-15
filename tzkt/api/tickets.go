package api

import (
	"context"

	"github.com/dipdup-net/go-lib/tzkt/data"
)

// GetTicketsCount - Returns a total number of tickets.
func (tzkt *API) GetTicketsCount(ctx context.Context, filters map[string]string) (uint64, error) {
	return tzkt.count(ctx, "/v1/tickets/count", filters)
}

// GetTickets - Returns a list of tickets.
func (tzkt *API) GetTickets(ctx context.Context, filters map[string]string) (data []data.Ticket, err error) {
	err = tzkt.json(ctx, "/v1/tickets", filters, false, &data)
	return
}

// GetTicketBalancesCount - Returns a total number of ticket balances.
func (tzkt *API) GetTicketBalancesCount(ctx context.Context, filters map[string]string) (uint64, error) {
	return tzkt.count(ctx, "/v1/tickets/balances/count", filters)
}

// GetTicketBalances - Returns a list of ticket balances.
func (tzkt *API) GetTicketBalances(ctx context.Context, filters map[string]string) (data []data.TicketBalance, err error) {
	err = tzkt.json(ctx, "/v1/tickets/balances", filters, false, &data)
	return
}

// GetTicketTransfersCount - Returns the total number of ticket transfers.
func (tzkt *API) GetTicketTransfersCount(ctx context.Context, filters map[string]string) (uint64, error) {
	return tzkt.count(ctx, "/v1/tickets/transfers/count", filters)
}

// GetTicketTransfers - Returns a list of ticket transfers.
func (tzkt *API) GetTicketTransfers(ctx context.Context, filters map[string]string) (data []data.TicketTransfer, err error) {
	err = tzkt.json(ctx, "/v1/tickets/transfers", filters, false, &data)
	return
}
