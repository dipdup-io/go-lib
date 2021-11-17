# DipDup Go SDK

[![Tests](https://github.com/dipdup-net/metadata/workflows/Tests/badge.svg?)](https://github.com/dipdup-net/metadata/actions?query=workflow%3ATests)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

This library partially implements DipDup framework features and can be used for building Tezos indexers and dapps when performance and effective resource utilization are important.

## Packages

#### `cmdline`

Command line argument parser, compatible with [DipDup CLI](https://docs.dipdup.net/command-line).

```go
import "github.com/dipdup-net/go-lib/cmdline"

args := cmdline.Parse()
if args.Help {
	return
}
```

#### `config`

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

#### `node`

Simple Tezos RPC API wrapper.

```go
import "github.com/dipdup-net/go-lib/node"

rpc := node.NewNodeRPC(url, node.WithTimeout(timeout))
```

#### `state`

Managing DipDup index state.

```go
import "github.com/dipdup-net/go-lib/state"

s := state.State{}
```

where `State` structure is:
```go
// State -
type State struct {
	IndexName string `gorm:"primaryKey"`
	IndexType string
	Hash      string
	Level     uint64
	UpdatedAt int    `gorm:"autoUpdateTime"`
}
```

#### `tzkt`

TzKT API and Events wrapper.  
Read more about events and SignalR in the [doc](https://github.com/dipdup-net/go-lib/blob/master/tzkt/events/README.md)

```go
package main

import (
	"log"

	"github.com/dipdup-net/go-lib/tzkt/events"
)

func main() {
	tzkt := events.NewTzKT(events.BaseURL)
	if err := tzkt.Connect(); err != nil {
		log.Panic(err)
	}
	defer tzkt.Close()

	if err := tzkt.SubscribeToHead(); err != nil {
		panic(err)
	}

	if err := tzkt.SubscribeToBlocks(); err != nil {
		panic(err)
	}

	if err := tzkt.SubscribeToAccounts("KT1K4EwTpbvYN9agJdjpyJm4ZZdhpUNKB3F6"); err != nil {
		panic(err)
	}

	if err := tzkt.SubscribeToBigMaps(nil, "KT1K4EwTpbvYN9agJdjpyJm4ZZdhpUNKB3F6", ""); err != nil {
		panic(err)
	}

	if err := tzkt.SubscribeToOperations("KT1K4EwTpbvYN9agJdjpyJm4ZZdhpUNKB3F6", events.KindTransaction); err != nil {
		panic(err)
	}

	for msg := range tzkt.Listen() {
		switch msg.Type {
		case events.MessageTypeData:

			switch msg.Channel {
			case events.ChannelAccounts:
				items := msg.Body.([]events.Account)
				for _, item := range items {
					log.Println(item)
				}
			case events.ChannelBigMap:
				items := msg.Body.([]events.BigMapUpdate)
				for _, item := range items {
					log.Println(item)
				}
			case events.ChannelBlocks:
				items := msg.Body.([]events.Block)
				for _, item := range items {
					log.Println(item)
				}
			case events.ChannelHead:
				head := msg.Body.(events.Head)
				log.Println(head)
			case events.ChannelOperations:
				items := msg.Body.([]interface{})
				for _, item := range items {
					log.Println(item.(*events.Transaction))
				}
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
package main

import (
	"log"

	"github.com/dipdup-net/go-lib/tzkt/api"
)

func main() {
    tzkt := api.New("url here")
    
    head, err := tzkt.GetHead()
    if err != nil {
        log.Panic(err)
    }
    log.Println(head)
}
```