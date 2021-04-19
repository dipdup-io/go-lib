package state

import (
	"log"
	"os"
	"time"

	"github.com/pkg/errors"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Supported database kinds
const (
	DBKindSqlite     = "sqlite"
	DBKindPostgres   = "postgres"
	DBKindMysql      = "mysql"
	DBKindClickHouse = "clickhouse"
	DBKindSqlServer  = "sqlserver"
)

// OpenConnection -
func OpenConnection(kind, path string) (*gorm.DB, error) {
	var dialector gorm.Dialector
	switch kind {
	case DBKindSqlite:
		dialector = sqlite.Open(path)
	case DBKindPostgres:
		dialector = postgres.Open(path)
	case DBKindMysql:
		dialector = mysql.Open(path)
	case DBKindClickHouse:
		dialector = clickhouse.Open(path)
	case DBKindSqlServer:
		dialector = sqlserver.Open(path)
	default:
		return nil, errors.Errorf("Unsupported database %s", kind)
	}

	return gorm.Open(dialector, &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				IgnoreRecordNotFoundError: true,
			},
		),
	})
}
