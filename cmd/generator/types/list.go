package types

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/dipdup-net/go-lib/tzkt/api"
)

// List -
type List struct{}

// AsField -
func (List) AsField(name, path string, schema api.JSONSchema, isRequired bool, result *ContractTypeResult) (jen.Code, error) {
	tags := map[string]string{
		"json": name,
	}

	if isRequired {
		tags["validate"] = TagRequired
	}

	return jen.Id(fieldName("List", name)).Op("[]").Add(jen.Id(result.GetName("List", name))).Tag(tags), nil
}

// AsCode -
func (l List) AsCode(name, path string, schema api.JSONSchema, result *ContractTypeResult) (Code, error) {
	code, err := l.createList(name, path, schema, result)
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
func (l List) AsType(name, path string, schema api.JSONSchema, result *ContractTypeResult) (Code, error) {
	return l.createList(name, path, schema, result)
}

func (List) createList(name, path string, schema api.JSONSchema, result *ContractTypeResult) (Code, error) {
	var code Code
	code.Name = result.GetName("List", name)

	if schema.Items != nil {
		typ, err := selectType(*schema.Items)
		if err != nil {
			return code, err
		}

		itemType, err := typ.AsType(fmt.Sprintf("%s_item", name), getPath(path, name), *schema.Items, result)
		if err != nil {
			return code, err
		}

		code.Statement = jen.Op("[]").Add(itemType.Statement)
	}

	return code, nil
}
