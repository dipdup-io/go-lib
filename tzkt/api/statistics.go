package api

import (
	"context"

	"github.com/dipdup-net/go-lib/tzkt/data"
)

// Statistics -
func (tzkt *API) Statistics(ctx context.Context, filters map[string]string) (stats []data.Statistics, err error) {
	err = tzkt.json(ctx, "/v1/statistics", filters, false, &stats)
	return
}

// StatisticsDaily -
func (tzkt *API) StatisticsDaily(ctx context.Context, filters map[string]string) (stats []data.Statistics, err error) {
	err = tzkt.json(ctx, "/v1/statistics/daily", filters, false, &stats)
	return
}

// StatisticsCyclic -
func (tzkt *API) StatisticsCyclic(ctx context.Context, filters map[string]string) (stats []data.Statistics, err error) {
	err = tzkt.json(ctx, "/v1/statistics/cyclic", filters, false, &stats)
	return
}

// StatisticsCurrent -
func (tzkt *API) StatisticsCurrent(ctx context.Context, filters map[string]string) (stats data.Statistics, err error) {
	err = tzkt.json(ctx, "/v1/statistics/current", filters, false, &stats)
	return
}
