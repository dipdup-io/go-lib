package api

import (
	"context"

	"github.com/dipdup-net/go-lib/tzkt/data"
)

// GetAccountsMetadata -
func (tzkt *API) GetAccountsMetadata(ctx context.Context, filters map[string]string) ([]data.AccountMetadata, error) {
	var raw []data.Metadata[data.AccountMetadata]
	if err := tzkt.json(ctx, "/v1/extras/accounts", filters, true, &raw); err != nil {
		return nil, err
	}

	items := make([]data.AccountMetadata, len(raw))
	for i := range raw {
		items[i] = raw[i].Value
		items[i].Address = raw[i].Key
	}

	return items, nil
}
