package hasura

import (
	"reflect"
	"strings"
	"time"

	"github.com/dipdup-net/go-lib/config"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

// Create - creates hasura models
func Create(hasura config.Hasura, cfg config.Database, models ...interface{}) error {
	api := New(hasura.URL, hasura.Secret)

	log.Info("Waiting hasura is up...")
	for err := api.Health(); err != nil; err = api.Health() {
		time.Sleep(time.Second * 10)
	}

	metadata, err := Generate(cfg, models...)
	if err != nil {
		return err
	}

	log.Info("Fetching existing metadata...")
	export, err := api.ExportMetadata(metadata)
	if err != nil {
		return err
	}

	log.Info("Merging metadata...")
	tables := make(map[string]struct{})
	dataTables := metadata["tables"].([]interface{})
	for i := range dataTables {
		dataTable, ok := dataTables[i].(map[string]interface{})
		if !ok {
			continue
		}
		table, ok := dataTable["table"].(map[string]interface{})
		if !ok {
			continue
		}
		name := table["name"].(string)
		tables[name] = struct{}{}
	}

	for _, table := range export.Tables {
		tableData, ok := table["table"].(map[string]interface{})
		if !ok {
			continue
		}

		name := tableData["name"]
		if _, ok := tables[name.(string)]; !ok {
			dataTables = append(dataTables, table)
		}
	}

	metadata["tables"] = dataTables

	log.Info("Replacing metadata...")
	return api.ReplaceMetadata(metadata)
}

// Generate - creates hasura table structure in JSON from `models`. `models` should be pointer to your table models. `cfg` is DB config from YAML.
func Generate(cfg config.Database, models ...interface{}) (map[string]interface{}, error) {
	tables := make([]interface{}, 0)
	schema := getSchema(cfg)
	for _, model := range models {
		table, err := generateOne(schema, model)
		if err != nil {
			return nil, err
		}
		tables = append(tables, table.HasuraSchema)
	}

	return formatMetadata(tables), nil
}

type table struct {
	Name         string
	Schema       string
	Columns      []string
	HasuraSchema map[string]interface{}
}

func newTable(schema, name string) table {
	return table{
		Columns: make([]string, 0),
		Schema:  schema,
		Name:    name,
	}
}
func generateOne(schema string, model interface{}) (table, error) {
	value := reflect.ValueOf(model)
	if value.Kind() != reflect.Ptr {
		return table{}, errors.Errorf("Model has to be pointer")
	}
	typ := reflect.TypeOf(model)
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	t := newTable(schema, getTableName(value, typ))
	t.HasuraSchema = formatTable(t.Name, t.Schema)
	t.Columns = getColumns(typ)

	if p, ok := t.HasuraSchema["select_permissions"]; ok {
		t.HasuraSchema["select_permissions"] = append(p.([]interface{}), formatSelectPermissions(t.Columns...))
	} else {
		t.HasuraSchema["select_permissions"] = []interface{}{
			formatSelectPermissions(t.Columns...),
		}
	}
	t.HasuraSchema["object_relationships"] = []interface{}{}
	t.HasuraSchema["array_relationships"] = []interface{}{}

	return t, nil
}

func formatSelectPermissions(columns ...string) map[string]interface{} {
	return map[string]interface{}{
		"role": "user",
		"permission": map[string]interface{}{
			"columns":            columns,
			"filter":             map[string]interface{}{},
			"allow_aggregations": true,
		},
	}
}

func formatTable(name, schema string) map[string]interface{} {
	return map[string]interface{}{
		"table": map[string]interface{}{
			"schema": schema,
			"name":   name,
		},
		"object_relationships": []interface{}{},
		"array_relationships":  []interface{}{},
		"select_permissions":   []interface{}{},
	}
}

func formatMetadata(tables []interface{}) map[string]interface{} {
	return map[string]interface{}{
		"version": 2,
		"tables":  tables,
	}
}

func getTableName(value reflect.Value, typ reflect.Type) string {
	if _, ok := typ.MethodByName("TableName"); !ok {
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

// TODO: parsing schema from connection string
func getSchema(cfg config.Database) string {
	return "public"
}

func getColumns(typ reflect.Type) []string {
	columns := make([]string, 0)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if !field.Anonymous {
			tag := field.Tag.Get("gorm")
			if !strings.HasPrefix(tag, "-") {
				columns = append(columns, strcase.ToSnake(field.Name))
			}
		} else {
			cols := getColumns(field.Type)
			columns = append(columns, cols...)
		}
	}
	return columns
}
