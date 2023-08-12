package database

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

// PartitionBy -
type PartitionBy int

const (
	PartitionByMonth PartitionBy = iota + 1
	PartitionByYear
)

type params struct {
	id      string
	start   time.Time
	end     time.Time
	success bool
}

// RangePartitionManager -
type RangePartitionManager struct {
	conn Database
	by   PartitionBy

	lastId string
}

// NewPartitionManager -
func NewPartitionManager(conn Database, by PartitionBy) RangePartitionManager {
	return RangePartitionManager{
		conn: conn,
		by:   by,
	}
}

const createPartitionTemplate = `CREATE TABLE IF NOT EXISTS ? PARTITION OF ? FOR VALUES FROM (?) TO (?);`

func monthBoundaries(current time.Time) (time.Time, time.Time) {
	start := time.Date(current.Year(), current.Month(), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	return start, end
}

func yearBoundaries(current time.Time) (time.Time, time.Time) {
	start := time.Date(current.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(1, 0, 0)

	return start, end
}

func monthPartitionId(currentTime time.Time) string {
	return fmt.Sprintf("%d_%02d", currentTime.Year(), currentTime.Month())
}

func yearPartitionId(currentTime time.Time) string {
	return fmt.Sprintf("%d", currentTime.Year())
}

func (pm *RangePartitionManager) getParameters(currentTime time.Time) (params, error) {
	var p params

	switch pm.by {
	case PartitionByMonth:
		p.id = monthPartitionId(currentTime)
	case PartitionByYear:
		p.id = yearPartitionId(currentTime)
	default:
		return p, errors.Errorf("unknown partition by: %d", pm.by)
	}

	p.success = p.id != pm.lastId
	if !p.success {
		return p, nil
	}

	switch pm.by {
	case PartitionByMonth:
		p.start, p.end = monthBoundaries(currentTime)
	case PartitionByYear:
		p.start, p.end = yearBoundaries(currentTime)
	default:
		return p, errors.Errorf("unknown partition by: %d", pm.by)
	}

	return p, nil
}

// CreatePartition -
func (pm *RangePartitionManager) CreatePartition(ctx context.Context, currentTime time.Time, tableName string) error {
	p, err := pm.getParameters(currentTime)
	if err != nil {
		return err
	}
	if !p.success {
		return nil
	}

	partitionName := fmt.Sprintf("%s_%s", tableName, p.id)
	if _, err := pm.conn.Exec(
		ctx,
		createPartitionTemplate,
		bun.Ident(partitionName),
		bun.Ident(tableName),
		p.start.Format(time.RFC3339Nano),
		p.end.Format(time.RFC3339Nano),
	); err != nil {
		return err
	}

	pm.lastId = p.id
	return nil
}

// CreatePartitions -
func (pm *RangePartitionManager) CreatePartitions(ctx context.Context, currentTime time.Time, tableNames ...string) error {
	p, err := pm.getParameters(currentTime)
	if err != nil {
		return err
	}
	if !p.success {
		return nil
	}

	for _, tableName := range tableNames {
		partitionName := fmt.Sprintf("%s_%s", tableName, p.id)
		if _, err := pm.conn.Exec(
			ctx,
			createPartitionTemplate,
			bun.Ident(partitionName),
			bun.Ident(tableName),
			p.start.Format(time.RFC3339Nano),
			p.end.Format(time.RFC3339Nano),
		); err != nil {
			return err
		}
	}

	pm.lastId = p.id
	return nil
}
