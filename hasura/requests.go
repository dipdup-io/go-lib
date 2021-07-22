package hasura

import "github.com/pkg/errors"

type request struct {
	Type string      `json:"type"`
	Args interface{} `json:"args"`
}

// Permission -
type Permission struct {
	Columns   Columns     `json:"columns"`
	Limit     uint64      `json:"limit"`
	AllowAggs bool        `json:"allow_aggregations"`
	Filter    interface{} `json:"filter,omitempty"`
}

// Metadata -
type Metadata struct {
	Version          int               `json:"version"`
	Tables           []Table           `json:"tables"`
	QueryCollections []QueryCollection `json:"query_collections,omitempty"`
}

func newMetadata(version int, tables []Table) *Metadata {
	return &Metadata{
		Version: version,
		Tables:  tables,
	}
}

// Table -
type Table struct {
	ObjectRelationships []interface{}      `json:"object_relationships"`
	ArrayRelationships  []interface{}      `json:"array_relationships"`
	SelectPermissions   []SelectPermission `json:"select_permissions"`
	Schema              TableSchema        `json:"table"`
}

func newMetadataTable(name, schema string) Table {
	return Table{
		ObjectRelationships: make([]interface{}, 0),
		ArrayRelationships:  make([]interface{}, 0),
		SelectPermissions:   make([]SelectPermission, 0),
		Schema: TableSchema{
			Name:   name,
			Schema: schema,
		},
	}
}

// TableSchema -
type TableSchema struct {
	Schema string `json:"schema"`
	Name   string `json:"name"`
}

// SelectPermission -
type SelectPermission struct {
	Role       string     `json:"role"`
	Permission Permission `json:"permission"`
}

// Columns -
type Columns []string

// UnmarshalJSON -
func (columns *Columns) UnmarshalJSON(data []byte) error {
	var val interface{}
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}

	*columns = make(Columns, 0)
	switch typ := val.(type) {
	case string:
		*columns = append(*columns, typ)
	case []interface{}:
		for i := range typ {
			if s, ok := typ[i].(string); ok {
				*columns = append(*columns, s)
			}
		}
	default:
		return errors.Errorf("Invalid columns type: %T", typ)
	}
	return nil
}

// MarshalJSON -
func (columns Columns) MarshalJSON() ([]byte, error) {
	if len(columns) == 1 && columns[0] == "*" {
		return []byte(`"*"`), nil
	}
	s := []string(columns)
	return json.Marshal(s)
}

// QueryCollection -
type QueryCollection struct {
	Definition Definition `json:"definition"`
	Name       string     `json:"name"`
}

// Definition -
type Definition struct {
	Queries []Query `json:"queries"`
}

// Query -
type Query struct {
	Name           string `json:"name"`
	Query          string `json:"query,omitempty"`
	CollectionName string `json:"collection_name,omitempty"`
}
