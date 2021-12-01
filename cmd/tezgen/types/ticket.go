package types

import (
	"github.com/dave/jennifer/jen"
	"github.com/dipdup-net/go-lib/tzkt/api"
)

// Ticket -
type Ticket struct{}

// AsField -
func (Ticket) AsField(name, path string, schema api.JSONSchema, isRequired bool, result *ContractTypeResult) (jen.Code, error) {
	tags := map[string]string{
		"json": name,
	}

	if isRequired {
		tags["validate"] = TagRequired
	}

	return jen.Id(fieldName("Ticket", name)).String().Tag(tags), nil
}

// AsCode -
func (Ticket) AsCode(name, path string, schema api.JSONSchema, result *ContractTypeResult) (Code, error) {
	typName := result.GetName("Ticket", name)
	return Code{
		Statement: jen.Comment(typName).Line().Type().Id(typName).String().Line(),
		Name:      typName,
	}, nil
}

// AsType -
func (Ticket) AsType(name, path string, schema api.JSONSchema, result *ContractTypeResult) (Code, error) {
	return Code{
		Statement: jen.String(),
		Name:      name,
	}, nil
}
