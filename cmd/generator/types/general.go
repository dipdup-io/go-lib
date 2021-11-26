package types

import (
	"fmt"
	"go/token"

	"github.com/dave/jennifer/jen"
	"github.com/dipdup-net/go-lib/tzkt/api"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
)

// TODO: custom validation by type
// TODO: re-using same types

var reservedNames = map[string]string{
	"default": "DefaultEntrypoint",
}

// Type -
type Type interface {
	AsField(name, path string, schema api.JSONSchema, isRequired bool, result *ContractTypeResult) (jen.Code, error)
	AsCode(name, path string, schema api.JSONSchema, result *ContractTypeResult) (Code, error)
	AsType(name, path string, schema api.JSONSchema, result *ContractTypeResult) (Code, error)
}

// Code -
type Code struct {
	Statement jen.Code
	Name      string
}

// Generate -
func Generate(name string, schema api.JSONSchema, result *ContractTypeResult) (string, error) {
	typ, err := selectType(schema)
	if err != nil {
		return "", err
	}

	code, err := typ.AsCode(name, "", schema, result)
	if err != nil {
		return "", err
	}
	if code.Statement != nil {
		result.File.Add(code.Statement)
	}
	return code.Name, nil
}

// ContractTypeResult -
type ContractTypeResult struct {
	File        *jen.File
	Entrypoints map[string]EntrypointData
	PackageName string

	names   map[string]struct{}
	bigMaps map[string]string
	counter int64
}

// GetName -
func (result *ContractTypeResult) GetName(name string) string {
	if reserved, ok := reservedNames[name]; ok {
		name = reserved
	}
	name = strcase.ToCamel(name)

	if _, exists := result.names[name]; exists {
		result.counter++
		name = fmt.Sprintf("%s%d", name, result.counter)
		name = result.GetName(name)
	}

	result.names[name] = struct{}{}
	return name
}

// EntrypointData -
type EntrypointData struct {
	Type string
	Var  string
}

// GenerateContractTypes -
func GenerateContractTypes(schema api.ContractJSONSchema, packageName string) (ContractTypeResult, error) {
	result := ContractTypeResult{
		File:        jen.NewFile(packageName),
		Entrypoints: make(map[string]EntrypointData),
		PackageName: packageName,

		names:   make(map[string]struct{}),
		bigMaps: make(map[string]string),
	}

	result.File.PackageComment("DO NOT EDIT!!!")
	result.File.ImportName("encoding/json", "json")

	if err := generateForContract(schema, &result); err != nil {
		return result, err
	}

	return result, nil
}

func generateForContract(schema api.ContractJSONSchema, result *ContractTypeResult) error {
	for _, entrypoint := range schema.Entrypoints {
		entrypointType, err := Generate(entrypoint.Name, entrypoint.Parameter, result)
		if err != nil {
			return err
		}
		entrypointName := entrypoint.Name
		if token.Lookup(entrypointName).IsKeyword() {
			entrypointName = fmt.Sprintf("%s_entrypoint", entrypointName)
		}

		result.Entrypoints[entrypoint.Name] = EntrypointData{
			Type: entrypointType,
			Var:  strcase.ToLowerCamel(entrypointName),
		}
	}

	for i := range schema.BigMaps {
		if err := GenerateBigMap(schema.BigMaps[i], result); err != nil {
			return err
		}
	}

	if _, err := Generate("storage", schema.Storage, result); err != nil {
		return err
	}

	return nil
}

func selectType(schema api.JSONSchema) (Type, error) {
	switch schema.Type {
	case "object":
		switch schema.Comment {
		case "unit":
			return new(Unit), nil
		case "map":
			return new(Map), nil
		default:
			return new(Object), nil
		}
	case TypeString:
		switch schema.Comment {
		case "address":
			return new(Address), nil
		case "bytes", "bls12_381_g1", "bls12_381_g2":
			return new(Bytes), nil
		case "contract":
			return new(Contract), nil
		case "mutez":
			return new(Mutez), nil
		case "int":
			return new(Int), nil
		case "key_hash":
			return new(KeyHash), nil
		case "lambda":
			return new(Lambda), nil
		case "nat":
			return new(Nat), nil
		case "sapling_transaction":
			return new(SaplingTransaction), nil
		case TypeString, "signature", "key", "never", "chain_id", "bls12_381_fr", "ticket", "operation":
			return new(String), nil
		case "timestamp":
			return new(Timestamp), nil
		default:
			return nil, errors.Errorf("unknown comment for string: %s", schema.Comment)
		}
	case "array":
		switch schema.Comment {
		case "list":
			return new(List), nil
		case "set":
			return new(Set), nil
		default:
			return nil, errors.Errorf("unknown comment for array: %s", schema.Comment)
		}
	case "integer":
		return new(Int), nil
	case "boolean":
		return new(Bool), nil
	default:
		switch schema.Comment {
		case "option":
			return new(Option), nil
		case "big_map":
			return new(BigMap), nil
		case "sapling_state":
			return new(SaplingState), nil
		case "or":
			return new(Or), nil
		}
		return nil, errors.Errorf("unknown type: %s", schema.Type)
	}
}

func isSimpleType(comment string) bool {
	switch comment {
	case TypeString, "bytes", "address", "mutez", "int", "nat", "timestamp", "unit", "signature":
		return true
	}
	return false
}

func fieldName(name string) string {
	return strcase.ToCamel(name)
}

func getPath(path string, name string) string {
	if path == "" {
		return name
	}
	return fmt.Sprintf("%s.%s", path, name)
}
