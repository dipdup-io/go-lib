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
	Version     string                `yaml:"version" validate:"required"`
	Database    Database              `yaml:"database" validate:"required"`
	DataSources map[string]DataSource `yaml:"datasources,omitempty"`
	Contracts   map[string]Contract   `yaml:"contracts,omitempty"`
	Hasura      *Hasura               `yaml:"hasura,omitempty" validate:"omitempty"`
	Prometheus  *Prometheus           `yaml:"prometheus,omitempty" validate:"omitempty"`
}

// Substitute -
func (c *Config) Substitute() error {
	return nil
}

// DataSource -
type DataSource struct {
	Kind              string       `yaml:"kind"`
	URL               string       `yaml:"url" validate:"required,url"`
	Credentials       *Credentials `yaml:"credentials,omitempty" validate:"omitempty"`
	Timeout           uint         `yaml:"timeout" validate:"omitempty"`
	RequestsPerSecond int          `yaml:"rps" validate:"omitempty,min=1"`
}

// Contracts -
type Contract struct {
	Address  string `yaml:"address" validate:"required,len=36"`
	TypeName string `yaml:"typename"`
}

// Database
type Database struct {
	Path            string `yaml:"path,omitempty"`
	Kind            string `yaml:"kind" validate:"required,oneof=sqlite postgres mysql clickhouse elasticsearch"`
	Host            string `yaml:"host" validate:"required_with=Port User Database"`
	Port            int    `yaml:"port" validate:"required_with=Host User Database,gt=-1,lt=65535"`
	User            string `yaml:"user" validate:"required_with=Host Port Database"`
	Password        string `yaml:"password"`
	Database        string `yaml:"database" validate:"required_with=Host Port User"`
	SchemaName      string `yaml:"schema_name"`
	ApplicationName string `yaml:"application_name"`
}

// Hasura -
type Hasura struct {
	URL                string        `yaml:"url" validate:"required,url"`
	Secret             string        `yaml:"admin_secret" validate:"required"`
	RowsLimit          uint64        `yaml:"select_limit" validate:"gt=0"`
	EnableAggregations bool          `yaml:"allow_aggregation"`
	Source             *HasuraSource `yaml:"source"`
	Rest               *bool         `yaml:"rest"`
	UnauthorizedRole   string        `yaml:"unauthorized_role"`
}

type HasuraSource struct {
	Name                  string `yaml:"name" validate:"required"`
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
	URL string `yaml:"url" validate:"required"`
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
