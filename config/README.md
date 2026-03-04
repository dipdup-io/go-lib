# config

YAML configuration parser for DipDup indexers with built-in validation and environment variable substitution.

```bash
go get github.com/dipdup-io/go-lib/config
```

## Overview

The package provides two usage modes:

- **Standard config** — embed `config.Config` into your own struct and use `Parse()`
- **Custom config** — define any struct, implement the `Configurable` interface, call `Parse()`

## Config file format

```yaml
version: "1.0"

database:
  kind: postgres      # postgres | sqlite | mysql | clickhouse | elasticsearch
  host: localhost
  port: 5432
  user: dipdup
  password: ${DB_PASSWORD}
  database: dipdup_db
  schema_name: public
  application_name: my-indexer
  max_open_connections: 10
  max_idle_connections: 5
  max_lifetime_connections: 60   # seconds

datasources:
  tzkt_mainnet:
    kind: tzkt
    url: https://api.tzkt.io
    rps: 10
    timeout: 30
    credentials:
      api_key:
        header: X-API-Key
        key: ${TZKT_API_KEY}

contracts:
  my_contract:
    address: KT1K4EwTpbvYN9agJdjpyJm4ZZdhpUNKB3F6
    typename: my_contract

hasura:
  url: http://localhost:8080
  admin_secret: ${HASURA_SECRET}
  select_limit: 100
  allow_aggregation: true
  source:
    name: default
    use_prepared_statements: false
    isolation_level: read-committed
  rest: true

prometheus:
  url: ":9090"
```

### Environment variable substitution

Any YAML value can reference an environment variable:

```yaml
password: ${DB_PASSWORD}            # required — fails if variable is not set
password: ${DB_PASSWORD:-secret}    # optional — uses "secret" if variable is not set
```

Substitution works recursively across all nested fields before YAML is decoded.

## Usage

### Extending the built-in config

```go
import "github.com/dipdup-io/go-lib/config"

type MyConfig struct {
    config.Config `yaml:",inline"`
    Indexer IndexerConfig `yaml:"indexer" validate:"required"`
}

type IndexerConfig struct {
    Workers int    `yaml:"workers" validate:"required,min=1"`
    Queue   string `yaml:"queue"   validate:"required"`
}

// Substitute is called after YAML decode, before validation.
// Use it to derive computed fields or validate cross-field logic.
func (c *MyConfig) Substitute() error {
    return nil
}

var cfg MyConfig
if err := config.Parse("config.yaml", &cfg); err != nil {
    panic(err)
}
```

### Custom validator

```go
import (
    "github.com/dipdup-io/go-lib/config"
    "github.com/go-playground/validator/v10"
)

v := validator.New()
v.RegisterValidation("myRule", myRuleFunc)

var cfg MyConfig
if err := config.ParseWithValidator("config.yaml", v, &cfg); err != nil {
    panic(err)
}
```

Pass `nil` as the validator to skip validation entirely.

## Types reference

### `Database`

| Field | Type | Description |
|-------|------|-------------|
| `kind` | string | `postgres`, `sqlite`, `mysql`, `clickhouse`, `elasticsearch` |
| `host` / `port` / `user` / `password` / `database` | string/int | Connection details |
| `path` | string | DSN string (alternative to individual fields) |
| `schema_name` | string | PostgreSQL search path |
| `application_name` | string | `application_name` parameter sent to Postgres |
| `max_open_connections` | int | Default: `4 × GOMAXPROCS` |
| `max_idle_connections` | int | |
| `max_lifetime_connections` | int | Seconds |

### `DataSource`

| Field | Type | Description |
|-------|------|-------------|
| `kind` | string | Arbitrary label (e.g. `tzkt`, `rpc`) |
| `url` | string | Required, must be a valid URL |
| `rps` | int | Max requests per second |
| `timeout` | uint | Request timeout in seconds |
| `credentials` | object | Optional `user` or `api_key` block |

### `Credentials`

```yaml
credentials:
  # Basic auth
  user:
    name: admin
    password: secret

  # API key header
  api_key:
    header: X-API-Key
    key: abc123
```

### `Alias[T]`

A generic helper for config fields that accept either a short name (string) or a full inline struct:

```yaml
# Short form — resolved by name
datasource: tzkt_mainnet

# Inline form
datasource:
  kind: tzkt
  url: https://api.tzkt.io
```

```go
type MyConfig struct {
    Source config.Alias[config.DataSource] `yaml:"datasource"`
}
```
