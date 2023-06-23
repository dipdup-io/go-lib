package database

import (
	"context"
	"github.com/dipdup-net/go-lib/hasura"
	"github.com/go-pg/pg/v10"
	"reflect"
	"strings"
)

func makeComments(ctx context.Context, conn PgGoConnection, model interface{}) error {
	modelType := reflect.TypeOf(model)
	var tableName pg.Safe

	for i := 0; i < modelType.NumField(); i++ {
		fieldType := modelType.Field(i)

		if fieldType.Name == "tableName" {
			var ok bool
			tableName, ok = getPgName(fieldType)
			if !ok {
				tableName = pg.Safe(hasura.ToSnakeCase(modelType.Name()))
			}

			pgCommentTag, ok := getPgComment(fieldType)
			if !ok {
				continue
			}

			if _, err := conn.DB().ExecContext(ctx,
				`COMMENT ON TABLE ? IS ?`,
				tableName, pgCommentTag); err != nil {
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
			columnName = pg.Safe(hasura.ToSnakeCase(fieldType.Name))
		}

		if _, err := conn.DB().ExecContext(ctx,
			`COMMENT ON COLUMN ?.? IS ?`,
			tableName, columnName, pgCommentTag); err != nil {
			return err
		}
	}

	return nil
}

func getPgName(fieldType reflect.StructField) (name pg.Safe, ok bool) {
	pgTag, ok := fieldType.Tag.Lookup("pg")
	if !ok {
		return "", false
	}

	tags := strings.Split(pgTag, ",")

	if tags[0] != "" {
		name = pg.Safe(tags[0])
		return name, ok
	}

	return "", false
}

func getPgComment(fieldType reflect.StructField) (pg.Safe, bool) {
	pgCommentTag, ok := fieldType.Tag.Lookup("pg-comment")

	if ok {
		return pg.Safe(pgCommentTag), ok
	}

	return "", false
}
