package types

import (
	"github.com/dave/jennifer/jen"
	"github.com/dipdup-net/go-lib/tzkt/api"
)

// Contract -
type Contract struct{}

// AsField -
func (Contract) AsField(name, path string, schema api.JSONSchema, isRequired bool, result *ContractTypeResult) (jen.Code, error) {
	tags := map[string]string{
		"json": name,
	}

	if isRequired {
		tags["validate"] = TagRequired
	}

	return jen.Id(fieldName(name)).Add(jen.Id("Contract")).Tag(tags), nil
}

// AsCode -
func (Contract) AsCode(name, path string, schema api.JSONSchema, result *ContractTypeResult) (Code, error) {
	typName := result.GetName(name)
	return Code{
		Statement: jen.Comment(typName).Line().Type().Id(typName).Add(jen.Id("Contract")).Line(),
		Name:      typName,
	}, nil
}

// AsType -
func (Contract) AsType(name, path string, schema api.JSONSchema, result *ContractTypeResult) (Code, error) {
	return Code{
		Statement: jen.Add(jen.Id("Contract")),
		Name:      name,
	}, nil
}
