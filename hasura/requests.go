package hasura

import (
	"os"

	"github.com/pkg/errors"
)

type Request struct {
	Type string      `json:"type"`
	Args interface{} `json:"args"`
}

type versionedRequest struct {
	Type    string      `json:"type"`
	Version int         `json:"int"`
	Args    interface{} `json:"args"`
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
	Sources          []Source          `json:"sources"`
	QueryCollections []QueryCollection `json:"query_collections,omitempty"`
	RestEndpoints    []interface{}     `json:"rest_endpoints,omitempty"`
}

func newMetadata(version int, sources []Source) *Metadata {
	return &Metadata{
		Version: version,
		Sources: sources,
	}
}

// Configuration -
type Configuration struct {
	ConnectionInfo ConnectionInfo `json:"connection_info"`
}

// ConnectionInfo -
type ConnectionInfo struct {
	UsePreparedStatements bool        `json:"use_prepared_statements"`
	IsolationLevel        string      `json:"isolation_level"`
	DatabaseUrl           DatabaseUrl `json:"database_url"`
}

// DatabaseUrl -
type DatabaseUrl string

// DatabaseUrlFromEnv -
type DatabaseUrlFromEnv struct {
	FromEnv string `json:"from_env"`
}

// UnmarshalJSON -
func (d *DatabaseUrl) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*d = DatabaseUrl(s)
		return nil
	}

	var fromEnv DatabaseUrlFromEnv
	if err := json.Unmarshal(data, &fromEnv); err != nil {
		return err
	}

	*d = DatabaseUrl(os.Getenv(fromEnv.FromEnv))
	return nil
}

// Source -
type Source struct {
	Name          string        `json:"name"`
	Kind          string        `json:"kind"`
	Tables        []Table       `json:"tables"`
	Configuration Configuration `json:"configuration"`
}

// Table -
type Table struct {
	ObjectRelationships []Relationship     `json:"object_relationships"`
	ArrayRelationships  []Relationship     `json:"array_relationships"`
	SelectPermissions   []SelectPermission `json:"select_permissions"`
	Configuration       TableConfiguration `json:"configuration"`
	Schema              TableSchema        `json:"table"`
}

// Relationship -
type Relationship struct {
	Table   PGTable           `json:"table"`
	Name    string            `json:"name"`
	Comment string            `json:"comment,omitempty"`
	Using   RelationshipUsing `json:"using"`
}

// RelationshipUsing -
type RelationshipUsing struct {
	Manual *ManualRelationship `json:"manual_configuration,omitempty"`
	FK     *FKRelationship     `json:"foreign_key_constraint_on,omitempty"`
}

// PGTable -
type PGTable struct {
	Name   string `json:"name"`
	Schema string `json:"schema"`
}

// UnmarshalJSON -
func (t *PGTable) UnmarshalJSON(data []byte) error {
	type buf PGTable
	if err := json.Unmarshal(data, (*buf)(t)); err == nil {
		return err
	}

	return json.Unmarshal(data, &t.Name)
}

// MarshalJSON -
func (t *PGTable) MarshalJSON() ([]byte, error) {
	if t.Schema == "" {
		return json.Marshal(t.Name)
	}

	type buf PGTable
	return json.Marshal((*buf)(t))
}

// ManualRelationship -
type ManualRelationship struct {
	RemoteTable   PGTable           `json:"remote_table"`
	ColumnMapping map[string]string `json:"column_mapping"`
}

// FKRelationship -
type FKRelationship struct {
	Table  string `json:"table"`
	Column string `json:"column"`
}

func newMetadataTable(name, schema string) Table {
	return Table{
		ObjectRelationships: make([]Relationship, 0),
		ArrayRelationships:  make([]Relationship, 0),
		SelectPermissions:   make([]SelectPermission, 0),
		Schema: TableSchema{
			Name:   name,
			Schema: schema,
		},
	}
}

// TableConfiguration -
type TableConfiguration struct {
	Comment           *string           `json:"comment"`
	CustomRootFields  map[string]string `json:"custom_root_fields"`
	CustomColumnNames map[string]string `json:"custom_column_names"`
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
