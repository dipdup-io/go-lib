# go-lib
General instruments for DipDup in golang

## Packages

* `cmdline` - parser for default dipdup keys

```go
import "github.com/dipdup-net/go-lib/cmdline"

args := cmdline.Parse()
if args.Help {
	return
}
```

* `config`

```go
import "github.com/dipdup-net/go-lib/config"

type MyConfig struct {
	config.Config `yaml:",inline"`
    // Custom field here
    Fields Fields `yaml:"fields"`
}

// Validate - required by Configurable interface
func (c *MyConfig) Validate() error {
    return c.Fields.Validate() // if needed
}

// Substitute - required by Configurable interface
func (c *MyConfig) Substitute() error {
    return nil
}

type Fields struct {
    First string `yaml:"first"`
}

// Validate -
func (f *Fields) Validate() error {
    return nil
}

var cfg MyConfig
if err := config.Parse("config.yaml", &cfg); err != nil {
    panic(err)
}
```

* `node` - package for accessing to tezos node

```go
import "github.com/dipdup-net/go-lib/node"

rpc := node.NewNodeRPC(url, node.WithTimeout(timeout))
```

* `state` - package with DipDup state model

```go
import "github.com/dipdup-net/go-lib/state"

s := state.State{}
```

`State` structure is 
```go
// State -
type State struct {
	IndexName string `gorm:"primaryKey"`
	IndexType string
	Hash      string
	Level     uint64
}
```

* `tzkt` - package with API and Events wrapper for TzKT.

You can find docs on Events wrapper [here](tzkt/events/README.md).

Example usage of events

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

	if err := tzkt.SubscribeToBlocks(); err != nil {
		log.Panic(err)
	}

	for msg := range tzkt.Listen() {
		log.Println(msg)
	}
}

```

Example usage of API wrapper

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