package types

import (
	"github.com/dave/jennifer/jen"
	"github.com/dipdup-net/go-lib/tzkt/api"
)

// Address -
type Address struct{}

// AsField -
func (Address) AsField(name, path string, schema api.JSONSchema, isRequired bool, result *ContractTypeResult) (jen.Code, error) {
	tags := map[string]string{
		"json": name,
	}

	if isRequired {
		tags["validate"] = TagRequired
	}

	return jen.Id(fieldName("Address", name)).Add(
		jen.Qual("github.com/dipdup-net/go-lib/tools/tezgen", "Address"),
	).Tag(tags), nil
}

// AsCode -
func (Address) AsCode(name, path string, schema api.JSONSchema, result *ContractTypeResult) (Code, error) {
	typName := result.GetName("Address", name)
	return Code{
		Statement: jen.Comment(typName).Line().Type().Id(typName).Add(
			jen.Qual("github.com/dipdup-net/go-lib/tools/tezgen", "Address"),
		).Line(),
		Name: typName,
	}, nil
}

// AsType -
func (Address) AsType(name, path string, schema api.JSONSchema, result *ContractTypeResult) (Code, error) {
	return Code{
		Statement: jen.Add(
			jen.Qual("github.com/dipdup-net/go-lib/tools/tezgen", "Address"),
		),
		Name: name,
	}, nil
}
