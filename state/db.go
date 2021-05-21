package state

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dipdup-net/go-lib/config"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// OpenConnection -
func OpenConnection(cfg config.Database) (*gorm.DB, error) {
	var dialector gorm.Dialector
	switch cfg.Kind {
	case config.DBKindSqlite:
		dialector = sqlite.Open(cfg.Path)
	case config.DBKindPostgres:
		connString := cfg.Path
		if connString == "" {
			connString = fmt.Sprintf(
				"host=%s user=%s password=%s dbname=%s port=%d",
				cfg.Host, cfg.User, cfg.Password, cfg.Database, cfg.Port,
			)
		}
		dialector = postgres.Open(connString)
	case config.DBKindMysql:
		connString := cfg.Path
		if connString == "" {
			connString = fmt.Sprintf(
				"%s:%s@tcp(%s:%d)/%s",
				cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database,
			)
		}
		dialector = mysql.Open(connString)
	default:
		return nil, errors.Errorf("Unsupported database kind %s", cfg.Kind)
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
