# Go wrapper for Hasura metadata methods

Golang's library for registration data source in Hasura service

## Configuration

```yaml
hasura:
  url: string             # hasura url, required
  admin_secret: string    # required
  select_limit: int       
  allow_aggregation: bool
  source:
    name: string                  # name of data source, required. For more info, [hasura docs](https://hasura.io/docs/latest/api-reference/metadata-api/source/#metadata-pg-add-source-syntax).
    database_host: string         # host of datasource, if omitted, used host from database config
    use_prepared_statements: bool # if set to true the server prepares statement before executing on the source database (default: false)
    isolation_level: bool         # The transaction isolation level in which the queries made to the source will be run with (options: read-committed | repeatable-read | serializable) (default: read-committed)
  add_source: bool        # should data source be added?
  rest: bool              # should REST endpoints be created?
```
