package database

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"runtime"
	"time"

	"github.com/dipdup-io/go-lib/config"
	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

// Bun -
type Bun struct {
	sqldb *sql.DB
	conn  *bun.DB
	pool  *pgxpool.Pool
}

// NewBun -
func NewBun() *Bun {
	return new(Bun)
}

// DB -
func (db *Bun) DB() *bun.DB {
	return db.conn
}

func (db *Bun) Pool() *pgxpool.Pool {
	return db.pool
}

// Connect -
func (db *Bun) Connect(ctx context.Context, cfg config.Database) error {
	if cfg.Kind != config.DBKindPostgres {
		return errors.Wrap(ErrUnsupportedDatabaseType, cfg.Kind)
	}

	dsn := cfg.Path
	if dsn == "" {
		values := make(url.Values)
		values.Set("sslmode", "disable")
		if cfg.ApplicationName != "" {
			values.Set("application_name", cfg.ApplicationName)
		}

		dsn = fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s",
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.Database,
		)

		if len(values) > 0 {
			dsn = fmt.Sprintf("%s?%s", dsn, values.Encode())
		}
	}

	connCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return errors.Wrap(err, "parse postgres config")
	}

	connCfg.ConnConfig.RuntimeParams["TimeZone"] = "UTC"

	connCfg.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		conn.TypeMap().RegisterType(&pgtype.Type{
			Name:  "timestamp",
			OID:   pgtype.TimestampOID,
			Codec: &pgtype.TimestampCodec{ScanLocation: time.UTC},
		})
		conn.TypeMap().RegisterType(&pgtype.Type{
			Name:  "timestamptz",
			OID:   pgtype.TimestamptzOID,
			Codec: &pgtype.TimestamptzCodec{ScanLocation: time.UTC},
		})

		rows, err := conn.Query(ctx, "SELECT typname, oid FROM pg_type WHERE typtype = 'e'")
		if err != nil {
			return errors.Wrap(err, "query enum types")
		}
		defer rows.Close()
		for rows.Next() {
			var (
				name string
				oid  uint32
			)
			if err := rows.Scan(&name, &oid); err != nil {
				return errors.Wrap(err, "scan row")
			}
			conn.TypeMap().RegisterType(&pgtype.Type{
				Name:  name,
				OID:   oid,
				Codec: pgtype.TextCodec{},
			})
		}

		if err = rows.Err(); err != nil {
			return errors.Wrap(err, "rows iteration failed")
		}

		return nil
	}

	maxOpenConns := 4 * runtime.GOMAXPROCS(0)
	if cfg.MaxOpenConnections > 0 {
		maxOpenConns = cfg.MaxOpenConnections
	}
	connCfg.MaxConns = int32(maxOpenConns)

	maxLifetime := time.Minute
	if cfg.MaxLifetimeConnections > 0 {
		maxLifetime = time.Duration(cfg.MaxLifetimeConnections) * time.Second
	}
	connCfg.MaxConnLifetime = maxLifetime

	pool, err := pgxpool.NewWithConfig(ctx, connCfg)
	if err != nil {
		return errors.Wrap(err, "create pgxpool")
	}

	db.pool = pool
	db.sqldb = stdlib.OpenDBFromPool(pool)
	db.conn = bun.NewDB(db.sqldb, pgdialect.New())

	return nil
}

// Close -
func (db *Bun) Close() error {
	if err := db.conn.Close(); err != nil {
		return err
	}
	return db.sqldb.Close()
}

// Exec -
func (db *Bun) Exec(ctx context.Context, query string, args ...any) (int64, error) {
	result, err := db.conn.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
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

// CreateTable -
func (db *Bun) CreateTable(ctx context.Context, model any, opts ...CreateTableOption) error {
	if model == nil {
		return nil
	}
	var options CreateTableOptions
	for i := range opts {
		opts[i](&options)
	}

	query := db.DB().
		NewCreateTable().
		Model(model)

	if options.ifNotExists {
		query = query.IfNotExists()
	}

	if options.partitionBy != "" {
		query = query.PartitionBy(options.partitionBy)
	}

	if options.temporary {
		query = query.Temp()
	}

	_, err := query.Exec(ctx)
	return err
}
