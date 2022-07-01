package types

import (
	"github.com/dave/jennifer/jen"
	"github.com/dipdup-net/go-lib/tzkt/data"
)

// Contract -
type Contract struct{}

// AsField -
func (Contract) AsField(name, path string, schema data.JSONSchema, isRequired bool, result *ContractTypeResult) (jen.Code, error) {
	tags := map[string]string{
		"json": name,
	}

	if isRequired {
		tags["validate"] = TagRequired
	}

	return jen.Id(fieldName("Contract", name)).Add(
		jen.Qual("github.com/dipdup-net/go-lib/tools/tezgen", "Contract"),
	).Tag(tags), nil
}

// AsCode -
func (Contract) AsCode(name, path string, schema data.JSONSchema, result *ContractTypeResult) (Code, error) {
	typName := result.GetName("Contract", name)
	return Code{
		Statement: jen.Comment(typName).Line().Type().Id(typName).Add(
			jen.Qual("github.com/dipdup-net/go-lib/tools/tezgen", "Contract"),
		).Line(),
		Name: typName,
	}, nil
}

// AsType -
func (Contract) AsType(name, path string, schema data.JSONSchema, result *ContractTypeResult) (Code, error) {
	return Code{
		Statement: jen.Add(
			jen.Qual("github.com/dipdup-net/go-lib/tools/tezgen", "Contract"),
		),
		Name: name,
	}, nil
}
