package config

import (
	"fmt"
	"io/ioutil"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// Config
type Config struct {
	Version     string                `yaml:"version" validate:"required"`
	Database    Database              `yaml:"database" validate:"required"`
	DataSources map[string]DataSource `yaml:"datasources"`
	Contracts   map[string]Contract   `yaml:"contracts"`
	Hasura      Hasura                `yaml:"hasura"`
}

// Substitute -
func (c *Config) Substitute() error {
	return nil
}

// DataSource -
type DataSource struct {
	Kind string `yaml:"kind"`
	URL  string `yaml:"url" validate:"required,url"`
}

// Contracts -
type Contract struct {
	Address  string `yaml:"address" validate:"required,len=36"`
	TypeName string `yaml:"typename"`
}

// Database
type Database struct {
	Path       string `yaml:"path" validate:"required_if=Kind sqlite"`
	Kind       string `yaml:"kind" validate:"required,oneof=sqlite postgres mysql"`
	Host       string `yaml:"host" validate:"required_unless=Kind sqlite"`
	Port       int    `yaml:"port" validate:"required_unless=Kind sqlite,gt=-1,lt=65535"`
	User       string `yaml:"user" validate:"required_unless=Kind sqlite"`
	Password   string `yaml:"password" validate:"required_unless=Kind sqlite"`
	Database   string `yaml:"database" validate:"required_unless=Kind sqlite"`
	SchemaName string `yaml:"schema_name"`
}

// Hasura -
type Hasura struct {
	URL                string `yaml:"url" validate:"required,url"`
	Secret             string `yaml:"admin_secret" validate:"required"`
	RowsLimit          uint64 `yaml:"select_limit" validate:"gt=0,lt=1000"`
	EnableAggregations bool   `yaml:"allow_aggregation"`
	Rest               *bool  `yaml:"rest"`
}

// Load - load default config from `filename`
func Load(filename string) (*Config, error) {
	if filename == "" {
		return nil, errors.Errorf("you have to provide configuration filename")
	}

	src, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "reading file %s", filename)
	}

	expanded := expandEnv(string(src))

	var c Config
	if err := yaml.Unmarshal([]byte(expanded), &c); err != nil {
		return nil, err
	}

	if err := c.Substitute(); err != nil {
		return nil, errors.Wrap(err, "Substitute")
	}

	return &c, validator.New().Struct(c)
}

// Parse -
func Parse(filename string, output Configurable) error {
	if filename == "" {
		return fmt.Errorf("you have to provide configuration filename")
	}

	src, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("reading file %s error: %w", filename, err)
	}

	expanded := expandEnv(string(src))

	if err := yaml.Unmarshal([]byte(expanded), output); err != nil {
		return err
	}

	if err := output.Substitute(); err != nil {
		return err
	}
	return validator.New().Struct(output)
}
