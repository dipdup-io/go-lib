# tzkt

[TzKT](https://tzkt.io/) REST API client and real-time events (SignalR) client for Tezos blockchain data.

```bash
go get github.com/dipdup-io/go-lib/tzkt
```

Full TzKT API documentation: [api.tzkt.io](https://api.tzkt.io/)

## REST API

### Creating a client

```go
import "github.com/dipdup-io/go-lib/tzkt/api"

tzkt := api.New("https://api.tzkt.io")
// or for a custom network:
tzkt := api.New("https://api.ghostnet.tzkt.io")
```

### Blocks

```go
head, err := tzkt.GetHead(ctx)
block, err := tzkt.GetBlock(ctx, 1234567)
blocks, err := tzkt.GetBlocks(ctx, api.BlocksFilter{Level: api.Int64Filter{Gt: 1000000}})
```

### Operations

```go
// Transactions for a contract
ops, err := tzkt.GetTransactions(ctx, api.TransactionsFilter{
    Target:  api.StringFilter{Eq: "KT1..."},
    Limit:   100,
})

// Originations
origs, err := tzkt.GetOriginations(ctx, api.OriginationsFilter{})

// Delegations
dels, err := tzkt.GetDelegations(ctx, api.DelegationsFilter{})
```

### Contracts and storage

```go
// Contract info
contract, err := tzkt.GetContract(ctx, "KT1...")

// Current storage
storage, err := tzkt.GetContractStorage(ctx, "KT1...")
```

### Big maps

```go
bigmap, err := tzkt.GetBigMap(ctx, 1234)
keys, err := tzkt.GetBigMapKeys(ctx, 1234, api.BigMapKeysFilter{Active: true})
updates, err := tzkt.GetBigMapUpdates(ctx, api.BigMapUpdatesFilter{BigMap: api.Int64Filter{Eq: 1234}})
```

### Tokens

```go
tokens, err := tzkt.GetTokens(ctx, api.TokensFilter{})
balances, err := tzkt.GetTokenBalances(ctx, api.TokenBalancesFilter{
    Account: api.StringFilter{Eq: "tz1..."},
})
transfers, err := tzkt.GetTokenTransfers(ctx, api.TokenTransfersFilter{})
```

### Delegates

```go
delegate, err := tzkt.GetDelegate(ctx, "tz1...")
delegates, err := tzkt.GetDelegates(ctx, api.DelegatesFilter{Active: true})
```

---

## Real-time events

Subscribe to live blockchain data via TzKT's SignalR WebSocket API.

### Connecting and subscribing

```go
import (
    "github.com/dipdup-io/go-lib/tzkt/data"
    "github.com/dipdup-io/go-lib/tzkt/events"
)

tzkt := events.NewTzKT(data.BaseEventsURL)

ctx, cancel := context.WithCancel(context.Background())
defer cancel()

if err := tzkt.Connect(ctx); err != nil {
    panic(err)
}
defer tzkt.Close()

// Subscribe to channels
tzkt.SubscribeToHead()
tzkt.SubscribeToBlocks()
tzkt.SubscribeToOperations("KT1...", data.KindTransaction)
tzkt.SubscribeToBigMaps(nil, "KT1...", "ledger")
tzkt.SubscribeToTokenTransfers("", "KT1...", "")
tzkt.SubscribeToTokenBalances("", "KT1...", "")
tzkt.SubscribeToAccounts("tz1...")
tzkt.SubscribeToCycles(2)
```

### Handling messages

```go
for msg := range tzkt.Listen() {
    switch msg.Type {
    case events.MessageTypeState:
        log.Println("subscription confirmed, current level:", msg.State)

    case events.MessageTypeData:
        switch msg.Channel {
        case events.ChannelHead:
            head := msg.Body.(data.Head)
        case events.ChannelBlocks:
            blocks := msg.Body.([]data.Block)
        case events.ChannelOperations:
            ops := msg.Body.([]any) // cast individual items to *data.Transaction etc.
        case events.ChannelBigMap:
            updates := msg.Body.([]data.BigMapUpdate)
        case events.ChannelTransfers:
            transfers := msg.Body.([]data.Transfer)
        case events.ChannelTokenBalances:
            balances := msg.Body.([]data.TokenBalance)
        case events.ChannelAccounts:
            accounts := msg.Body.([]data.Account)
        case events.ChannelCycles:
            cycle := msg.Body.(data.Cycle)
        }

    case events.MessageTypeReorg:
        log.Println("chain reorg at level", msg.State)
    }
}
```

### Message structure

```go
type Message struct {
    Channel string      // "head" | "blocks" | "operations" | "bigmaps" | ...
    Type    MessageType // 0=state, 1=data, 2=reorg
    State   uint64      // current chain level
    Body    any         // typed payload depending on Channel
}
```

### Available subscriptions

| Method | Channel | Body type |
|--------|---------|-----------|
| `SubscribeToHead()` | `head` | `data.Head` |
| `SubscribeToBlocks()` | `blocks` | `[]data.Block` |
| `SubscribeToOperations(addr, kind)` | `operations` | `[]any` |
| `SubscribeToBigMaps(ptr, addr, path)` | `bigmaps` | `[]data.BigMapUpdate` |
| `SubscribeToTokenTransfers(from, to, contract)` | `transfers` | `[]data.Transfer` |
| `SubscribeToTokenBalances(account, contract, token)` | `tokenbalances` | `[]data.TokenBalance` |
| `SubscribeToAccounts(addr)` | `accounts` | `[]data.Account` |
| `SubscribeToCycles(depth)` | `cycles` | `data.Cycle` |

### Operation kinds

```go
data.KindTransaction
data.KindOrigination
data.KindDelegation
data.KindReveal
data.KindEndorsement
// and more in the data package
```

---

## Low-level SignalR client

If you need to build a custom WebSocket client or reuse the SignalR transport in another package:

```go
import "github.com/dipdup-io/go-lib/tzkt/events/signalr"

client := signalr.NewSignalR("https://api.tzkt.io/v1/ws")
if err := client.Connect(ctx); err != nil {
    panic(err)
}
```
