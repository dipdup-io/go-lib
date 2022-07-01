package types

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/dipdup-net/go-lib/tzkt/data"
	"github.com/pkg/errors"
)

// Map -
type Map struct{}

// AsField -
func (m Map) AsField(name, path string, schema data.JSONSchema, isRequired bool, result *ContractTypeResult) (jen.Code, error) {
	tags := map[string]string{
		"json": name,
	}

	if isRequired {
		tags["validate"] = TagRequired
	}

	code, err := m.createType(name, path, schema, result)
	if err != nil {
		return code.Statement, err
	}

	return jen.Id(code.Name).Add(code.Statement).Tag(tags), nil
}

// AsCode -
func (m Map) AsCode(name, path string, schema data.JSONSchema, result *ContractTypeResult) (Code, error) {
	code, err := m.createType(name, path, schema, result)
	if err != nil {
		return code, err
	}

	statement := jen.Comment(code.Name).Line().Type().Id(code.Name)
	return Code{
		Statement: statement.Add(code.Statement),
		Name:      code.Name,
	}, nil
}

// AsType -
func (m Map) AsType(name, path string, schema data.JSONSchema, result *ContractTypeResult) (Code, error) {
	return m.createType(name, path, schema, result)
}

func (Map) createType(name, path string, schema data.JSONSchema, result *ContractTypeResult) (Code, error) {
	var code Code
	code.Name = result.GetName("Map", name)

	if schema.PropertyNames == nil {
		return code, errors.Errorf("nil property names in map: %s", code.Name)
	}
	if schema.AdditionalProperties.Value == nil {
		return code, errors.Errorf("nil additional properties in map: %s", code.Name)
	}

	keyType, err := selectType(*schema.PropertyNames)
	if err != nil {
		return code, err
	}

	valType, err := selectType(*schema.AdditionalProperties.Value)
	if err != nil {
		return code, err
	}

	newPath := getPath(path, name)

	if isSimpleType(schema.PropertyNames.Comment) {
		key, err := keyType.AsType(fmt.Sprintf("%sKey", code.Name), newPath, *schema.PropertyNames, result)
		if err != nil {
			return code, err
		}

		value, err := valType.AsType(fmt.Sprintf("%sValue", code.Name), newPath, *schema.AdditionalProperties.Value, result)
		if err != nil {
			return code, err
		}

		code.Statement = jen.Map(key.Statement).Add(value.Statement)
	} else {
		key, err := keyType.AsCode(fmt.Sprintf("%sKey", code.Name), newPath, *schema.PropertyNames, result)
		if err != nil {
			return code, err
		}

		value, err := valType.AsCode(fmt.Sprintf("%sValue", code.Name), newPath, *schema.AdditionalProperties.Value, result)
		if err != nil {
			return code, err
		}

		code.Statement = jen.Map(jen.Id(key.Name)).Add(jen.Id(value.Name)).Line().Add(key.Statement, value.Statement).Line()
	}

	return code, nil
}
