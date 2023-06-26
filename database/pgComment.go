package database

import (
	"context"
	"github.com/dipdup-net/go-lib/hasura"
	"reflect"
	"strings"
)

func MakeComments(ctx context.Context, sc SchemeCommenter, models ...interface{}) error {
	for _, model := range models {
		modelType := reflect.TypeOf(model)
		var tableName string

		for i := 0; i < modelType.NumField(); i++ {
			fieldType := modelType.Field(i)

			if fieldType.Name == "tableName" {
				var ok bool
				tableName, ok = getPgName(fieldType)
				if !ok {
					tableName = hasura.ToSnakeCase(modelType.Name())
				}

				pgCommentTag, ok := getPgComment(fieldType)
				if !ok {
					continue
				}

				if err := sc.MakeTableComment(ctx, tableName, pgCommentTag); err != nil {
					return err
				}

				continue
			}

			pgCommentTag, ok := getPgComment(fieldType)
			if !ok {
				continue
			}

			columnName, ok := getPgName(fieldType)
			if !ok {
				columnName = hasura.ToSnakeCase(fieldType.Name)
			}

			if err := sc.MakeColumnComment(ctx, tableName, columnName, pgCommentTag); err != nil {
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

func getPgComment(fieldType reflect.StructField) (string, bool) {
	pgCommentTag, ok := fieldType.Tag.Lookup("pg-comment")

	if ok {
		return pgCommentTag, ok
	}

	return "", false
}
