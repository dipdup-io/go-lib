package database

// CreateTableOptions -
type CreateTableOptions struct {
	partitionBy string
	ifNotExists bool
	temporary   bool
}

// CreateTableOption -
type CreateTableOption func(opts *CreateTableOptions)

// WithPartitioning -
func WithPartitioning(by string) CreateTableOption {
	return func(opts *CreateTableOptions) {
		opts.partitionBy = by
	}
}

// WithIfNotExists -
func WithIfNotExists() CreateTableOption {
	return func(opts *CreateTableOptions) {
		opts.ifNotExists = true
	}
}

// WithTemporary -
func WithTemporary() CreateTableOption {
	return func(opts *CreateTableOptions) {
		opts.temporary = true
	}
}
