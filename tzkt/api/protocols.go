package api

import (
	"context"

	"github.com/dipdup-net/go-lib/tzkt/data"
)

// Protocols - Returns a list of protocols.
func (tzkt *API) Protocols(ctx context.Context, args map[string]string) (protocols []data.Protocol, err error) {
	err = tzkt.json(ctx, "/v1/protocols", args, false, &protocols)
	return
}
