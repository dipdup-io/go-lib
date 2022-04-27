package api

import "context"

// GetAccountsMetadata -
func (tzkt *API) GetAccountsMetadata(ctx context.Context, filters map[string]string) ([]AccountMetadata, error) {
	var data []Metadata[AccountType]
	if err := tzkt.json(ctx, "/v1/metadata/accounts", filters, true, &data); err != nil {
		return nil, err
	}

	items := make([]AccountMetadata, len(data))
	for i := range data {
		items[i] = data[i].Value.Profile
		items[i].Address = data[i].Key
	}

	return items, nil
}
