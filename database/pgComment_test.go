package database

import (
	"context"
	"github.com/dipdup-net/go-lib/mocks"
	"github.com/go-pg/pg/v10"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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
	defer mockCtrl.Finish()

	model := Ballot{}

	// Assert prepare
	expectedParams := toInterfaceSlice([]pg.Safe{"ballots", "Ballot table"})
	mockPgDB.
		EXPECT().
		ExecContext(ctx, "COMMENT ON TABLE ? IS ?",
			gomock.Eq(expectedParams)).
		Times(1).
		Return(nil, nil)

	// Act
	err := makeComments(ctx, pgGo, model)

	// Assert
	assert.Empty(t, err)
}

func TestMakeCommentsWithTableNameWithoutPgComment(t *testing.T) {
	type Ballot struct {
		//nolint
		tableName struct{} `pg:"ballots"`
		Ballot    string   `json:"ballot"`
	}

	mockCtrl, mockPgDB, pgGo, ctx := createPgDbMock(t)
	defer mockCtrl.Finish()

	model := Ballot{}

	// Assert prepare
	expectedParams := toInterfaceSlice([]pg.Safe{"ballots", "Ballot table"})
	mockPgDB.
		EXPECT().
		ExecContext(ctx, "COMMENT ON TABLE ? IS ?",
			gomock.Eq(expectedParams)).
		Times(0).
		Return(nil, nil)

	// Act
	err := makeComments(ctx, pgGo, model)

	// Assert
	assert.Empty(t, err)
}

func TestMakeCommentsFieldWithPgComment(t *testing.T) {
	type Ballot struct {
		//nolint
		tableName struct{} `pg:"ballots"`
		Ballot    string   `json:"ballot" pg-comment:"This is field comment"`
	}

	mockCtrl, mockPgDB, pgGo, ctx := createPgDbMock(t)
	defer mockCtrl.Finish()

	model := Ballot{}

	// Assert prepare
	expectedParams := toInterfaceSlice([]pg.Safe{"ballots", "ballot", "This is field comment"})
	mockPgDB.
		EXPECT().
		ExecContext(ctx, "COMMENT ON COLUMN ?.? IS ?",
			gomock.Eq(expectedParams)).
		Times(1).
		Return(nil, nil)

	// Act
	err := makeComments(ctx, pgGo, model)

	// Assert
	assert.Empty(t, err)
}

func TestMakeCommentsWithTableNameAndFieldsWithPgComment(t *testing.T) {
	type Ballot struct {
		//nolint
		tableName       struct{}    `pg:"ballots" pg-comment:"Ballot table"`
		CreatedAt       int64       `json:"-" pg-comment:"This is field comment"`
		UpdatedAt       int64       `json:"-" pg-comment:"This is field comment"`
		Network         string      `json:"network" pg:",pk" pg-comment:"This is field comment"`
		Hash            string      `json:"hash" pg:",pk" pg-comment:"This is field comment"`
		Branch          string      `json:"branch" pg-comment:"This is field comment"`
		Status          string      `json:"status" pg-comment:"This is field comment"`
		Kind            string      `json:"kind" pg-comment:"This is field comment"`
		Signature       string      `json:"signature" pg-comment:"This is field comment"`
		Protocol        string      `json:"protocol" pg-comment:"This is field comment"`
		Level           uint64      `json:"level" pg-comment:"This is field comment"`
		Errors          interface{} `json:"errors,omitempty" pg:"type:jsonb" pg-comment:"This is field comment"`
		ExpirationLevel *uint64     `json:"expiration_level" pg-comment:"This is field comment"`
		Raw             interface{} `json:"raw,omitempty" pg:"type:jsonb" pg-comment:"This is field comment"`
		Ballot          string      `json:"ballot" pg-comment:"This is field comment"`
		Period          int64       `json:"period" pg-comment:"This is field comment"`
	}

	mockCtrl, mockPgDB, pgGo, ctx := createPgDbMock(t)
	defer mockCtrl.Finish()

	model := Ballot{}

	// Assert prepare
	expectedParams := toInterfaceSlice([]pg.Safe{"ballots", "Ballot table"})
	commentOnTableCall := mockPgDB.
		EXPECT().
		ExecContext(ctx, "COMMENT ON TABLE ? IS ?",
			gomock.Eq(expectedParams)).
		Times(1).
		Return(nil, nil)

	mockPgDB.
		EXPECT().
		ExecContext(ctx, "COMMENT ON COLUMN ?.? IS ?", gomock.Any()).
		Times(15).
		After(commentOnTableCall).
		Return(nil, nil)

	// Act
	err := makeComments(ctx, pgGo, model)

	// Assert
	assert.Empty(t, err)
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
