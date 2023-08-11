package database

import (
	"context"
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

// Gorm -
type Gorm struct {
	conn *gorm.DB
}

// NewGorm -
func NewGorm() *Gorm {
	return new(Gorm)
}

// DB -
func (db *Gorm) DB() *gorm.DB {
	return db.conn
}

// Connect -
func (db *Gorm) Connect(ctx context.Context, cfg config.Database) error {
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
		return errors.Wrap(ErrUnsupportedDatabaseType, cfg.Kind)
	}

	conn, err := gorm.Open(dialector, &gorm.Config{
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
		return err
	}
	db.conn = conn

	return nil
}

// Close -
func (db *Gorm) Close() error {
	sql, err := db.conn.DB()
	if err != nil {
		return err
	}
	return sql.Close()
}

// Ping -
func (db *Gorm) Ping(ctx context.Context) error {
	if db.conn == nil {
		return ErrConnectionIsNotInitialized
	}
	sql, err := db.conn.DB()
	if err != nil {
		return err
	}

	if err = sql.PingContext(ctx); err != nil {
		return err
	}

	return nil
}

// State -
func (db *Gorm) State(ctx context.Context, indexName string) (*State, error) {
	var s State
	err := db.conn.WithContext(ctx).Where("index_name = ?", indexName).First(s).Error
	return &s, err
}

// CreateState -
func (db *Gorm) CreateState(ctx context.Context, s *State) error {
	return db.conn.WithContext(ctx).Create(s).Error
}

// UpdateState -
func (db *Gorm) UpdateState(ctx context.Context, s *State) error {
	return db.conn.WithContext(ctx).Save(s).Error
}

// DeleteState -
func (db *Gorm) DeleteState(ctx context.Context, s *State) error {
	return db.conn.WithContext(ctx).Delete(s).Error
}

// MakeTableComment -
func (db *Gorm) MakeTableComment(ctx context.Context, name string, comment string) error {
	return db.conn.WithContext(ctx).Exec(
		`COMMENT ON TABLE ? IS ?`,
		name,
		comment).Error
}

// MakeColumnComment -
func (db *Gorm) MakeColumnComment(ctx context.Context, tableName string, columnName string, comment string) error {
	return db.conn.WithContext(ctx).Exec(
		`COMMENT ON COLUMN ?.? IS ?`,
		tableName,
		columnName,
		comment).Error
}
