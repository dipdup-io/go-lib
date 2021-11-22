package node

import (
	"context"
)

// RequestOpts -
type RequestOpts struct {
	ctx context.Context
}

func newRequestOpts(opts ...RequestOption) RequestOpts {
	req := RequestOpts{context.Background()}

	for i := range opts {
		opts[i](&req)
	}
	return req
}

// RequestOption -
type RequestOption func(*RequestOpts)

// WithContext -
func WithContext(ctx context.Context) RequestOption {
	return func(opts *RequestOpts) {
		opts.ctx = ctx
	}
}
