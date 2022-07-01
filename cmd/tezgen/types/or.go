package types

import (
	"errors"
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/dipdup-net/go-lib/tzkt/data"
	"github.com/iancoleman/strcase"
)

// Object -
type Or struct{}

// AsField -
func (Or) AsField(name, path string, schema data.JSONSchema, isRequired bool, result *ContractTypeResult) (jen.Code, error) {
	tags := map[string]string{
		"json": name,
	}

	if isRequired {
		tags["validate"] = TagRequired
	}

	return jen.Id(fieldName("Or", name)).Add(jen.Id(result.GetName("Or", name))).Tag(tags), nil
}

// AsCode -
func (obj Or) AsCode(name, path string, schema data.JSONSchema, result *ContractTypeResult) (Code, error) {
	code, err := obj.createOr(name, path, schema, result)
	if err != nil {
		return code, err
	}

	return Code{
		Statement: jen.Comment(code.Name).Line().Type().Id(code.Name).Add(code.Statement),
		Name:      code.Name,
	}, nil
}

// AsType -
func (obj Or) AsType(name, path string, schema data.JSONSchema, result *ContractTypeResult) (Code, error) {
	return obj.createOr(name, path, schema, result)
}

func (obj Or) createOr(name, path string, schema data.JSONSchema, result *ContractTypeResult) (Code, error) {
	if len(schema.OneOf) != 2 {
		return Code{}, errors.New("invalid oneOf key for or type")
	}

	leftField, leftStatement, err := obj.createSide(name, "left", path, schema.OneOf[0], result)
	if err != nil {
		return Code{}, err
	}
	rightField, rightStatement, err := obj.createSide(name, "right", path, schema.OneOf[1], result)
	if err != nil {
		return Code{}, err
	}

	typName := result.GetName("Or", name)

	unmarshalJSON := obj.getUnmarshalJSON(typName)
	marshalJSON := obj.getMarshalJSON(typName)

	statement := jen.Struct(
		leftField,
		rightField,
	).Line().Line().Add(unmarshalJSON).Line().Add(marshalJSON).Line()

	if leftStatement != nil {
		statement.Add(leftStatement)
	}
	if rightStatement != nil {
		statement.Add(rightStatement)
	}
	return Code{
		Statement: statement,
		Name:      typName,
	}, nil
}

func (obj Or) createSide(name, side, path string, schema data.JSONSchema, result *ContractTypeResult) (jen.Code, jen.Code, error) {
	typ, err := selectType(schema)
	if err != nil {
		return nil, nil, err
	}

	if isSimpleType(schema.Comment) {
		code, err := typ.AsField(side, path, schema, false, result)
		return code, nil, err
	}

	sideName := fmt.Sprintf("%s_%s", side, name)
	newPath := getPath(path, name)
	code, err := typ.AsCode(sideName, newPath, schema, result)
	if err != nil {
		return nil, nil, err
	}
	field := jen.Id(strcase.ToCamel(side)).Op("*").Id(code.Name)
	return field, code.Statement, nil
}

func (obj Or) getUnmarshalJSON(typName string) jen.Code {
	return jen.Comment("UnmarshalJSON").Line().
		Func().Params(
		jen.Id("or").Op("*").Id(typName),
	).Id("UnmarshalJSON").Params(
		jen.Id("data").Op("[]").Byte(),
	).Error().Block(
		jen.If(
			jen.Err().Op(":=").Id("json.Unmarshal").
				Call(
					jen.Id("data"),
					jen.Id("or.Left"),
				),
			jen.Err().Op("!=").Nil(),
		).Block(
			jen.Return(jen.Nil()),
		),

		jen.Return(
			jen.Id("json.Unmarshal").
				Call(
					jen.Id("data"),
					jen.Id("or.Right"),
				),
		),
	)
}

func (obj Or) getMarshalJSON(typName string) jen.Code {
	return jen.Comment("MarshalJSON").Line().
		Func().Params(
		jen.Id("or").Op("*").Id(typName),
	).Id("MarshalJSON").Params().Params(
		jen.Op("[]").Byte(),
		jen.Error(),
	).Block(
		jen.If(
			jen.Id("or.Left").Op("!=").Nil().Block(
				jen.Return(
					jen.Id("json.Marshal").
						Call(
							jen.Id("or.Left"),
						),
				),
			),
		),
		jen.If(
			jen.Id("or.Right").Op("!=").Nil().Block(
				jen.Return(
					jen.Id("json.Marshal").
						Call(
							jen.Id("or.Right"),
						),
				),
			),
		),
		jen.Return(
			jen.Nil(),
			jen.Qual("errors", "New").Call(
				jen.Lit("left and right in or type are not initialized"),
			),
		),
	)
}
