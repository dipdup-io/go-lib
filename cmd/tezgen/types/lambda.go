package types

import (
	"github.com/dave/jennifer/jen"
	"github.com/dipdup-net/go-lib/tzkt/api"
)

// Lambda -
type Lambda struct{}

// AsField -
func (Lambda) AsField(name, path string, schema api.JSONSchema, isRequired bool, result *ContractTypeResult) (jen.Code, error) {
	tags := map[string]string{
		"json": name,
	}

	if isRequired {
		tags["validate"] = TagRequired
	}

	return jen.Id(fieldName("Lambda", name)).String().Tag(tags), nil
}

// AsCode -
func (Lambda) AsCode(name, path string, schema api.JSONSchema, result *ContractTypeResult) (Code, error) {
	typName := result.GetName("Lambda", name)
	return Code{
		Statement: jen.Comment(typName).Line().Type().Id(typName).String().Line(),
		Name:      typName,
	}, nil
}

// AsType -
func (Lambda) AsType(name, path string, schema api.JSONSchema, result *ContractTypeResult) (Code, error) {
	return Code{
		Statement: jen.String(),
		Name:      name,
	}, nil
}
