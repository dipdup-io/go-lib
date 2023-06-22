package database

import (
	"context"
	"github.com/dipdup-net/go-lib/hasura"
	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

func makeComments(ctx context.Context, conn PgGoConnection, model interface{}) error {
	typ := reflect.TypeOf(model)

	// 1. go through fields
	// 2. if tableName field -
	//		2.1 read value from pg tab, if not exist from model type name with snake case convertion - remember as tableName
	// 		2.2 read value from pg-comment tag, if not exist continue
	//		2.3 set comment with SQL statement
	// 3. other
	//		3.1 read comment from pg-comment, it not exist continue
	// 		3.2 read pg tag first value if it exist, if not - then snake case name of field and set as columnName
	// 		3.3 set comment with SQL statement
	for i := 0; i < typ.NumField(); i++ {
		fieldType := typ.Field(i)
		pgTag, ok := fieldType.Tag.Lookup("pg")
		if !ok {
			continue
		}

		tags := strings.Split(pgTag, ",")

		var name string
		for i := range tags {
			if i == 0 {
				if name == "" {
					name = hasura.ToSnakeCase(fieldType.Name)
				} else {
					name = tags[i]
				}
				continue
			}

			parts := strings.Split(tags[i], ":")
			if parts[0] == "comment" {
				if len(parts) != 2 {
					return errors.Errorf("invalid comments format: %s", pgTag)
				}
				if fieldType.Name == "tableName" {
					// typ.Name() to
					if _, err := conn.DB().ExecContext(ctx, `COMMENT ON TABLE ? IS ?`, pg.Safe(typ.Name()), parts[1]); err != nil {
						return err
					}
				} else {
					if _, err := conn.DB().ExecContext(ctx, `COMMENT ON COLUMN ?.? IS ?`, pg.Safe(typ.Name()), pg.Safe(name), parts[1]); err != nil {
						return err
					}
				}
				continue
			}
		}
	}

	return nil
}
