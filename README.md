# DipDup Go SDK

[![Tests](https://github.com/dipdup-io/go-lib/workflows/Tests/badge.svg)](https://github.com/dipdup-io/go-lib/actions?query=workflow%3ATests)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Go SDK for building indexers and dapps with the [DipDup](https://dipdup.net) framework. Provides performance-oriented building blocks for blockchain data processing.

## Repository structure

The SDK is organized as a Go monorepo. Each subdirectory is an independent Go module with its own versioning:

| Module | Import path | Description |
|--------|-------------|-------------|
| `config` | `github.com/dipdup-io/go-lib/config` | DipDup YAML config parser with validation |
| `database` | `github.com/dipdup-io/go-lib/database` | PostgreSQL connection management via [bun](https://bun.uptrace.dev/) |
| `hasura` | `github.com/dipdup-io/go-lib/hasura` | [Hasura](https://hasura.io/) metadata API wrapper |
| `node` | `github.com/dipdup-io/go-lib/node` | Tezos RPC node client |
| `prometheus` | `github.com/dipdup-io/go-lib/prometheus` | Prometheus metrics service |
| `tools` | `github.com/dipdup-io/go-lib/tools` | Tezos data types, AST, crypto, forge/unforge utilities |
| `tzkt` | `github.com/dipdup-io/go-lib/tzkt` | [TzKT](https://tzkt.io/) REST API and SignalR events client |

Each module can be imported independently — you only pull in what you need.

## Modules

### `config`

DipDup YAML configuration parser. Supports environment variable substitution and validation via [`go-playground/validator`](https://github.com/go-playground/validator).

```go
import "github.com/dipdup-io/go-lib/config"

type MyConfig struct {
    config.Config `yaml:",inline"`
    Fields Fields `yaml:"fields" validate:"required"`
}

func (c *MyConfig) Substitute() error { return nil }

var cfg MyConfig
if err := config.Parse("config.yaml", &cfg); err != nil {
    panic(err)
}
```

---

### `database`

PostgreSQL connection management built on [bun](https://bun.uptrace.dev/). Provides a common `Database` interface with state tracking, connection waiting, and partition management.

```go
import "github.com/dipdup-io/go-lib/database"

db := database.NewBun()
if err := db.Connect(ctx, cfg.Database); err != nil {
    panic(err)
}
defer db.Close()

database.Wait(ctx, db, 5*time.Second)
```

---

### `hasura`

Wrapper for the [Hasura](https://hasura.io/) metadata API. Handles schema tracking, relationship creation, and permission setup programmatically.

See [`hasura/README.md`](hasura/README.md) for full documentation.

---

### `node`

Tezos node RPC client covering the full Tezos RPC specification. Includes typed interfaces for each API group (`BlockAPI`, `ChainAPI`, `ContextAPI`, etc.) to simplify mocking in tests.

```go
import "github.com/dipdup-io/go-lib/node"

rpc := node.NewMainRPC("https://rpc.tzkt.io/mainnet")
block, err := rpc.Block(ctx, "head")
```

See [`node/README.md`](node/README.md) for full documentation.

---

### `prometheus`

Prometheus metrics service that wraps `prometheus/client_golang`. Manages counters, histograms, and gauges by name, and runs an HTTP `/metrics` endpoint.

```go
import "github.com/dipdup-io/go-lib/prometheus"

svc := prometheus.NewService(cfg.Prometheus)
svc.RegisterCounter("indexer_operations_total", "Total operations indexed", "kind")
svc.Start()
defer svc.Close()

svc.IncrementCounter("indexer_operations_total", map[string]string{"kind": "transaction"})
```

---

### `tools`

Core Tezos utilities. Contains multiple sub-packages:

| Sub-package | Description |
|-------------|-------------|
| `tools/ast` | Michelson AST: parse, fold, encode/decode to JSON, Miguel, and storage formats |
| `tools/base` | Base node type used across AST |
| `tools/consts` | Tezos primitive constants |
| `tools/contract` | Contract parser and interface detection (FA1, FA1.2, FA2) |
| `tools/crypto` | Key parsing, signature verification, ed25519/secp256k1/p256 |
| `tools/encoding` | Base58Check encoding for addresses, keys, hashes |
| `tools/forge` | Binary forge/unforge for operations, Michelson values, and transactions |
| `tools/formatter` | Michelson source code formatter |
| `tools/tezerrors` | Typed Tezos error structures |
| `tools/tezgen` | Runtime types for `tezgen`-generated contract bindings |
| `tools/translator` | Converts Michelson values between different representations |
| `tools/types` | Shared type definitions |

---

### `tzkt`

[TzKT](https://tzkt.io/) API and real-time events client.

**REST API:**

```go
import "github.com/dipdup-io/go-lib/tzkt/api"

tzkt := api.New("https://api.tzkt.io")
head, err := tzkt.GetHead(ctx)
```

**Real-time events (SignalR):**

```go
import (
    "github.com/dipdup-io/go-lib/tzkt/data"
    "github.com/dipdup-io/go-lib/tzkt/events"
)

tzkt := events.NewTzKT(data.BaseEventsURL)
if err := tzkt.Connect(ctx); err != nil {
    panic(err)
}
defer tzkt.Close()

if err := tzkt.SubscribeToOperations("KT1...", data.KindTransaction); err != nil {
    panic(err)
}

for msg := range tzkt.Listen() {
    // handle msg
}
```

See [`tzkt/events/README.md`](tzkt/events/README.md) for all subscription types.

---

## Migration guide

### From `github.com/dipdup-net/go-lib` (v0.5.x and earlier)

The repository was split into independent modules and the organization name changed from `dipdup-net` to `dipdup-io`. Migration requires two steps.

#### 1. Update `go.mod`

Replace the single dependency with individual modules for what you actually use:

```diff
- require github.com/dipdup-net/go-lib v0.5.0
+ require (
+     github.com/dipdup-io/go-lib/config   v1.0.0
+     github.com/dipdup-io/go-lib/database v1.0.0
+     github.com/dipdup-io/go-lib/tzkt     v1.0.0
+     // add only what you need
+ )
```

Then run:

```bash
go mod tidy
```

#### 2. Update import paths

Replace all occurrences of the old import prefix in your `.go` files:

```bash
# macOS / BSD sed
find . -name "*.go" -exec sed -i '' \
  's|github.com/dipdup-net/go-lib|github.com/dipdup-io/go-lib|g' {} +

# GNU sed (Linux)
find . -name "*.go" -exec sed -i \
  's|github.com/dipdup-net/go-lib|github.com/dipdup-io/go-lib|g' {} +
```

The sub-package paths themselves are unchanged — only the module prefix changes:

| Before | After |
|--------|-------|
| `github.com/dipdup-net/go-lib/config` | `github.com/dipdup-io/go-lib/config` |
| `github.com/dipdup-net/go-lib/database` | `github.com/dipdup-io/go-lib/database` |
| `github.com/dipdup-net/go-lib/hasura` | `github.com/dipdup-io/go-lib/hasura` |
| `github.com/dipdup-net/go-lib/node` | `github.com/dipdup-io/go-lib/node` |
| `github.com/dipdup-net/go-lib/prometheus` | `github.com/dipdup-io/go-lib/prometheus` |
| `github.com/dipdup-net/go-lib/tools/...` | `github.com/dipdup-io/go-lib/tools/...` |
| `github.com/dipdup-net/go-lib/tzkt/...` | `github.com/dipdup-io/go-lib/tzkt/...` |

## Development

This repository uses [Go workspaces](https://go.dev/doc/tutorial/workspaces) for local development across all modules. No manual `replace` directives are needed:

```bash
git clone https://github.com/dipdup-io/go-lib
cd go-lib

# build all modules
go build ./...

# test all modules
go test ./...
```

## License

[MIT](LICENSE)
