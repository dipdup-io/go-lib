package database

import (
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/stretchr/testify/mock"
	"sort"
	"testing"
)

type PgGoMock struct {
	conn *PgDBMock
}

func (p *PgGoMock) DB() PgDB {
	return p.conn
}

func newPgGoMock() *PgGoMock {
	return &PgGoMock{
		conn: &PgDBMock{},
	}
}

type PgDBMock struct {
	mock.Mock
}

func (db *PgDBMock) ExecContext(ctx context.Context, query interface{}, params ...interface{}) (pg.Result, error) {
	args := db.Called(ctx, query, params)

	return nil, args.Error(0)
}

func TestMakeCommentsWithTableName(t *testing.T) {
	type Ballot struct {
		//nolint
		tableName struct{} `pg:"ballots" pg-comment:"Ballot table"`
		Ballot    string   `json:"ballot"`
	}

	pgGo := newPgGoMock()
	ctx := context.Background()
	pgGo.conn.On("ExecContext",
		ctx, mock.Anything, mock.Anything).Return(nil)
	model := Ballot{}

	makeComments(ctx, pgGo, model)

	// assert params of ExecContext
	pgGo.conn.AssertCalled(t, "ExecContext",
		ctx, "COMMENT ON TABLE ? IS ?", []pg.Safe{"ballots", "Ballot table"})

	//paramsMatcher := mock.MatchedBy(func(params []string) bool {
	//	return IsEqual(params, []string{"ballots", "Ballot table"})
	//})
	//pgGo.conn.AssertCalled(t, "ExecContext", paramsMatcher)
}

func IsEqual(a1 []string, a2 []string) bool {
	sort.Strings(a1)
	sort.Strings(a2)
	if len(a1) == len(a2) {
		for i, v := range a1 {
			if v != a2[i] {
				return false
			}
		}
	} else {
		return false
	}
	return true
}
