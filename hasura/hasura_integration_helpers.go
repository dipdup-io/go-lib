package hasura

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/dipdup-net/go-lib/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type IntegrationHelpers struct {
	api    *API
	config *config.Hasura
}

type ExpectedMetadata struct {
	Tables []ExpectedTable `yaml:"tables" validate:"required"`
}

type ExpectedTable struct {
	Name    string   `yaml:"name" validate:"reuired"`
	Columns []string `yaml:"columns" validate:"required"`
}

func NewIntegrationHelpers(config *config.Hasura) *IntegrationHelpers {
	api := New(config.URL, config.Secret)
	return &IntegrationHelpers{
		api:    api,
		config: config,
	}
}

func TestExpectedMetadataWithActual(
	t *testing.T,
	configPath string,
	expectedMetadataPath string,
) {
	var cfg config.Config
	if err := config.Parse(configPath, &cfg); err != nil {
		log.Err(err).Msg("") // or fail
		return
	}

	integrationHelpers := NewIntegrationHelpers(cfg.Hasura)

	ctx, _ := context.WithCancel(context.Background())

	metadata, err := integrationHelpers.GetMetadata(ctx)
	if err != nil {
		t.Fatalf("Error with getting hasura metadata %e", err)
	}

	expectedMetadata, err := integrationHelpers.ParseExpectedMetadata(expectedMetadataPath)
	if err != nil {
		t.Fatalf("Error with parsing expected metadata: %e", err)
	}

	// Go through `expectedMetadata` and assert that each object
	// in that array is in `metadata` with corresponding columns.
	for _, expectedTable := range expectedMetadata.Tables {
		metadataTableColumns, err := getTableColumns(metadata, expectedTable.Name, "user") // todo: read role from config
		if err != nil {
			t.Fatalf("Error with searching expectedTable in metadata: %s\n%e", expectedTable.Name, err)
		}

		if !elementsMatch(expectedTable, metadataTableColumns) {
			t.Errorf(
				"Table columns do not match: %s\nexpected: %s\nactual: %s",
				expectedTable.Name,
				expectedTable.Columns,
				metadataTableColumns,
			)
		}
	}
}

func elementsMatch(expectedTable ExpectedTable, metadataTable Columns) bool {
	if len(expectedTable.Columns) != len(metadataTable) {
		return false
	}

	hasuraColumns := make(map[string]int)

	for _, columnName := range metadataTable {
		hasuraColumns[columnName] = 0
	}

	for _, expectedColumn := range expectedTable.Columns {
		if _, ok := hasuraColumns[expectedColumn]; !ok {
			return false
		}
	}

	return true
}

func getTableColumns(metadata Metadata, tableName string, role string) (Columns, error) {
	for _, source := range metadata.Sources {
		for _, table := range source.Tables {
			if table.Schema.Name == tableName {
				for _, selectPermission := range table.SelectPermissions {
					if selectPermission.Role == role {
						return selectPermission.Permission.Columns, nil
					}
				}
			}
		}
	}

	return nil, errors.Errorf("Table %s for role %s was not found", tableName, role)
}

func (i *IntegrationHelpers) GetMetadata(ctx context.Context) (Metadata, error) {
	return i.api.ExportMetadata(ctx)
}

func (i *IntegrationHelpers) ParseExpectedMetadata(filename string) (*ExpectedMetadata, error) {
	buf, err := readFile(filename)
	if err != nil {
		return nil, err
	}

	var output ExpectedMetadata
	if err := yaml.NewDecoder(buf).Decode(&output); err != nil {
		return nil, err
	}

	//return validator.New().Struct(&output), nil
	return &output, nil
}

func readFile(filename string) (*bytes.Buffer, error) {
	if filename == "" {
		return nil, errors.Errorf("you have to provide configuration filename")
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "reading file %s", filename)
	}
	//expanded, err := expandVariables(data)
	//if err != nil {
	//	return nil, err
	//}
	return bytes.NewBuffer(data), nil
}
