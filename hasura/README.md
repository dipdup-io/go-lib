# hasura

Programmatic Hasura metadata management for DipDup indexers. Automatically tracks tables, creates relationships, sets permissions, and optionally generates REST endpoints — all from your Go model structs.

```bash
go get github.com/dipdup-io/go-lib/hasura
```

## Quickstart

```go
import "github.com/dipdup-io/go-lib/hasura"

err := hasura.Create(ctx, cfg.Hasura, cfg.Database, []any{
    (*Transfer)(nil),
    (*Token)(nil),
    (*Account)(nil),
}, nil)
```

`Create` performs the full setup sequence:

1. Waits for Hasura to become healthy
2. Adds the PostgreSQL data source
3. Tracks all listed tables
4. Derives object/array relationships from bun `rel:` tags
5. Creates `select` permissions for the configured role
6. Optionally creates REST endpoints for tracked tables

## Configuration

```yaml
hasura:
  url: http://localhost:8080
  admin_secret: myadminsecret
  select_limit: 100           # default row limit for queries
  allow_aggregation: true     # enable aggregate queries
  unauthorized_role: user     # role used for anonymous access (default: "user")
  rest: true                  # create REST endpoints for each table

  source:
    name: default                    # data source name in Hasura
    database_host: ""                # override host (uses database.host if empty)
    use_prepared_statements: false
    isolation_level: read-committed  # read-committed | repeatable-read | serializable
```

## Relationships from struct tags

Relationships are derived automatically from bun ORM tags on your models:

```go
type Transfer struct {
    bun.BaseModel `bun:"transfers"`

    ID       int64    `bun:",pk,autoincrement"`
    SenderID int64    `bun:"sender_id"`
    Sender   *Account `bun:"rel:belongs-to,join:sender_id=id"`

    TokenID int64  `bun:"token_id"`
    Token   *Token `bun:"rel:belongs-to,join:token_id=id"`
}

type Account struct {
    bun.BaseModel `bun:"accounts"`
    ID        int64       `bun:",pk,autoincrement"`
    Transfers []*Transfer `bun:"rel:has-many,join:id=sender_id"`
}
```

Supported relations: `belongs-to` (object relationship) and `has-many` (array relationship).

## Generating metadata without applying

To inspect or export the metadata JSON instead of applying it:

```go
metadata, err := hasura.Generate(ctx, cfg.Hasura, cfg.Database, []any{
    (*Transfer)(nil),
    (*Token)(nil),
}, nil)
```

## REST endpoints

When `rest: true` is set, a REST endpoint is created for each tracked table at:

```
GET /api/rest/<table_name>
```

GraphQL query files from a `queries/` directory are also registered if present.

## Low-level API client

The `API` struct provides direct access to individual Hasura metadata API calls:

```go
api := hasura.NewAPI("http://localhost:8080", "myadminsecret")

healthy, err := api.Health(ctx)
err = api.TrackTable(ctx, sourceName, schema, tableName)
err = api.CreateSelectPermissions(ctx, sourceName, tableName, role, filter, limit, allowAggregation)
err = api.DropSelectPermissions(ctx, sourceName, tableName, role)
err = api.CreateRestEndpoint(ctx, name, url, query, collectionName)

meta, err := api.ExportMetadata(ctx)
err  = api.ReplaceMetadata(ctx, meta)
```
