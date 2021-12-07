package database

import (
	"context"
	"database/sql/driver"
	"io"
	"time"

	"github.com/dipdup-net/go-lib/config"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Database -
type Database interface {
	Connect(ctx context.Context, cfg config.Database) error

	StateRepository

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
	logrus.Info("Waiting database is up and runnning")
	if err := db.Ping(ctx); err == nil {
		return
	}

	ticker := time.NewTicker(checkPeriod)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := db.Ping(ctx); err == nil {
				logrus.Warn(err)
				continue
			}
			return
		}
	}
}
