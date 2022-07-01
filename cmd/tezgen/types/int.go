package types

import (
	"github.com/dave/jennifer/jen"
	"github.com/dipdup-net/go-lib/tzkt/data"
)

// Int -
type Int struct{}

// AsField -
func (Int) AsField(name, path string, schema data.JSONSchema, isRequired bool, result *ContractTypeResult) (jen.Code, error) {
	tags := map[string]string{
		"json": name,
	}

	if isRequired {
		tags["validate"] = TagRequired
	}

	return jen.Id(fieldName(schema.Comment, name)).Add(
		jen.Qual("github.com/dipdup-net/go-lib/tools/tezgen", "Int"),
	).Tag(tags), nil
}

// AsCode -
func (Int) AsCode(name, path string, schema data.JSONSchema, result *ContractTypeResult) (Code, error) {
	typ := result.GetName(schema.Comment, name)
	return Code{
		Statement: jen.Comment(typ).Line().Type().Id(typ).Add(
			jen.Qual("github.com/dipdup-net/go-lib/tools/tezgen", "Int"),
		).Line(),
		Name: typ,
	}, nil
}

// AsType -
func (Int) AsType(name, path string, schema data.JSONSchema, result *ContractTypeResult) (Code, error) {
	return Code{
		Statement: jen.Add(
			jen.Qual("github.com/dipdup-net/go-lib/tools/tezgen", "Int"),
		),
		Name: result.GetName(schema.Comment, name),
	}, nil
}
