package types

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/dipdup-net/go-lib/tzkt/api"
	"github.com/pkg/errors"
)

// Option -
type Option struct{}

// AsField -
func (Option) AsField(name, path string, schema api.JSONSchema, isRequired bool, result *ContractTypeResult) (jen.Code, error) {
	tags := map[string]string{
		"json": fmt.Sprintf("%s,omitempty", name),
	}

	if isRequired {
		tags["validate"] = TagRequired
	}

	return jen.Id(fieldName(name)).Op("[]").Add(jen.Id(result.GetName(name))).Tag(tags), nil
}

// AsCode -
func (opt Option) AsCode(name, path string, schema api.JSONSchema, result *ContractTypeResult) (Code, error) {
	optType, err := opt.createOption(name, path, schema, result)
	if err != nil {
		return Code{}, err
	}

	statement := jen.Comment(optType.Name).Line().Type().Id(optType.Name)
	return Code{
		Statement: statement.Add(optType.Statement),
		Name:      optType.Name,
	}, nil
}

// AsType -
func (opt Option) AsType(name, path string, schema api.JSONSchema, result *ContractTypeResult) (Code, error) {
	return opt.createOption(name, path, schema, result)
}

func (Option) createOption(name, path string, schema api.JSONSchema, result *ContractTypeResult) (Code, error) {
	var code Code
	code.Name = result.GetName(name)

	if len(schema.OneOf) != 2 {
		return code, errors.Errorf("invalid oneOf field for option: %s", code.Name)
	}

	var entityType api.JSONSchema
	for i := range schema.OneOf {
		switch schema.OneOf[i].Type {
		case "null":
		default:
			entityType = schema.OneOf[i]
		}
	}

	newPath := getPath(path, name)

	option, err := selectType(entityType)
	if err != nil {
		return code, err
	}
	optType, err := option.AsType(fmt.Sprintf("%s_option", name), newPath, entityType, result)
	if err != nil {
		return code, err
	}

	statement := jen.Op("*").Add(optType.Statement)

	if !isSimpleType(entityType.Comment) {
		optCode, err := option.AsCode(fmt.Sprintf("%s_option", name), newPath, entityType, result)
		if err != nil {
			return code, err
		}
		statement = statement.Add(optCode.Statement)
	}

	return Code{
		Statement: statement.Line(),
		Name:      code.Name,
	}, nil
}
