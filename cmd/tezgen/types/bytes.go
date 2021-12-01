package types

import (
	"github.com/dave/jennifer/jen"
	"github.com/dipdup-net/go-lib/tzkt/api"
)

// Bytes -
type Bytes struct{}

// AsField -
func (Bytes) AsField(name, path string, schema api.JSONSchema, isRequired bool, result *ContractTypeResult) (jen.Code, error) {
	tags := map[string]string{
		"json": name,
	}

	if isRequired {
		tags["validate"] = TagRequired
	}

	return jen.Id(fieldName("Bytes", name)).Add(
		jen.Qual("github.com/dipdup-net/go-lib/tools/tezgen", "Bytes"),
	).Tag(tags), nil
}

// AsCode -
func (Bytes) AsCode(name, path string, schema api.JSONSchema, result *ContractTypeResult) (Code, error) {
	typName := result.GetName("Bytes", name)
	return Code{
		Statement: jen.Comment(typName).Line().Type().Id(typName).Add(
			jen.Qual("github.com/dipdup-net/go-lib/tools/tezgen", "Bytes"),
		).Line(),
		Name: typName,
	}, nil
}

// AsType -
func (Bytes) AsType(name, path string, schema api.JSONSchema, result *ContractTypeResult) (Code, error) {
	return Code{
		Statement: jen.Add(
			jen.Qual("github.com/dipdup-net/go-lib/tools/tezgen", "Bytes"),
		),
		Name: result.GetName("Bytes", name),
	}, nil
}
