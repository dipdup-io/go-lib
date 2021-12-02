package config

import (
	"bufio"
	"bytes"
	"io"
	"os"

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
	Hasura      *Hasura               `yaml:"hasura" validate:"omitempty"`
	Prometheus  *Prometheus           `yaml:"prometheus" validate:"omitempty"`
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

	f, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "reading file %s", filename)
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	buffer := bytes.NewBuffer(make([]byte, 0))
	part := make([]byte, 1024)

	for {
		count, err := reader.Read(part)
		if err != nil {
			if err == io.EOF {
				return buffer, nil
			} else {
				return nil, err
			}
		}
		expanded, err := expandVariables(part[:count])
		if err != nil {
			return nil, err
		}
		if _, err := buffer.Write(expanded); err != nil {
			return nil, err
		}
	}
}
