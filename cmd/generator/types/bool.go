package types

import (
	"github.com/dave/jennifer/jen"
	"github.com/dipdup-net/go-lib/tzkt/api"
)

// Bool -
type Bool struct{}

// AsField -
func (Bool) AsField(name, path string, schema api.JSONSchema, isRequired bool, result *ContractTypeResult) (jen.Code, error) {
	tags := map[string]string{
		"json": name,
	}

	if isRequired {
		tags["validate"] = TagRequired
	}

	return jen.Id(fieldName("Bool", name)).Bool().Tag(tags), nil
}

// AsCode -
func (Bool) AsCode(name, path string, schema api.JSONSchema, result *ContractTypeResult) (Code, error) {
	typ := result.GetName("Bool", name)
	return Code{
		Statement: jen.Comment(typ).Line().Type().Id(typ).Bool().Line(),
		Name:      typ,
	}, nil
}

// AsType -
func (Bool) AsType(name, path string, schema api.JSONSchema, result *ContractTypeResult) (Code, error) {
	return Code{
		Statement: jen.Bool(),
		Name:      result.GetName("Bool", name),
	}, nil
}
