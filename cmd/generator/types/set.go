package types

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/dipdup-net/go-lib/tzkt/api"
	"github.com/pkg/errors"
)

// Set -
type Set struct{}

// AsField -
func (s Set) AsField(name, path string, schema api.JSONSchema, isRequired bool, result *ContractTypeResult) (jen.Code, error) {
	tags := map[string]string{
		"json": name,
	}

	if isRequired {
		tags["validate"] = TagRequired
	}

	code, err := s.createType(name, path, schema, result)
	if err != nil {
		return code.Statement, err
	}

	return jen.Id(fieldName("Set", name)).Add(code.Statement).Tag(tags), nil
}

// AsCode -
func (s Set) AsCode(name, path string, schema api.JSONSchema, result *ContractTypeResult) (Code, error) {
	code, err := s.createType(name, path, schema, result)
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
func (s Set) AsType(name, path string, schema api.JSONSchema, result *ContractTypeResult) (Code, error) {
	return s.createType(name, path, schema, result)
}

func (Set) createType(name, path string, schema api.JSONSchema, result *ContractTypeResult) (Code, error) {
	var code Code
	code.Name = result.GetName("Set", name)

	if schema.Items == nil {
		return code, errors.Errorf("nil items in set: %s", code.Name)
	}

	itemsType, err := selectType(*schema.Items)
	if err != nil {
		return code, err
	}

	newPath := getPath(path, name)

	if isSimpleType(schema.Items.Comment) {
		item, err := itemsType.AsType(fmt.Sprintf("%s_item", name), newPath, *schema.Items, result)
		if err != nil {
			return code, err
		}

		code.Statement = jen.Map(item.Statement).Struct().Line()
	} else {
		item, err := itemsType.AsCode(fmt.Sprintf("%s_item", name), newPath, *schema.Items, result)
		if err != nil {
			return code, err
		}

		code.Statement = jen.Map(jen.Id(item.Name)).Struct().Line().Add(item.Statement).Line()
	}

	return code, nil
}
