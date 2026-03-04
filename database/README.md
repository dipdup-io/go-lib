# database

PostgreSQL connection management for DipDup indexers. Built on [bun](https://bun.uptrace.dev/) ORM and [pgx](https://github.com/jackc/pgx) connection pool.

```bash
go get github.com/dipdup-io/go-lib/database
```

## Quickstart

```go
import (
    "github.com/dipdup-io/go-lib/config"
    "github.com/dipdup-io/go-lib/database"
)

db := database.NewBun()
if err := db.Connect(ctx, cfg.Database); err != nil {
    panic(err)
}
defer db.Close()

// Wait until the database is reachable (useful on startup)
database.Wait(ctx, db, 5*time.Second)

// Access the underlying bun.DB for queries
conn := db.DB()
```

## `Database` interface

All functionality is exposed through this interface, making it easy to mock in tests:

```go
type Database interface {
    Connect(ctx context.Context, cfg config.Database) error
    Exec(ctx context.Context, query string, args ...any) (int64, error)
    CreateTable(ctx context.Context, model any, opts ...CreateTableOption) error

    StateRepository   // index state CRUD
    SchemeCommenter   // table/column PostgreSQL comments

    driver.Pinger
    io.Closer
}
```

## State management

`State` tracks the sync progress of each index by name and level:

```go
type State struct {
    IndexName string `bun:",pk"`
    IndexType string
    Hash      string
    Level     uint64
    UpdatedAt int
}
```

```go
// Create initial state
if err := db.CreateState(ctx, &database.State{
    IndexName: "my_index",
    IndexType: "operation",
    Level:     0,
}); err != nil {
    panic(err)
}

// Read
state, err := db.State(ctx, "my_index")

// Update after processing a block
state.Level = newLevel
if err := db.UpdateState(ctx, state); err != nil {
    panic(err)
}
```

## Table management

### Creating tables

```go
type Transfer struct {
    bun.BaseModel `bun:"transfers"`
    ID       int64  `bun:",pk,autoincrement"`
    Sender   string `bun:"sender,notnull"`
    Receiver string `bun:"receiver,notnull"`
    Amount   int64  `bun:"amount,notnull"`
}

err := db.CreateTable(ctx, (*Transfer)(nil),
    database.WithIfNotExists(),
)
```

Available options:

| Option | Description |
|--------|-------------|
| `WithIfNotExists()` | Add `IF NOT EXISTS` clause |
| `WithPartitionBy(expr)` | Set `PARTITION BY` clause |
| `WithTemporary()` | Create a `TEMP` table |

### PostgreSQL table and column comments

Comments are set from struct tags and surfaced in Hasura's schema introspection:

```go
type Transfer struct {
    bun.BaseModel `bun:"transfers"         comment:"Token transfers"`
    ID       int64  `bun:",pk"             comment:"Primary key"`
    Sender   string `bun:"sender,notnull"  comment:"Sender address"`
}

// Reflect struct tags and apply all comments in one call
if err := database.MakeComments(ctx, db, (*Transfer)(nil)); err != nil {
    panic(err)
}
```

`MakeComments` recursively handles embedded structs. You can also set comments manually:

```go
db.MakeTableComment(ctx, "transfers", "Token transfers")
db.MakeColumnComment(ctx, "transfers", "sender", "Sender address")
```

## Connection pooling

`Bun` uses `pgxpool` internally. Pool parameters are read from `config.Database`:

| Config field | Default |
|---|---|
| `max_open_connections` | `4 × GOMAXPROCS` |
| `max_lifetime_connections` | 60 seconds |

All timestamps are forced to UTC on every connection via `pgtype` type registration, so no timezone surprises.

## Raw pool access

If you need to bypass bun and run raw pgx queries:

```go
pool := db.Pool() // *pgxpool.Pool
conn, err := pool.Acquire(ctx)
defer conn.Release()
```
