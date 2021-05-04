package config

import "github.com/pkg/errors"

// Supported database kinds
const (
	DBKindSqlite   = "sqlite"
	DBKindPostgres = "postgres"
	DBKindMysql    = "mysql"
)

var (
	ErrUnsupportedDB = errors.New("Unsupported database")
	ErrMissingField  = errors.New("Missing field")
)
