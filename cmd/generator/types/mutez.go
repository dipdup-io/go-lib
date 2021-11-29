package types

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/dipdup-net/go-lib/tzkt/api"
)

// Mutez -
type Mutez struct{}

// AsField -
func (Mutez) AsField(name, path string, schema api.JSONSchema, isRequired bool, result *ContractTypeResult) (jen.Code, error) {
	tags := map[string]string{
		"json": fmt.Sprintf("%s,string", name),
	}

	if isRequired {
		tags["validate"] = TagRequired
	}

	return jen.Id(fieldName("Mutez", name)).Int64().Tag(tags), nil
}

// AsCode -
func (Mutez) AsCode(name, path string, schema api.JSONSchema, result *ContractTypeResult) (Code, error) {
	typ := result.GetName("Mutez", name)
	return Code{
		Statement: jen.Comment(typ).Line().Type().Id(typ).Int64().Line(),
		Name:      typ,
	}, nil
}

// AsType -
func (Mutez) AsType(name, path string, schema api.JSONSchema, result *ContractTypeResult) (Code, error) {
	return Code{
		Statement: jen.Int64().Line(),
		Name:      result.GetName("Mutez", name),
	}, nil
}
