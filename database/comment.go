package database

import (
	"context"
	"github.com/dipdup-net/go-lib/hasura"
	"reflect"
	"strings"
)

func MakeComments(ctx context.Context, sc SchemeCommenter, models ...interface{}) error {
	if models == nil {
		return nil
	}

	for _, model := range models {
		if model == nil {
			continue
		}

		modelType := reflect.TypeOf(model)

		if reflect.ValueOf(model).Kind() == reflect.Ptr {
			modelType = modelType.Elem()
		}

		var tableName string

		for i := 0; i < modelType.NumField(); i++ {
			fieldType := modelType.Field(i)

			if fieldType.Name == "tableName" {
				var ok bool
				tableName, ok = getPgName(fieldType)
				if !ok {
					tableName = hasura.ToSnakeCase(modelType.Name())
				}

				comment, ok := getComment(fieldType)
				if !ok {
					continue
				}

				if err := sc.MakeTableComment(ctx, tableName, comment); err != nil {
					return err
				}

				continue
			}

			comment, ok := getComment(fieldType)
			if !ok || comment == "" {
				continue
			}

			columnName, ok := getPgName(fieldType)

			if columnName == "-" {
				continue
			}

			if !ok {
				columnName = hasura.ToSnakeCase(fieldType.Name)
			}

			if err := sc.MakeColumnComment(ctx, tableName, columnName, comment); err != nil {
				return err
			}
		}
	}
	return nil
}

func getPgName(fieldType reflect.StructField) (name string, ok bool) {
	pgTag, ok := fieldType.Tag.Lookup("pg")
	if !ok {
		return "", false
	}

	tags := strings.Split(pgTag, ",")

	if tags[0] != "" {
		name = tags[0]
		return name, ok
	}

	return "", false
}

func getComment(fieldType reflect.StructField) (string, bool) {
	commentTag, ok := fieldType.Tag.Lookup("comment")

	if ok {
		return commentTag, ok
	}

	return "", false
}
