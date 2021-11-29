package main

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/dipdup-net/go-lib/cmd/generator/types"
	"github.com/dipdup-net/go-lib/tzkt/api"
	"github.com/iancoleman/strcase"
)

type templateContext struct {
	PackageName     string
	TypeName        string
	EntrypointTypes map[string]types.EntrypointData
	BigMaps         map[string]types.BigMapData
	Contract        string
}

var (
	//go:embed template/*.tmpl
	templates embed.FS
)

// Generate -
func Generate(schema api.ContractJSONSchema, name, contract, dest string) error {
	if dest == "" {
		output, err := os.Getwd()
		if err != nil {
			return err
		}
		dest = output
	}
	destDir := filepath.Join(dest, strcase.ToSnake(name))
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	packageName := strings.ToLower(strcase.ToLowerCamel(name))

	if err := generateDefaultTypes(packageName, destDir); err != nil {
		return err
	}

	result, err := generateContractTypes(schema, packageName, destDir)
	if err != nil {
		return err
	}

	if err := generateContractObject(packageName, contract, destDir, result); err != nil {
		return err
	}

	return nil
}

func generateContractObject(name, contract, dest string, result types.ContractTypeResult) error {
	className := strcase.ToCamel(name)
	return generateFromTemplate("contract", dest, templateContext{
		PackageName:     result.PackageName,
		TypeName:        className,
		Contract:        contract,
		EntrypointTypes: result.Entrypoints,
		BigMaps:         result.BigMaps,
	})
}

func generateContractTypes(schema api.ContractJSONSchema, packageName, dest string) (types.ContractTypeResult, error) {
	result, err := types.GenerateContractTypes(schema, packageName)
	if err != nil {
		return result, err
	}

	f, err := os.OpenFile(filepath.Join(dest, "contract_types.go"), os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return result, err
	}
	if err := result.File.Render(f); err != nil {
		return result, err
	}
	defer f.Close()
	return result, nil
}

func generateDefaultTypes(packageName, dest string) error {
	return generateFromTemplate("types", dest, templateContext{PackageName: packageName})
}

func generateFromTemplate(templateFileName, dest string, ctx interface{}) error {
	tmpl, err := template.New("").Funcs(template.FuncMap{
		"ToLower": strings.ToLower,
		"TrimStorage": func(x string) string {
			return strings.TrimPrefix(x, "storage.")
		},
	}).ParseFS(templates, "template/*")
	if err != nil {
		return err
	}
	targetFile := filepath.Join(dest, fmt.Sprintf("%s.go", templateFileName))
	templateFile, err := os.OpenFile(targetFile, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	if err := tmpl.ExecuteTemplate(templateFile, fmt.Sprintf("%s.tmpl", templateFileName), ctx); err != nil {
		return err
	}

	return templateFile.Close()
}
