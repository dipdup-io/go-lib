# Golang code generator for TzKT API

Application generates Golang code for TzKT contract types. It requests JSON schema from TzKT API for your contract and generates code for processing TzKT API or events entities.

## Usage

To install binary

```bash
go get github.com/dipdup-net/go-lib/cmd/generator
```

```bash
generator -n my_contract -c KT1...
```

Args:

* `c` - contract address. For example, `KT1WxV6DDSFogKDg9DeAZZZr1HnVvKadpd3S`. Required if `f` is not set.
* `n` - your contract name. Optional. Default: `my_contract`.
* `u` - base TzKT API URL. Optional. Default: `https://api.tzkt.io/`.
* `o` - output directory. Optional. Default: current directory.
* `f` - path to JSON schema file. Required if `c` is not set.

## Output

Application creates directory according to contract name pointed in `n` command-line arg. It creates 3 files in the directory:

* `types.go` - default Tezos types
* `contract_types.go` - custom contract types
* `contract.go` - contract TzKT handler

## Usage of generated code

```go
package main

import (
    "context"
    "<YOUR_OUTPUT_DIRECTORY/atomex>"
    "log"
    "os"
    "signals"
)

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    atx := atomex.New("https://api.tzkt.io")
    if err := atx.Subscribe(ctx); err != nil {
        log.Panic(err)
    }

    initiates, err := atx.GetInitiate(ctx, atomex.Page{
        Limit: 2,
    })
    if err != nil {
        log.Panic(err)
    }
    log.Println(initiates)

    storage, err := atx.GetStorage(ctx)
    if err != nil {
        log.Panic(err)
    }
    log.Println(storage)

    signals := make(chan os.Signal, 1)
    signal.Notify(signals, os.Interrupt)

    for {
        select {
        case <-signals:
            cancel()
            if err := atx.Close(); err != nil {
                log.Panic(err)
            }
            close(signals)
            return
        case initiate := <-atx.InitiateEvents():
            log.Println(initiate)
        case add := <-atx.AddEvents():
            log.Println(add)
        case redeem := <-atx.RedeemEvents():
            log.Println(redeem)
        case refund := <-atx.RefundEvents():
            log.Println(refund)
        case update := <-atx.BigMap0Updates():
            log.Println(update)
        }
    }
}
```