package types

import (
	"github.com/dave/jennifer/jen"
	"github.com/dipdup-net/go-lib/tzkt/data"
)

// String -
type String struct{}

// AsField -
func (String) AsField(name, path string, schema data.JSONSchema, isRequired bool, result *ContractTypeResult) (jen.Code, error) {
	tags := map[string]string{
		"json": name,
	}

	if isRequired {
		tags["validate"] = TagRequired
	}

	return jen.Id(fieldName(schema.Comment, name)).String().Tag(tags), nil
}

// AsCode -
func (String) AsCode(name, path string, schema data.JSONSchema, result *ContractTypeResult) (Code, error) {
	typ := result.GetName(schema.Comment, name)
	return Code{
		Statement: jen.Comment(typ).Line().Type().Id(typ).String().Line(),
		Name:      typ,
	}, nil
}

// AsType -
func (String) AsType(name, path string, schema data.JSONSchema, result *ContractTypeResult) (Code, error) {
	return Code{
		Statement: jen.String(),
		Name:      result.GetName(schema.Comment, name),
	}, nil
}
