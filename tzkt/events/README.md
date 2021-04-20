# TzKT events wrapper
Golang library for TzKT events API

## Usage 


### TzKT client

First of all, import library

```golang
import events "github.com/dipdup-net/go-lib/tzkt/events"
```

Then create `TzKT` client, connect to server and subscribe to channels.

```golang
tzkt := events.NewTzKT(events.BaseURL)
if err := tzkt.Connect(); err != nil {
    log.Panic(err)
}
defer tzkt.Close()

if err := tzkt.SubscribeToBlocks(); err != nil {
    log.Panic(err)
}
```

Now, you can listen message channel

```golang
for msg := range tzkt.Listen() {
    log.Println(msg)
}
```

Message is struct with fields:

```golang
type Message struct {
	Channel string       // is channel name: head, block or operations
	Type    MessageType  // is message type: 0, 1 or 2 (state, data, reorg)
	State   uint64       // is current level
	Body    interface{}  // is map or array of data depending of channel
}
```


### SignalR

If you want to write custom client or re-use SignalR in another package you can import only signalr and use it.


```golang
import "github.com/dipdup-net/go-lib/tzkt/events/signalr"

client := signalr.NewSignalR("https://api.tzkt.io/v1/events")
```