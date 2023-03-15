# DipDup Go SDK

[![Tests](https://github.com/dipdup-net/metadata/workflows/Tests/badge.svg?)](https://github.com/dipdup-net/metadata/actions?query=workflow%3ATests)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

This library partially implements DipDup framework features and can be used for building Tezos indexers and dapps when performance and effective resource utilization are important.

## Packages

### `config`

DipDup YAML [configuration](https://docs.dipdup.net/config-file-reference) parser. You can validate config by `validate` tag from [validator package](https://github.com/go-playground/validator).

```go
import "github.com/dipdup-net/go-lib/config"

type MyConfig struct {
	config.Config `yaml:",inline"`
    // Custom field here
    Fields Fields `yaml:"fields" validate:"required"`
}

// Substitute - required by Configurable interface
func (c *MyConfig) Substitute() error {
    return nil
}

type Fields struct {
    First string `yaml:"first"`
}

var cfg MyConfig
if err := config.Parse("config.yaml", &cfg); err != nil {
    panic(err)
}
```

### `node`

Simple Tezos RPC API wrapper. Docs you can find [here](node/README.md)

### `database`

Managing DipDup database connection. Default interface contains custom method `Connect` and extends 3 interfaces `driver.Pinger`,  `StateRepository` and `io.Closer`.


```go
// Database -
type Database interface {
	Connect(ctx context.Context, cfg config.Database) error

	StateRepository
	
	driver.Pinger
	io.Closer
}

// StateRepository -
type StateRepository interface {
	State(name string) (State, error)
	UpdateState(state State) error
	CreateState(state State) error
	DeleteState(state State) error
}
```

where `State` structure is:

```go
// State -
type State struct {
	//nolint
	tableName struct{} `gorm:"-" pg:"dipdup_state" json:"-"`

	IndexName string `gorm:"primaryKey" pg:",pk" json:"index_name"`
	IndexType string `json:"index_type"`
	Hash      string `json:"hash,omitempty"`
	Level     uint64 `json:"level"`
	UpdatedAt int    `gorm:"autoUpdateTime"`
}
```

There are 2 default implementations of `Database` interface:
* `Gorm` - database connection via [gorm](https://gorm.io/)
* `PgGo` - database connection via [pg-go](https://pg.uptrace.dev/)

There is method `Wait` which waiting until database connection will be established.

Example of usage:

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

db := database.NewPgGo()
if err := db.Connect(ctx, cfg); err != nil {
	panic(err)
}
defer db.Close()

database.Wait(ctx, db, 5*time.Second)

var yourModel struct{}
conn := db.DB()
if err := conn.WithContext(ctx).Model(&yourModel).Limit(10).Select(); err != nil {
	panic(err)
}
```


### `tzkt`

TzKT API and Events wrapper.  
Read more about events and SignalR in the [doc](https://github.com/dipdup-net/go-lib/blob/master/tzkt/events/README.md)

```go
package main

import (
	"context"
	"log"

	"github.com/dipdup-net/go-lib/tzkt/data"
	"github.com/dipdup-net/go-lib/tzkt/events"
)

func main() {
	tzkt := events.NewTzKT(data.BaseEventsURL)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := tzkt.Connect(ctx); err != nil {
		log.Panic(err)
	}
	defer tzkt.Close()

	if err := tzkt.SubscribeToHead(); err != nil {
		panic(err)
	}

	if err := tzkt.SubscribeToBlocks(); err != nil {
		panic(err)
	}

	if err := tzkt.SubscribeToCycles(2); err != nil {
		panic(err)
	}

	if err := tzkt.SubscribeToAccounts("KT1K4EwTpbvYN9agJdjpyJm4ZZdhpUNKB3F6"); err != nil {
		panic(err)
	}

	if err := tzkt.SubscribeToBigMaps(nil, "KT1K4EwTpbvYN9agJdjpyJm4ZZdhpUNKB3F6", ""); err != nil {
		panic(err)
	}

	if err := tzkt.SubscribeToOperations("KT1K4EwTpbvYN9agJdjpyJm4ZZdhpUNKB3F6", data.KindTransaction); err != nil {
		panic(err)
	}

	if err := tzkt.SubscribeToTokenTransfers("", "", ""); err != nil {
		panic(err)
	}

	if err := tzkt.SubscribeToTokenBalances("", "", ""); err != nil {
		panic(err)
	}

	for msg := range tzkt.Listen() {
		switch msg.Type {
		case events.MessageTypeData:

			switch msg.Channel {
			case events.ChannelAccounts:
				items := msg.Body.([]data.Account)
				for _, item := range items {
					log.Println(item)
				}
			case events.ChannelBigMap:
				items := msg.Body.([]data.BigMapUpdate)
				for _, item := range items {
					log.Println(item)
				}
			case events.ChannelBlocks:
				items := msg.Body.([]data.Block)
				for _, item := range items {
					log.Println(item)
				}
			case events.ChannelHead:
				head := msg.Body.(data.Head)
				log.Println(head)
			case events.ChannelOperations:
				items := msg.Body.([]any)
				for _, item := range items {
					log.Println(item.(*data.Transaction))
				}
			case events.ChannelTokenBalances:
				items := msg.Body.([]data.TokenBalance)
				for _, item := range items {
					log.Println(item)
				}
			case events.ChannelTransfers:
				items := msg.Body.([]data.Transfer)
				for _, item := range items {
					log.Println(item)
				}
			case events.ChannelCycles:
				log.Println(msg.Body.(data.Cycle))
			}

		case events.MessageTypeReorg:
			log.Print("reorg")
		case events.MessageTypeState:
			log.Print("initialized")
		case events.MessageTypeSubscribed:
			log.Println("subscribed", msg)
		}
	}
}

```

Example usage of the API wrapper:

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

tzkt := api.New("url here")

head, err := tzkt.GetHead(ctx)
if err != nil {
	log.Panic(err)
}
```

### `hasura`

Go wrapper for Hasura metadata methods, read docs [here](hasura/README.md)
