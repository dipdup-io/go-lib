package config

import (
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// Config
type Config struct {
	Version     string                `yaml:"version"`
	Database    Database              `yaml:"database"`
	DataSources map[string]DataSource `yaml:"datasources"`
	Contracts   map[string]Contract   `yaml:"contracts"`
	Hasura      Hasura                `yaml:"hasura"`
}

// Validate -
func (c *Config) Validate() error {
	return c.Database.Validate()
}

// Substitute -
func (c *Config) Substitute() error {
	return nil
}

// DataSource -
type DataSource struct {
	Kind string `yaml:"kind"`
	URL  string `yaml:"url"`
}

// Contracts -
type Contract struct {
	Address  string `yaml:"address"`
	TypeName string `yaml:"typename"`
}

// Database
type Database struct {
	Path       string `yaml:"path"`
	Kind       string `yaml:"kind"`
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	Database   string `yaml:"database"`
	SchemaName string `yaml:"schema_name"`
}

// Hasura -
type Hasura struct {
	URL                string `yaml:"url"`
	Secret             string `yaml:"admin_secret"`
	RowsLimit          uint64 `yaml:"select_limit"`
	EnableAggregations bool   `yaml:"allow_aggregation"`
	Rest               *bool  `yaml:"rest"`
}

// Validate -
func (db *Database) Validate() error {
	if db.Kind == DBKindSqlite {
		if db.Path == "" {
			return errors.Wrap(ErrMissingField, "Path")
		}
		return nil
	} else if db.Kind == DBKindPostgres || db.Kind == DBKindMysql {
		if db.Host == "" {
			return errors.Wrap(ErrMissingField, "Host")
		}
		if db.Port == 0 {
			return errors.Wrap(ErrMissingField, "Port")
		}
		if db.User == "" {
			return errors.Wrap(ErrMissingField, "User")
		}
		if db.Password == "" {
			return errors.Wrap(ErrMissingField, "Password")
		}
		if db.Database == "" {
			return errors.Wrap(ErrMissingField, "Database")
		}
		return nil
	}
	return errors.Wrap(ErrUnsupportedDB, db.Kind)
}

// Load - load default config from `filename`
func Load(filename string) (*Config, error) {
	if filename == "" {
		return nil, fmt.Errorf("you have to provide configuration filename")
	}

	src, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("reading file %s error: %w", filename, err)
	}

	expanded := expandEnv(string(src))

	var c Config
	if err := yaml.Unmarshal([]byte(expanded), &c); err != nil {
		return nil, err
	}

	return &c, c.Substitute()
}

// LoadAndValidate - load config from `filename` and validate it
func LoadAndValidate(filename string) (*Config, error) {
	cfg, err := Load(filename)
	if err != nil {
		return nil, err
	}
	return cfg, cfg.Validate()
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
	return output.Validate()
}
