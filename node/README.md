# node

Tezos RPC node client covering the full Tezos RPC specification. All API groups are exposed as interfaces for easy mocking in tests.

```bash
go get github.com/dipdup-io/go-lib/node
```

## Quickstart

```go
import "github.com/dipdup-io/go-lib/node"

rpc := node.NewMainRPC("https://rpc.tzkt.io/mainnet")

ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

block, err := rpc.Block(ctx, "head")
if err != nil {
    panic(err)
}
log.Printf("%+v", block)
```

## Constructors

| Constructor | Description |
|---|---|
| `NewRPC(url, chain)` | Full client with explicit chain ID |
| `NewMainRPC(url)` | Shorthand for mainnet (`chain = "main"`) |
| `NewMainBlockRPC(url)` | Block API only |
| `NewMainChainRPC(url)` | Chain API only |
| `NewMainContextRPC(url)` | Context API only |

## API groups

The `RPC` struct composes all groups. Each group can also be instantiated independently if you need only part of the API.

### Block API — `BlockAPI`

```go
block, err := rpc.Block(ctx, "head")      // latest block
block, err := rpc.Block(ctx, "1234567")   // by level
block, err := rpc.Block(ctx, "<hash>")    // by hash

ops,  err := rpc.BlockOperations(ctx, "head")
hash, err := rpc.BlockHash(ctx, "head")
```

### Chain API — `ChainAPI`

```go
chainID,    err := rpc.ChainID(ctx)
checkpoint, err := rpc.Checkpoint(ctx)
```

### Context API — `ContextAPI`

```go
storage,   err := rpc.ContractStorage(ctx, "head", "KT1...")
delegate,  err := rpc.Delegate(ctx, "head", "tz1...")
delegates, err := rpc.Delegates(ctx, "head")
constants, err := rpc.Constants(ctx, "head")
```

### Protocols API — `ProtocolsAPI`

```go
protocol,  err := rpc.Protocol(ctx, "head")
protocols, err := rpc.Protocols(ctx)
```

### Network API — `NetworkAPI`

```go
peers, err := rpc.Peers(ctx)
conns, err := rpc.Connections(ctx)
stat,  err := rpc.Stat(ctx)
```

### Inject API — `InjectAPI`

```go
hash, err := rpc.InjectOperation(ctx, signedBytes)
```

## Interfaces for mocking

Each API group has a corresponding interface. Use them in your own structs to enable substitution in tests:

```go
type MyIndexer struct {
    rpc node.BlockAPI
}
```

Available interfaces: `BlockAPI`, `ChainAPI`, `ContextAPI`, `ConfigAPI`, `GeneralAPI`, `ProtocolsAPI`, `NetworkAPI`, `InjectAPI`, `RpcAPI` (all-in-one).

Generate mocks with `mockgen`:

```bash
mockgen -source=node/interface.go -destination=mocks/node_mock.go
```
