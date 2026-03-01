package config

import (
	"bytes"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// Config
type Config struct {
	Version     string                `validate:"required"          yaml:"version"`
	Database    Database              `validate:"required"          yaml:"database"`
	DataSources map[string]DataSource `yaml:"datasources,omitempty"`
	Contracts   map[string]Contract   `yaml:"contracts,omitempty"`
	Hasura      *Hasura               `validate:"omitempty"         yaml:"hasura,omitempty"`
	Prometheus  *Prometheus           `validate:"omitempty"         yaml:"prometheus,omitempty"`
}

// Substitute -
func (c *Config) Substitute() error {
	return nil
}

// DataSource -
type DataSource struct {
	Kind              string       `yaml:"kind"`
	URL               string       `validate:"required,url"    yaml:"url"`
	Credentials       *Credentials `validate:"omitempty"       yaml:"credentials,omitempty"`
	Timeout           uint         `validate:"omitempty"       yaml:"timeout"`
	RequestsPerSecond int          `validate:"omitempty,min=1" yaml:"rps"`
}

// Contracts -
type Contract struct {
	Address  string `validate:"required,len=36" yaml:"address"`
	TypeName string `yaml:"typename"`
}

// Database
type Database struct {
	Path                   string `yaml:"path,omitempty"`
	Kind                   string `validate:"required,oneof=sqlite postgres mysql clickhouse elasticsearch" yaml:"kind"`
	Host                   string `validate:"required_with=Port User Database"                              yaml:"host"`
	Port                   int    `validate:"required_with=Host User Database,gt=-1,lt=65535"               yaml:"port"`
	User                   string `validate:"required_with=Host Port Database"                              yaml:"user"`
	Password               string `yaml:"password"`
	Database               string `validate:"required_with=Host Port User"                                  yaml:"database"`
	SchemaName             string `yaml:"schema_name"`
	ApplicationName        string `yaml:"application_name"`
	MaxOpenConnections     int    `yaml:"max_open_connections"`
	MaxIdleConnections     int    `yaml:"max_idle_connections"`
	MaxLifetimeConnections int    `yaml:"max_lifetime_connections"`
}

// Hasura -
type Hasura struct {
	URL                string        `validate:"required,url"  yaml:"url"`
	Secret             string        `validate:"required"      yaml:"admin_secret"`
	RowsLimit          uint64        `validate:"gt=0"          yaml:"select_limit"`
	EnableAggregations bool          `yaml:"allow_aggregation"`
	Source             *HasuraSource `yaml:"source"`
	Rest               *bool         `yaml:"rest"`
	UnauthorizedRole   string        `yaml:"unauthorized_role"`
}

type HasuraSource struct {
	Name                  string `validate:"required"            yaml:"name"`
	DatabaseHost          string `yaml:"database_host"`
	UsePreparedStatements bool   `yaml:"use_prepared_statements"`
	IsolationLevel        string `yaml:"isolation_level"`
}

// UnmarshalYAML -
func (h *Hasura) UnmarshalYAML(unmarshal func(interface{}) error) error {
	h.Source = &HasuraSource{
		Name:                  "default",
		UsePreparedStatements: false,
		IsolationLevel:        "read-committed",
	}

	h.UnauthorizedRole = "user"

	type plain Hasura
	return unmarshal((*plain)(h))
}

// Prometheus -
type Prometheus struct {
	URL string `validate:"required" yaml:"url"`
}

// Load - load default config from `filename`
func Load(filename string) (*Config, error) {
	buf, err := readFile(filename)
	if err != nil {
		return nil, err
	}
	var c Config
	if err := yaml.NewDecoder(buf).Decode(&c); err != nil {
		return nil, err
	}

	if err := c.Substitute(); err != nil {
		return nil, errors.Wrap(err, "Substitute")
	}

	return &c, validator.New().Struct(c)
}

// Parse -
func Parse(filename string, output Configurable) error {
	buf, err := readFile(filename)
	if err != nil {
		return err
	}

	if err := yaml.NewDecoder(buf).Decode(output); err != nil {
		return err
	}

	if err := output.Substitute(); err != nil {
		return err
	}
	return validator.New().Struct(output)
}

// ParseWithValidator - parse config with custom validator. If validator is nil validation will be skipped.
func ParseWithValidator(filename string, val *validator.Validate, output Configurable) error {
	buf, err := readFile(filename)
	if err != nil {
		return err
	}

	if err := yaml.NewDecoder(buf).Decode(output); err != nil {
		return err
	}

	if err := output.Substitute(); err != nil {
		return err
	}
	if val != nil {
		return val.Struct(output)
	}
	return nil
}

func readFile(filename string) (*bytes.Buffer, error) {
	if filename == "" {
		return nil, errors.Errorf("you have to provide configuration filename")
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "reading file %s", filename)
	}
	expanded, err := expandVariables(data)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(expanded), nil
}
