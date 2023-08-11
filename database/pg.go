package database

//go:generate mockgen -destination ../mocks/mock_SchemeCommenter.go -package mocks github.com/dipdup-net/go-lib/database SchemeCommenter

import (
	"context"
	"fmt"

	"github.com/dipdup-net/go-lib/config"
	pg "github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
)

// PgGo -
type PgGo struct {
	conn *pg.DB
}

// NewPgGo -
func NewPgGo() *PgGo {
	return new(PgGo)
}

// DB -
func (db *PgGo) DB() *pg.DB {
	return db.conn
}

// Connect -
func (db *PgGo) Connect(ctx context.Context, cfg config.Database) error {
	if cfg.Kind != config.DBKindPostgres {
		return errors.Wrap(ErrUnsupportedDatabaseType, cfg.Kind)
	}
	var conn *pg.DB
	if cfg.Path != "" {
		opt, err := pg.ParseURL(cfg.Path)
		if err != nil {
			return err
		}
		conn = pg.Connect(opt)
	} else {
		conn = pg.Connect(&pg.Options{
			Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			User:     cfg.User,
			Password: cfg.Password,
			Database: cfg.Database,
		})
	}
	db.conn = conn
	return nil
}

// Close -
func (db *PgGo) Close() error {
	return db.conn.Close()
}

// Ping -
func (db *PgGo) Ping(ctx context.Context) error {
	if db.conn == nil {
		return ErrConnectionIsNotInitialized
	}
	return db.conn.Ping(ctx)
}

// State -
func (db *PgGo) State(ctx context.Context, indexName string) (*State, error) {
	var s State
	err := db.conn.ModelContext(ctx, &s).Where("index_name = ?", indexName).Limit(1).Select()
	return &s, err
}

// CreateState -
func (db *PgGo) CreateState(ctx context.Context, s *State) error {
	_, err := db.conn.ModelContext(ctx, s).Insert()
	return err
}

// UpdateState -
func (db *PgGo) UpdateState(ctx context.Context, s *State) error {
	_, err := db.conn.ModelContext(ctx, s).Where("index_name = ?", s.IndexName).Update()
	return err
}

// DeleteState -
func (db *PgGo) DeleteState(ctx context.Context, s *State) error {
	_, err := db.conn.ModelContext(ctx, s).Where("index_name = ?", s.IndexName).Delete()
	return err
}

// MakeTableComment -
func (db *PgGo) MakeTableComment(ctx context.Context, name string, comment string) error {
	_, err := db.conn.ExecContext(ctx,
		`COMMENT ON TABLE ? IS ?`,
		pg.Safe(name),
		comment)

	return err
}

// MakeColumnComment -
func (db *PgGo) MakeColumnComment(ctx context.Context, tableName string, columnName string, comment string) error {
	_, err := db.conn.ExecContext(ctx,
		`COMMENT ON COLUMN ?.? IS ?`,
		pg.Safe(tableName),
		pg.Safe(columnName),
		comment)

	return err
}
