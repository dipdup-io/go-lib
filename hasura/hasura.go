package hasura

import (
	"context"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/dipdup-net/go-lib/config"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"

	"github.com/rs/zerolog/log"
)

const (
	allowedQueries = "allowed-queries"
)

func checkHealth(ctx context.Context, api *API) {
	log.Info().Msg("Waiting hasura is up and running")
	if err := api.Health(ctx); err == nil {
		return
	}
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := api.Health(ctx); err != nil {
				log.Warn().Err(err).Msg("")
				continue
			}
			return
		}
	}
}

// GenerateArgs -
type GenerateArgs struct {
	Config               *config.Hasura  `validate:"required"`
	DatabaseConfig       config.Database `validate:"required"`
	Views                []string        `validate:"omitempty"`
	CustomConfigurations []Request       `validate:"omitempty"`
	Models               []any           `validate:"omitempty"`
}

// Create - creates hasura models
func Create(ctx context.Context, args GenerateArgs) error {
	if args.Config == nil {
		return nil
	}

	if err := validator.New().Struct(args); err != nil {
		return err
	}

	api := New(args.Config.URL, args.Config.Secret)

	checkHealth(ctx, api)

	if args.Config.Source != nil {
		log.Info().Msg("Adding source...")
		if err := api.AddSource(ctx, args.Config, args.DatabaseConfig); err != nil {
			return err
		}
	}

	metadata, err := Generate(*args.Config, args.DatabaseConfig, args.Models...)
	if err != nil {
		return err
	}

	log.Info().Msg("Fetching existing metadata...")
	export, err := api.ExportMetadata(ctx)
	if err != nil {
		return err
	}

	// Find our source in the existing metadata
	var selectedSource *Source = nil
	for idx := range export.Sources {
		if export.Sources[idx].Name == args.Config.Source.Name {
			selectedSource = &export.Sources[idx]
			break
		}
	}
	if selectedSource == nil {
		return errors.Errorf("Source '%s' not found on exported metadata", args.Config.Source.Name)
	}

	log.Info().Msg("Merging metadata...")
	// Clear tables
	// TODO: maybe instead replace tables by name.
	selectedSource.Tables = make([]Table, 0)
	// Insert generated tables
	selectedSource.Tables = append(selectedSource.Tables, metadata.Sources[0].Tables...)

	if err := createQueryCollections(&export); err != nil {
		return err
	}

	log.Info().Msg("Replacing metadata...")
	if err := api.ReplaceMetadata(ctx, &export); err != nil {
		return err
	}

	if len(export.QueryCollections) > 0 && (args.Config.Rest == nil || *args.Config.Rest) {
		log.Info().Msg("Creating REST endpoints...")
		for _, query := range export.QueryCollections[0].Definition.Queries {
			if err := api.CreateRestEndpoint(ctx, query.Name, query.Name, query.Name, allowedQueries); err != nil {
				if e, ok := err.(APIError); !ok || !e.AlreadyExists() {
					return err
				}
			}
		}
	}

	log.Info().Msg("Tracking views...")
	for i := range args.Views {
		if err := api.TrackTable(ctx, args.Views[i], args.Config.Source.Name); err != nil {
			if !strings.Contains(err.Error(), "view/table already tracked") {
				return err
			}
		}
		if err := api.DropSelectPermissions(ctx, args.Views[i], args.Config.Source.Name, args.Config.UnauthorizedRole); err != nil {
			if e, ok := err.(APIError); !ok || !e.PermissionDenied() {
				log.Warn().Err(err).Msg("")
			}
		}
		if err := api.CreateSelectPermissions(ctx, args.Views[i], args.Config.Source.Name, args.Config.UnauthorizedRole, Permission{
			Limit:     args.Config.RowsLimit,
			AllowAggs: args.Config.EnableAggregations,
			Columns:   Columns{"*"},
			Filter:    map[string]interface{}{},
		}); err != nil {
			return err
		}
	}

	log.Info().Msg("Running custom configurations...")
	for _, conf := range args.CustomConfigurations {
		if err := api.CustomConfiguration(ctx, conf); err != nil {
			log.Warn().Err(err).Msg("")
		}
	}

	return nil
}

// Generate - creates hasura table structure in JSON from `models`. `models` should be pointer to your table models. `cfg` is DB config from YAML.
func Generate(hasura config.Hasura, cfg config.Database, models ...interface{}) (*Metadata, error) {
	schema := getSchema(cfg)
	source := Source{
		Name:   hasura.Source.Name,
		Tables: make([]Table, 0),
	}
	for _, model := range models {
		table, err := generateOne(hasura, schema, model)
		if err != nil {
			return nil, err
		}
		source.Tables = append(source.Tables, table.HasuraSchema)
	}

	return newMetadata(3, []Source{source}), nil
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

	t.HasuraSchema.SelectPermissions = append(t.HasuraSchema.SelectPermissions, formatSelectPermissions(hasura.RowsLimit, hasura.EnableAggregations, hasura.UnauthorizedRole, t.Columns...))

	if err := getRelationships(&t.HasuraSchema, t.Name, typ); err != nil {
		return t, err
	}

	return t, nil
}

type rel struct {
	tableName   string
	remoteField string
	field       string
	name        string
	typ         string
	comment     string
}

func getRelationships(t *Table, name string, typ reflect.Type) error {
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("hasura")
		if tag == "" {
			continue
		}

		var r rel
		attrs := strings.Split(tag, ",")
		for i := range attrs {
			keyValue := strings.Split(attrs[i], ":")
			if len(keyValue) != 2 {
				return errors.Errorf("invalid hasura tag: %s", tag)
			}
			switch keyValue[0] {
			case "table":
				r.tableName = keyValue[1]
			case "field":
				r.field = keyValue[1]
			case "type":
				r.typ = keyValue[1]
			case "name":
				r.name = keyValue[1]
			case "comment":
				r.comment = keyValue[1]
			case "remote_field":
				r.remoteField = keyValue[1]
			}
		}

		relationship := Relationship{
			Table: PGTable{
				Name: name,
			},
			Name:    r.name,
			Comment: r.comment,
			Using: RelationshipUsing{
				Manual: &ManualRelationship{
					RemoteTable: PGTable{Name: r.tableName},
					ColumnMapping: map[string]string{
						r.field: r.remoteField,
					},
				},
			},
		}

		switch r.typ {
		case "oto":
			t.ObjectRelationships = append(t.ObjectRelationships, relationship)
		case "otm", "mtm":
			t.ArrayRelationships = append(t.ArrayRelationships, relationship)
		}
	}

	return nil
}

func formatSelectPermissions(limit uint64, allowAggs bool, role string, columns ...string) SelectPermission {
	if limit == 0 {
		limit = 10
	}
	return SelectPermission{
		Role: role,
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
			if tag := field.Tag.Get("bun"); tag != "" {
				if values := strings.Split(tag, ","); len(values) > 0 {
					return values[0]
				}
			}
		}
		return ToSnakeCase(typ.Name())
	}
	res := value.MethodByName("TableName").Call([]reflect.Value{})
	if len(res) != 1 {
		return ToSnakeCase(typ.Name())
	}
	if res[0].Kind() != reflect.String {
		return ToSnakeCase(typ.Name())
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
					columns = append(columns, ToSnakeCase(field.Name))
				}
			} else if tag := field.Tag.Get("pg"); tag != "" {
				if !strings.HasPrefix(tag, "-") && field.Name != "tableName" && !strings.Contains(tag, "rel:") {
					columns = append(columns, ToSnakeCase(field.Name))
				}
			} else if tag := field.Tag.Get("bun"); tag != "" {
				if !strings.HasPrefix(tag, "-") && field.Name != "BaseModel" && !strings.Contains(tag, "rel:") {
					columns = append(columns, ToSnakeCase(field.Name))
				}
			} else {
				columns = append(columns, ToSnakeCase(field.Name))
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

	files, err := os.ReadDir("graphql")
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

		data, err := io.ReadAll(f)
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

// ReadCustomConfigs -
func ReadCustomConfigs(ctx context.Context, database config.Database, hasuraConfigDir string) ([]Request, error) {
	files, err := os.ReadDir(hasuraConfigDir)
	if err != nil {
		return nil, err
	}

	custom_configs := make([]Request, 0)
	for i := range files {
		if files[i].IsDir() || strings.HasPrefix(files[i].Name(), ".") {
			continue
		}

		path := fmt.Sprintf("%s/%s", hasuraConfigDir, files[i].Name())
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		conf := Request{}

		err = json.Unmarshal(raw, &conf)
		if err != nil {
			return nil, err
		}
		custom_configs = append(custom_configs, conf)
	}

	return custom_configs, nil
}
