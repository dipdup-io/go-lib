package database

import (
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/stretchr/testify/mock"
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

func (db *PgDBMock) ExecContext(_ context.Context, query interface{}, params ...interface{}) (pg.Result, error) {
	args := db.Called(query, params)

	return nil, args.Error(0)
}

func TestMakeCommentsWithTableName(t *testing.T) {
	type Ballot struct {
		//nolint
		tableName struct{} `pg:"ballots"`
		Ballot    string   `json:"ballot"`
	}

	pgGo := newPgGoMock()
	ctx := context.Background()
	pgGo.conn.On("ExecContext",
		ctx, mock.Anything, mock.Anything).Return(nil)
	model := Ballot{}

	makeComments(ctx, pgGo, model)

	// assert params of ExecContext
	pgGo.conn.AssertCalled(t, "ExecContext", ctx, `COMMENT ON TABLE ? IS ?`, "ballot", "")
}
