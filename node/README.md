# Tezos RPC client

The library realize almost all RPC methods of Tezos node.

## Usage

### Simple example

```go
rpc := node.NewRPC("https://rpc.tzkt.io/mainnet", "main")
ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
defer cancel()

block, err := rpc.Block(ctx, "head")
if err != nil {
	panic(err)
}
log.Printf("%##v", block)
```

You can use main RPC constructor where chain id set by default to `main`

```go
rpc := node.NewMainRPC("https://rpc.tzkt.io/mainnet")
```

### Usage certain API part

RPC struct contains some internal parts: `Chain`, `Block`, `Context`, `Config`, `General`, `Protocols`, `Inject` and `Network`. You can use it without creation of full RPC client.

```go
rpc := node.NewMainBlockRPC("https://rpc.tzkt.io/mainnet")
ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
defer cancel()

block, err := rpc.Block(ctx, "head")
if err != nil {
	panic(err)
}
log.Printf("%##v", block)
```

### Interfaces

For testing purpose RPC was wrapped by interfaces. Also each part of RPC has interface. You can mock it with code generation tools.

Interfaces list:
* `BlockAPI`
* `ChainAPI`
* `ContextAPI`
* `ConfigAPI`
* `GeneralAPI`
* `ProtocolsAPI`
* `NetworkAPI`
* `InjectAPI`