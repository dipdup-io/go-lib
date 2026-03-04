# prometheus

Prometheus metrics service for DipDup indexers. Wraps [`prometheus/client_golang`](https://github.com/prometheus/client_golang) with a named-metric registry and an HTTP `/metrics` endpoint.

```bash
go get github.com/dipdup-io/go-lib/prometheus
```

## Quickstart

```go
import "github.com/dipdup-io/go-lib/prometheus"

svc := prometheus.NewService(cfg.Prometheus) // cfg.Prometheus.URL = ":9090"
svc.Start()
defer svc.Close()
```

If `cfg.Prometheus` is `nil` or `URL` is empty, `NewService` returns a no-op service — all metric calls become safe no-ops, which is useful in test environments.

## Metric types

### Counter

Monotonically increasing value. Use for counting events (blocks processed, errors, etc.).

```go
svc.RegisterCounter("indexer_blocks_total", "Total blocks indexed", "network", "status")

// Increment by 1
svc.IncrementCounter("indexer_blocks_total", map[string]string{
    "network": "mainnet",
    "status":  "success",
})

// Access the underlying *prometheus.CounterVec for custom operations
vec := svc.Counter("indexer_blocks_total")
vec.With(prometheus.Labels{"network": "mainnet", "status": "failed"}).Add(3)
```

### Histogram

Samples observations into configurable buckets. Use for latencies and sizes.

```go
svc.RegisterHistogram("indexer_operation_duration_seconds", "Operation processing time", "kind")

svc.AddHistogramValue("indexer_operation_duration_seconds",
    map[string]string{"kind": "transaction"},
    elapsed.Seconds(),
)
```

### Gauge

Value that can go up and down. Use for queue lengths, active connections, current level, etc.

```go
svc.RegisterGauge("indexer_current_level", "Current indexed level", "network")

svc.SetGaugeValue("indexer_current_level", map[string]string{"network": "mainnet"}, float64(level))
svc.IncGaugeValue("indexer_current_level", map[string]string{"network": "mainnet"})
svc.DecGaugeValue("indexer_current_level", map[string]string{"network": "mainnet"})
```

## Go build info metrics

Register the standard Go build info collector (Go version, module info):

```go
svc.RegisterGoBuildMetrics()
```

This adds the `go_build_info` metric automatically collected by the Go runtime.

## Config reference

```yaml
prometheus:
  url: ":9090"   # listen address; leave empty to disable
```

The `/metrics` endpoint is served at `http://<url>/metrics`.
