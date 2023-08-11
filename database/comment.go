package database

import (
	"context"
	"reflect"
	"strings"

	"github.com/dipdup-net/go-lib/hasura"
	"github.com/pkg/errors"
)

// MakeComments -
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

			if fieldType.Name == "tableName" || fieldType.Name == "BaseModel" {
				var ok bool
				tableName, ok = getDatabaseTagName(fieldType)
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

			if fieldType.Anonymous {
				if err := makeEmbeddedComments(ctx, sc, tableName, fieldType.Type); err != nil {
					return err
				}
				continue
			}

			if err := makeFieldComment(ctx, sc, tableName, fieldType); err != nil {
				return err
			}
		}
	}
	return nil
}

func makeEmbeddedComments(ctx context.Context, sc SchemeCommenter, tableName string, t reflect.Type) error {
	for i := 0; i < t.NumField(); i++ {
		fieldType := t.Field(i)

		if fieldType.Anonymous {
			if err := makeEmbeddedComments(ctx, sc, tableName, fieldType.Type); err != nil {
				return err
			}

			continue
		}

		if fieldType.Name == "tableName" {
			return errors.New("Embedded type must not have tableName field.")
		}

		if err := makeFieldComment(ctx, sc, tableName, fieldType); err != nil {
			return err
		}
	}

	return nil
}

func makeFieldComment(ctx context.Context, sc SchemeCommenter, tableName string, fieldType reflect.StructField) error {
	comment, ok := getComment(fieldType)
	if !ok || comment == "" {
		return nil
	}

	columnName, ok := getDatabaseTagName(fieldType)

	if columnName == "-" {
		return nil
	}

	if !ok {
		columnName = hasura.ToSnakeCase(fieldType.Name)
	}

	if err := sc.MakeColumnComment(ctx, tableName, columnName, comment); err != nil {
		return err
	}

	return nil
}

func getDatabaseTagName(fieldType reflect.StructField) (name string, ok bool) {
	pgTag, pgOk := fieldType.Tag.Lookup("pg")
	bunTag, bunOk := fieldType.Tag.Lookup("bun")
	ok = pgOk || bunOk

	var tag string
	switch {
	case !pgOk && !bunOk:
		return "", false
	case pgOk && pgTag != "-":
		tag = pgTag
	case bunOk && bunTag != "-":
		tag = strings.TrimPrefix(bunTag, "table:")
	case pgOk:
		tag = pgTag
	case bunOk:
		tag = bunTag
	}

	tags := strings.Split(tag, ",")

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
