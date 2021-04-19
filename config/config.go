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
	Path string `yaml:"path"`
	Kind string `yaml:"kind"`
}

// Validate -
func (db *Database) Validate() error {
	for _, valid := range []string{
		DBKindClickHouse, DBKindMysql, DBKindPostgres, DBKindSqlServer, DBKindSqlite,
	} {
		if valid == db.Kind {
			return nil
		}
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
