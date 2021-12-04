package hasura

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/dipdup-net/go-lib/config"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

const (
	allowedQueries = "allowed-queries"
)

func checkHealth(ctx context.Context, api *API) {
	log.Info("Waiting hasura is up and running")
	if err := api.Health(ctx); err != nil {
		return
	}
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := api.Health(ctx); err != nil {
				log.Warn(err)
				continue
			}
			return
		}
	}
}

// Create - creates hasura models
func Create(ctx context.Context, hasura *config.Hasura, cfg config.Database, views []string, models ...interface{}) error {
	if hasura == nil {
		return nil
	}
	api := New(hasura.URL, hasura.Secret)

	checkHealth(ctx, api)

	metadata, err := Generate(*hasura, cfg, models...)
	if err != nil {
		return err
	}

	log.Info("Fetching existing metadata...")
	export, err := api.ExportMetadata(ctx, metadata)
	if err != nil {
		return err
	}

	log.Info("Merging metadata...")
	tables := make(map[string]struct{})
	for i := range metadata.Tables {
		tables[metadata.Tables[i].Schema.Name] = struct{}{}
	}

	for _, table := range export.Tables {
		if _, ok := tables[table.Schema.Name]; !ok {
			metadata.Tables = append(metadata.Tables, table)
		}
	}

	if err := createQueryCollections(metadata); err != nil {
		return err
	}

	log.Info("Replacing metadata...")
	if err := api.ReplaceMetadata(ctx, metadata); err != nil {
		return err
	}

	if len(metadata.QueryCollections) > 0 && (hasura.Rest == nil || *hasura.Rest) {
		log.Info("Creating REST endpoints...")
		for _, query := range metadata.QueryCollections[0].Definition.Queries {
			if err := api.CreateRestEndpoint(ctx, query.Name, query.Name, query.Name, allowedQueries); err != nil {
				if e, ok := err.(APIError); !ok || !e.AlreadyExists() {
					return err
				}
			}
		}
	}

	log.Info("Tracking views...")
	for i := range views {
		if err := api.TrackTable(ctx, "public", views[i]); err != nil {
			if !strings.Contains(err.Error(), "view/table already tracked") {
				return err
			}
		}
		if err := api.DropSelectPermissions(ctx, views[i], "user"); err != nil {
			log.Warn(err)
		}
		if err := api.CreateSelectPermissions(ctx, views[i], "user", Permission{
			Limit:     hasura.RowsLimit,
			AllowAggs: hasura.EnableAggregations,
			Columns:   Columns{"*"},
			Filter:    map[string]interface{}{},
		}); err != nil {
			return err
		}
	}

	return nil
}

// Generate - creates hasura table structure in JSON from `models`. `models` should be pointer to your table models. `cfg` is DB config from YAML.
func Generate(hasura config.Hasura, cfg config.Database, models ...interface{}) (*Metadata, error) {
	tables := make([]Table, 0)
	schema := getSchema(cfg)
	for _, model := range models {
		table, err := generateOne(hasura, schema, model)
		if err != nil {
			return nil, err
		}
		tables = append(tables, table.HasuraSchema)
	}

	return newMetadata(2, tables), nil
}

type table struct {
	Name         string
	Schema       string
	Columns      []string
	HasuraSchema Table
}

func newTable(schema, name string) table {
	return table{
		Columns: make([]string, 0),
		Schema:  schema,
		Name:    name,
	}
}

func generateOne(hasura config.Hasura, schema string, model interface{}) (table, error) {
	value := reflect.ValueOf(model)
	if value.Kind() != reflect.Ptr {
		return table{}, errors.Errorf("Model has to be pointer")
	}
	typ := reflect.TypeOf(model)
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	t := newTable(schema, getTableName(value, typ))
	t.HasuraSchema = newMetadataTable(t.Name, t.Schema)
	t.Columns = getColumns(typ)

	t.HasuraSchema.SelectPermissions = append(t.HasuraSchema.SelectPermissions, formatSelectPermissions(hasura.RowsLimit, hasura.EnableAggregations, t.Columns...))

	return t, nil
}

func formatSelectPermissions(limit uint64, allowAggs bool, columns ...string) SelectPermission {
	if limit == 0 {
		limit = 10
	}
	return SelectPermission{
		Role: "user",
		Permission: Permission{
			Columns:   columns,
			Filter:    map[string]interface{}{},
			AllowAggs: allowAggs,
			Limit:     limit,
		},
	}
}

func getTableName(value reflect.Value, typ reflect.Type) string {
	if _, ok := typ.MethodByName("TableName"); !ok {
		if field, exists := typ.FieldByName("tableName"); exists {
			if tag := field.Tag.Get("pg"); tag != "" {
				if values := strings.Split(tag, ","); len(values) > 0 {
					return values[0]
				}
			}
		}
		return strcase.ToSnake(typ.Name())
	}
	res := value.MethodByName("TableName").Call([]reflect.Value{})
	if len(res) != 1 {
		return strcase.ToSnake(typ.Name())
	}
	if res[0].Kind() != reflect.String {
		return strcase.ToSnake(typ.Name())
	}
	return res[0].String()
}

func getSchema(cfg config.Database) string {
	return "public"
}

func getColumns(typ reflect.Type) []string {
	columns := make([]string, 0)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if !field.Anonymous {
			if tag := field.Tag.Get("gorm"); tag != "" {
				if !strings.HasPrefix(tag, "-") {
					columns = append(columns, strcase.ToSnake(field.Name))
				}
			} else if tag := field.Tag.Get("pg"); tag != "" {
				if !strings.HasPrefix(tag, "-") && field.Name != "tableName" {
					columns = append(columns, strcase.ToSnake(field.Name))
				}
			} else {
				columns = append(columns, strcase.ToSnake(field.Name))
			}
		} else {
			cols := getColumns(field.Type)
			columns = append(columns, cols...)
		}
	}
	return columns
}

func createQueryCollections(metadata *Metadata) error {
	if metadata == nil {
		return nil
	}

	files, err := ioutil.ReadDir("graphql")
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	queries := make([]Query, 0)
	for i := range files {
		name := files[i].Name()
		if !strings.HasSuffix(name, ".graphql") {
			continue
		}

		queryName := strings.TrimSuffix(name, ".graphql")

		f, err := os.Open(fmt.Sprintf("graphql/%s", name))
		if err != nil {
			return err
		}
		defer f.Close()

		data, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}

		queries = append(queries, Query{
			Name:  queryName,
			Query: string(data),
		})
	}

	if len(metadata.QueryCollections) > 0 && len(metadata.QueryCollections[0].Definition.Queries) > 0 {
		metadata.QueryCollections[0].Definition.Queries = mergeQueries(queries, metadata.QueryCollections[0].Definition.Queries)
	} else {
		metadata.QueryCollections = []QueryCollection{
			{
				Name: allowedQueries,
				Definition: Definition{
					Queries: queries,
				},
			},
		}
	}

	return nil
}

func mergeQueries(a []Query, b []Query) []Query {
	for i := range a {
		var found bool
		for j := range b {
			if b[j].Name == a[i].Name {
				found = true
				break
			}
		}

		if !found {
			b = append(b, a[i])
		}
	}
	return b
}
