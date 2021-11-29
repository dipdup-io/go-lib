package types

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/dipdup-net/go-lib/tzkt/api"
)

// Int -
type Int struct{}

// AsField -
func (Int) AsField(name, path string, schema api.JSONSchema, isRequired bool, result *ContractTypeResult) (jen.Code, error) {
	tags := map[string]string{
		"json": fmt.Sprintf("%s,string", name),
	}

	if isRequired {
		tags["validate"] = TagRequired
	}

	return jen.Id(fieldName("Int", name)).Int64().Tag(tags), nil
}

// AsCode -
func (Int) AsCode(name, path string, schema api.JSONSchema, result *ContractTypeResult) (Code, error) {
	typ := result.GetName("Int", name)
	return Code{
		Statement: jen.Comment(typ).Line().Type().Id(typ).Int64().Line(),
		Name:      typ,
	}, nil
}

// AsType -
func (Int) AsType(name, path string, schema api.JSONSchema, result *ContractTypeResult) (Code, error) {
	return Code{
		Statement: jen.Int64(),
		Name:      result.GetName("Int", name),
	}, nil
}
