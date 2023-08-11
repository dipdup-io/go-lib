package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dipdup-net/go-lib/config"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

// Bun -
type Bun struct {
	sqldb *sql.DB
	conn  *bun.DB
}

// NewBun -
func NewBun() *Bun {
	return new(Bun)
}

// DB -
func (db *Bun) DB() *bun.DB {
	return db.conn
}

// Connect -
func (db *Bun) Connect(ctx context.Context, cfg config.Database) error {
	if cfg.Kind != config.DBKindPostgres {
		return errors.Wrap(ErrUnsupportedDatabaseType, cfg.Kind)
	}
	if cfg.Path != "" {
		db.sqldb = sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(cfg.Path)))
		db.conn = bun.NewDB(db.sqldb, pgdialect.New())
	} else {
		db.sqldb = sql.OpenDB(pgdriver.NewConnector(
			pgdriver.WithAddr(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)),
			pgdriver.WithDatabase(cfg.Database),
			pgdriver.WithPassword(cfg.Password),
			pgdriver.WithUser(cfg.User),
			pgdriver.WithInsecure(true),
		))
		db.conn = bun.NewDB(db.sqldb, pgdialect.New())
	}
	return nil
}

// Close -
func (db *Bun) Close() error {
	if err := db.conn.Close(); err != nil {
		return err
	}
	return db.sqldb.Close()
}

// Ping -
func (db *Bun) Ping(ctx context.Context) error {
	if db.conn == nil {
		return ErrConnectionIsNotInitialized
	}
	return db.conn.PingContext(ctx)
}

// State -
func (db *Bun) State(ctx context.Context, indexName string) (*State, error) {
	var s State
	err := db.conn.NewSelect().Model(&s).Where("index_name = ?", indexName).Limit(1).Scan(ctx)
	return &s, err
}

// CreateState -
func (db *Bun) CreateState(ctx context.Context, s *State) error {
	_, err := db.conn.NewInsert().Model(s).Exec(ctx)
	return err
}

// UpdateState -
func (db *Bun) UpdateState(ctx context.Context, s *State) error {
	_, err := db.conn.NewUpdate().Model(s).Where("index_name = ?", s.IndexName).Exec(ctx)
	return err
}

// DeleteState -
func (db *Bun) DeleteState(ctx context.Context, s *State) error {
	_, err := db.conn.NewDelete().Model(s).Where("index_name = ?", s.IndexName).Exec(ctx)
	return err
}

// MakeTableComment -
func (db *Bun) MakeTableComment(ctx context.Context, name string, comment string) error {
	_, err := db.conn.ExecContext(ctx,
		`COMMENT ON TABLE ? IS ?`,
		bun.Ident(name),
		comment)

	return err
}

// MakeColumnComment -
func (db *Bun) MakeColumnComment(ctx context.Context, tableName string, columnName string, comment string) error {
	_, err := db.conn.ExecContext(ctx,
		`COMMENT ON COLUMN ?.? IS ?`,
		bun.Ident(tableName),
		bun.Ident(columnName),
		comment)

	return err
}
