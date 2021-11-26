package types

import (
	"github.com/dave/jennifer/jen"
	"github.com/dipdup-net/go-lib/tzkt/api"
)

// Timestamp -
type Timestamp struct{}

// AsField -
func (Timestamp) AsField(name, path string, schema api.JSONSchema, isRequired bool, result *ContractTypeResult) (jen.Code, error) {
	tags := map[string]string{
		"json": name,
	}

	if isRequired {
		tags["validate"] = TagRequired
	}

	return jen.Id(fieldName(name)).Add(jen.Id("Timestamp")).Tag(tags), nil
}

// AsCode -
func (Timestamp) AsCode(name, path string, schema api.JSONSchema, result *ContractTypeResult) (Code, error) {
	typName := result.GetName(name)
	return Code{
		Statement: jen.Comment(typName).Line().Type().Id(typName).Add(jen.Id("Timestamp")).Line(),
		Name:      typName,
	}, nil
}

// AsType -
func (Timestamp) AsType(name, path string, schema api.JSONSchema, result *ContractTypeResult) (Code, error) {
	return Code{
		Statement: jen.Add(jen.Id("Timestamp")).Line(),
		Name:      result.GetName(name),
	}, nil
}
