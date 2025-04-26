package hasura

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/dipdup-net/go-lib/config"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type ExpectedMetadata struct {
	Tables []ExpectedTable `yaml:"tables" validate:"required"`
}

type ExpectedTable struct {
	Name    string   `yaml:"name" validate:"required"`
	Columns []string `yaml:"columns" validate:"required"`
}

func TestExpectedMetadataWithActual(
	t *testing.T,
	configPath string,
	expectedMetadataPath string,
) {
	var cfg config.Config
	if err := config.Parse(configPath, &cfg); err != nil {
		t.Fatalf("Error with reading configuratoin file: %s", err)
	}

	api := New(cfg.Hasura.URL, cfg.Hasura.Secret)
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	metadata, err := api.ExportMetadata(ctx)
	if err != nil {
		t.Fatalf("Error with getting hasura metadata %e", err)
	}

	expectedMetadata, err := parseExpectedMetadata(expectedMetadataPath)
	if err != nil {
		t.Fatalf("Error with parsing expected metadata: %e", err)
	}

	// Go through `expectedMetadata` and assert that each object
	// in that array is in `metadata` with corresponding columns.
	for _, expectedTable := range expectedMetadata.Tables {
		metadataTableColumns, err := getTableColumns(metadata, expectedTable.Name, cfg.Hasura.Source.Name, "user")
		if err != nil {
			t.Fatalf("Error with searching expectedTable in metadata: %s\n%s", expectedTable.Name, err)
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

	hasuraColumns := make(map[string]struct{})

	for _, columnName := range metadataTable {
		hasuraColumns[columnName] = struct{}{}
	}

	for _, expectedColumn := range expectedTable.Columns {
		if _, ok := hasuraColumns[expectedColumn]; !ok {
			return false
		}
	}

	return true
}

func getTableColumns(metadata Metadata, tableName string, sourceName string, role string) (Columns, error) {
	for _, source := range metadata.Sources {
		if source.Name != sourceName {
			continue
		}

		for _, table := range source.Tables {
			if table.Schema.Name != tableName {
				continue
			}

			for _, selectPermission := range table.SelectPermissions {
				if selectPermission.Role == role {
					return selectPermission.Permission.Columns, nil
				}
			}
		}
	}

	return nil, errors.Errorf("Table %s for role %s was not found", tableName, role)
}

func parseExpectedMetadata(filename string) (*ExpectedMetadata, error) {
	buf, err := readFile(filename)
	if err != nil {
		return nil, err
	}

	var output ExpectedMetadata
	if err := yaml.NewDecoder(buf).Decode(&output); err != nil {
		return nil, err
	}

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

	return bytes.NewBuffer(data), nil
}
