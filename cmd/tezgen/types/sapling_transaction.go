package types

import (
	"github.com/dave/jennifer/jen"
	"github.com/dipdup-net/go-lib/tzkt/data"
)

// SaplingTransaction -
type SaplingTransaction struct{}

// AsField -
func (SaplingTransaction) AsField(name, path string, schema data.JSONSchema, isRequired bool, result *ContractTypeResult) (jen.Code, error) {
	tags := map[string]string{
		"json": name,
	}

	if isRequired {
		tags["validate"] = TagRequired
	}

	return jen.Id(fieldName("SaplingTransaction", name)).Add(
		jen.Qual("github.com/dipdup-net/go-lib/tools/tezgen", "SaplingTransaction"),
	).Tag(tags), nil
}

// AsCode -
func (SaplingTransaction) AsCode(name, path string, schema data.JSONSchema, result *ContractTypeResult) (Code, error) {
	typName := result.GetName("SaplingTransaction", name)
	return Code{
		Statement: jen.Comment(typName).Line().Type().Id(typName).Add(
			jen.Qual("github.com/dipdup-net/go-lib/tools/tezgen", "SaplingTransaction"),
		).Line(),
		Name: typName,
	}, nil
}

// AsType -
func (SaplingTransaction) AsType(name, path string, schema data.JSONSchema, result *ContractTypeResult) (Code, error) {
	return Code{
		Statement: jen.Add(
			jen.Qual("github.com/dipdup-net/go-lib/tools/tezgen", "SaplingTransaction"),
		),
		Name: name,
	}, nil
}
