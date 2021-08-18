package state

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dipdup-net/go-lib/config"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// CheckConnection
func CheckConnection(db *gorm.DB) error {
	sql, err := db.DB()
	if err != nil {
		return err
	}

	if err = sql.Ping(); err != nil {
		return err
	}

	return nil
}

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

	db, err := gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				IgnoreRecordNotFoundError: true,
			},
		),
	})
	if err != nil {
		return nil, err
	}

	logrus.Info("Waiting database is up and runnning")
	for err := CheckConnection(db); err != nil; err = CheckConnection(db) {
		logrus.Warn("%v", err)
		time.Sleep(time.Second)
	}

	return db, nil
}
