package database

import (
	"context"
	"github.com/dipdup-net/go-lib/mocks"
	"github.com/go-pg/pg/v10"
	"github.com/golang/mock/gomock"
	"testing"
)

type PgGoMock struct {
	conn *mocks.MockPgDB
}

func (p *PgGoMock) DB() PgDB {
	return p.conn
}

func newPgGoMock(pgDB *mocks.MockPgDB) *PgGoMock {
	return &PgGoMock{
		conn: pgDB,
	}
}

func TestMakeCommentsWithTableName(t *testing.T) {
	type Ballot struct {
		//nolint
		tableName struct{} `pg:"ballots" pg-comment:"Ballot table"`
		Ballot    string   `json:"ballot"`
	}

	mockCtrl, mockPgDB, pgGo, ctx := createPgDbMock(t)
	model := Ballot{}

	// Assert params of ExecContext
	expectedParams := toInterfaceSlice([]pg.Safe{"ballots", "Ballot table"})
	mockPgDB.
		EXPECT().
		ExecContext(ctx, "COMMENT ON TABLE ? IS ?",
			gomock.Eq(expectedParams)).
		Return(nil, nil)

	// Act
	makeComments(ctx, pgGo, model)

	// Assert
	defer mockCtrl.Finish()
}

func createPgDbMock(t *testing.T) (*gomock.Controller, *mocks.MockPgDB, *PgGoMock, context.Context) {
	mockCtrl := gomock.NewController(t)
	mockPgDB := mocks.NewMockPgDB(mockCtrl)
	pgGo := newPgGoMock(mockPgDB)
	ctx := context.Background()

	return mockCtrl, mockPgDB, pgGo, ctx
}

func toInterfaceSlice(origin []pg.Safe) []interface{} {
	res := make([]interface{}, len(origin))

	for i := range origin {
		res[i] = origin[i]
	}

	return res
}
