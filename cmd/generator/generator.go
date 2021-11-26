package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/dipdup-net/go-lib/cmd/generator/types"
	"github.com/dipdup-net/go-lib/tzkt/api"
	"github.com/iancoleman/strcase"
)

type typesTemplateContext struct {
	PackageName string
}

type contractTemplateContext struct {
	PackageName     string
	TypeName        string
	EntrypointTypes map[string]types.EntrypointData
	Contract        string
}

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
	return generateFromTemplate(result.PackageName, "contract.tmpl", dest, contractTemplateContext{
		PackageName:     result.PackageName,
		TypeName:        className,
		Contract:        contract,
		EntrypointTypes: result.Entrypoints,
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
	return generateFromTemplate(packageName, "types.tmpl", dest, typesTemplateContext{packageName})
}

func generateFromTemplate(packageName, templateFileName, dest string, ctx interface{}) error {
	fileName := fmt.Sprintf("%s.go", strings.TrimSuffix(templateFileName, ".tmpl"))
	buf, err := ioutil.ReadFile(filepath.Join("./templates", templateFileName))
	if err != nil {
		return err
	}
	tmpl, err := template.New(packageName).Parse(string(buf))
	if err != nil {
		return err
	}
	targetFile := filepath.Join(dest, fileName)
	templateFile, err := os.OpenFile(targetFile, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(templateFile)
	if err := tmpl.Execute(w, ctx); err != nil {
		return err
	}
	w.Flush()

	return templateFile.Close()
}
