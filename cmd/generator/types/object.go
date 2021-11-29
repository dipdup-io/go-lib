package types

import (
	"github.com/dave/jennifer/jen"
	"github.com/dipdup-net/go-lib/tzkt/api"
)

// Object -
type Object struct{}

// AsField -
func (Object) AsField(name, path string, schema api.JSONSchema, isRequired bool, result *ContractTypeResult) (jen.Code, error) {
	tags := map[string]string{
		"json": name,
	}

	if isRequired {
		tags["validate"] = TagRequired
	}

	return jen.Id(fieldName("Object", name)).Add(jen.Id("Unit")).Tag(tags), nil
}

// AsCode -
func (obj Object) AsCode(name, path string, schema api.JSONSchema, result *ContractTypeResult) (Code, error) {
	code, err := obj.create(name, path, schema, result)
	if err != nil {
		return Code{}, err
	}

	statement := jen.Comment(code.Name).Line().Type().Id(code.Name)
	return Code{
		Statement: statement.Add(code.Statement),
		Name:      code.Name,
	}, nil
}

// AsType -
func (obj Object) AsType(name, path string, schema api.JSONSchema, result *ContractTypeResult) (Code, error) {
	return obj.create(name, path, schema, result)
}

func (Object) create(name, path string, schema api.JSONSchema, result *ContractTypeResult) (Code, error) {
	fields := make([]jen.Code, 0)
	statements := make([]jen.Code, 0)

	for propName, props := range schema.Properties {
		var required bool
		for j := range schema.Required {
			if schema.Required[j] == propName {
				required = true
				break
			}
		}

		typ, err := selectType(props)
		if err != nil {
			return Code{}, err
		}

		newPath := getPath(path, name)

		if isSimpleType(props.Comment) || props.Comment == "big_map" {
			field, err := typ.AsField(propName, newPath, props, required, result)
			if err != nil {
				return Code{}, err
			}
			fields = append(fields, field)
		} else {
			code, err := typ.AsCode(propName, newPath, props, result)
			if err != nil {
				return Code{}, err
			}

			if code.Statement != nil {
				statements = append(statements, code.Statement)
				tags := map[string]string{
					"json": propName,
				}

				if required {
					tags["validate"] = TagRequired
				}

				fields = append(fields, jen.Id(fieldName(props.Comment, propName)).Op(code.Name).Tag(tags))
			}
		}
	}

	return Code{
		Statement: jen.Struct(fields...).Line().Add(statements...).Line(),
		Name:      result.GetName("Object", name),
	}, nil
}
