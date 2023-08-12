package database

import (
	"context"
	"database/sql/driver"
	"io"
	"time"

	"github.com/dipdup-net/go-lib/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// SchemeCommenter -
type SchemeCommenter interface {
	MakeTableComment(ctx context.Context, name string, comment string) error
	MakeColumnComment(ctx context.Context, tableName string, columnName string, comment string) error
}

// Database -
type Database interface {
	Connect(ctx context.Context, cfg config.Database) error
	Exec(ctx context.Context, query string, args ...any) (int64, error)

	StateRepository
	SchemeCommenter

	driver.Pinger
	io.Closer
}

// errors
var (
	ErrConnectionIsNotInitialized = errors.New("connection is not initialized")
	ErrUnsupportedDatabaseType    = errors.New("unsupported database type")
)

// Wait -
func Wait(ctx context.Context, db driver.Pinger, checkPeriod time.Duration) {
	log.Info().Msg("Waiting database is up and runnning")
	if err := db.Ping(ctx); err == nil {
		return
	}

	ticker := time.NewTicker(checkPeriod)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := db.Ping(ctx); err != nil {
				log.Warn().Err(err).Msg("waiting...")
				continue
			}
			return
		}
	}
}
