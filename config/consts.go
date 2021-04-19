package config

import "github.com/pkg/errors"

// Supported database kinds
const (
	DBKindSqlite     = "sqlite"
	DBKindPostgres   = "postgres"
	DBKindMysql      = "mysql"
	DBKindClickHouse = "clickhouse"
	DBKindSqlServer  = "sqlserver"
)

var (
	ErrUnsupportedDB = errors.New("Unsupported database")
)
