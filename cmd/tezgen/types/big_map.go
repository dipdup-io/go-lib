package types

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/dipdup-net/go-lib/tzkt/api"
)

// GenerateBigMap -
func GenerateBigMap(bigMap api.BigMapJSONSchema, result *ContractTypeResult) error {
	keyTypeName, err := Generate(fmt.Sprintf("key_%s", bigMap.Name), bigMap.Key, result)
	if err != nil {
		return err
	}

	valueTypeName, err := Generate(fmt.Sprintf("value_%s", bigMap.Name), bigMap.Value, result)
	if err != nil {
		return err
	}

	typeName := result.GetName("BigMap", bigMap.Name)
	result.File.Comment(typeName).Line().Type().Id(typeName).Struct(
		jen.Id("Key").Add(jen.Id(keyTypeName)),
		jen.Id("Value").Add(jen.Id(valueTypeName)),
		jen.Id("Ptr").Add(jen.Op("*").Int64()),
	).Line().
		Comment("UnmarshalJSON").Line().
		Func().Params(
		jen.Id("b").Op("*").Id(typeName),
	).Id("UnmarshalJSON").Params(
		jen.Id("data").Index().Byte(),
	).Error().Block(
		jen.List(jen.Id("ptr"), jen.Err()).Op(":=").Qual("strconv", "ParseInt").Call(
			jen.String().Call(jen.Id("data")),
			jen.Lit(10),
			jen.Lit(64),
		).Line().
			If(
				jen.Err().Op("==").Nil(),
			).Block(
			jen.Id("b.Ptr").Op("=").Op("&").Id("ptr"),
			jen.Return(jen.Nil()),
		).Line().
			Id("parts").Op(":=").Index().Interface().Values(
			jen.Id("b.Key"), jen.Id("b.Value"),
		).Line().Return(
			jen.Id("json.Unmarshal").Call(
				jen.Id("data"),
				jen.Op("&").Id("parts"),
			),
		),
	)

	result.BigMaps[getPath("storage", bigMap.Path)] = BigMapData{
		KeyType:   keyTypeName,
		ValueType: valueTypeName,
		Type:      typeName,
	}

	return nil
}

// BigMap -
type BigMap struct{}

// AsField -
func (BigMap) AsField(name, path string, schema api.JSONSchema, isRequired bool, result *ContractTypeResult) (jen.Code, error) {
	tags := map[string]string{
		"json": name,
	}

	if isRequired {
		tags["validate"] = TagRequired
	}

	var typName string
	bmData, ok := result.BigMaps[getPath(path, name)]
	if !ok {
		typName = result.GetName("BigMap", name)
	} else {
		typName = bmData.Type
	}

	return jen.Id(fieldName("BigMap", name)).Add(jen.Id(typName)).Tag(tags), nil
}

// AsCode -
func (BigMap) AsCode(name, path string, schema api.JSONSchema, result *ContractTypeResult) (Code, error) {
	return Code{
		Statement: nil,
		Name:      result.GetName("BigMap", name),
	}, nil
}

// AsType -
func (BigMap) AsType(name, path string, schema api.JSONSchema, result *ContractTypeResult) (Code, error) {
	typName := result.GetName("BigMap", name)
	return Code{
		Statement: jen.Id(typName),
		Name:      typName,
	}, nil
}
