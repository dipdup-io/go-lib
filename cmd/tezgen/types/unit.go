package types

import (
	"github.com/dave/jennifer/jen"
	"github.com/dipdup-net/go-lib/tzkt/api"
)

// Unit -
type Unit struct{}

// AsField -
func (Unit) AsField(name, path string, schema api.JSONSchema, isRequired bool, result *ContractTypeResult) (jen.Code, error) {
	tags := map[string]string{
		"json": name,
	}

	if isRequired {
		tags["validate"] = TagRequired
	}

	return jen.Id(fieldName("Unit", name)).Add(
		jen.Qual("github.com/dipdup-net/go-lib/tools/tezgen", "Unit"),
	).Tag(tags), nil
}

// AsCode -
func (Unit) AsCode(name, path string, schema api.JSONSchema, result *ContractTypeResult) (Code, error) {
	typName := result.GetName("Unit", name)
	return Code{
		Statement: jen.Comment(typName).Line().Type().Id(typName).Add(
			jen.Qual("github.com/dipdup-net/go-lib/tools/tezgen", "Unit"),
		),
		Name: typName,
	}, nil
}

// AsType -
func (Unit) AsType(name, path string, schema api.JSONSchema, result *ContractTypeResult) (Code, error) {
	return Code{
		Statement: jen.Add(
			jen.Qual("github.com/dipdup-net/go-lib/tools/tezgen", "Unit"),
		).Line(),
		Name: result.GetName("Unit", name),
	}, nil
}
