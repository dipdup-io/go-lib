package database

import (
	"testing"

	"github.com/dipdup-net/go-lib/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
)

func initMocks(t *testing.T) (*gomock.Controller, *mocks.MockSchemeCommenter) {
	mockCtrl := gomock.NewController(t)
	mockSchemeCommenter := mocks.NewMockSchemeCommenter(mockCtrl)

	return mockCtrl, mockSchemeCommenter
}

func TestMakeCommentsWithTableName(t *testing.T) {
	type Ballot struct {
		//nolint
		tableName struct{} `pg:"ballots" comment:"Ballot table"`
		Ballot    string   `json:"ballot"`
	}

	mockCtrl, mockSC := initMocks(t)
	defer mockCtrl.Finish()

	model := Ballot{}

	// Assert prepare
	mockSC.
		EXPECT().
		MakeTableComment(t.Context(), "ballots", "Ballot table").
		Times(1).
		Return(nil)

	// Act
	err := MakeComments(t.Context(), mockSC, model)

	// Assert
	require.NoError(t, err)
}

func TestMakeCommentsWithoutPgComment(t *testing.T) {
	type Ballot struct {
		//nolint
		tableName struct{} `pg:"ballots"`
		Ballot    string   `json:"ballot"`
	}

	mockCtrl, mockSC := initMocks(t)
	defer mockCtrl.Finish()

	model := Ballot{}

	// Assert prepare
	mockSC.
		EXPECT().
		MakeTableComment(gomock.Any(), gomock.Any(), gomock.Any()).
		Times(0)

	mockSC.
		EXPECT().
		MakeColumnComment(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Times(0)

	// Act
	err := MakeComments(t.Context(), mockSC, model)

	// Assert
	require.NoError(t, err)
}

func TestMakeCommentsFieldWithPgComment(t *testing.T) {
	type Ballot struct {
		//nolint
		tableName struct{} `pg:"ballots"`
		Ballot    string   `json:"ballot" comment:"This is field comment"`
	}

	mockCtrl, mockSC := initMocks(t)
	defer mockCtrl.Finish()

	model := Ballot{}

	// Assert prepare
	mockSC.
		EXPECT().
		MakeColumnComment(t.Context(), "ballots", "ballot", "This is field comment").
		Times(1).
		Return(nil)

	// Act
	err := MakeComments(t.Context(), mockSC, model)

	// Assert
	require.NoError(t, err)
}

func TestMakeCommentsWithTableNameAndFieldsWithPgComment(t *testing.T) {
	type Ballot struct {
		//nolint
		tableName       struct{}    `pg:"ballots" comment:"Ballot table"`
		CreatedAt       int64       `json:"-" comment:"This is field comment"`
		UpdatedAt       int64       `json:"-" comment:"This is field comment"`
		Network         string      `json:"network" pg:",pk" comment:"This is field comment"`
		Hash            string      `json:"hash" pg:",pk" comment:"This is field comment"`
		Branch          string      `json:"branch" comment:"This is field comment"`
		Status          string      `json:"status" comment:"This is field comment"`
		Kind            string      `json:"kind" comment:"This is field comment"`
		Signature       string      `json:"signature" comment:"This is field comment"`
		Protocol        string      `json:"protocol" comment:"This is field comment"`
		Level           uint64      `json:"level" comment:"This is field comment"`
		Errors          interface{} `json:"errors,omitempty" pg:"type:jsonb" comment:"This is field comment"`
		ExpirationLevel *uint64     `json:"expiration_level" comment:"This is field comment"`
		Raw             interface{} `json:"raw,omitempty" pg:"type:jsonb" comment:"This is field comment"`
		Ballot          string      `json:"ballot" comment:"This is field comment"`
		Period          int64       `json:"period" comment:"This is field comment"`
	}

	mockCtrl, mockSC := initMocks(t)
	defer mockCtrl.Finish()

	model := Ballot{}

	// Assert prepare
	commentOnTableCall := mockSC.
		EXPECT().
		MakeTableComment(t.Context(), "ballots", "Ballot table").
		Times(1).
		Return(nil)

	mockSC.
		EXPECT().
		MakeColumnComment(t.Context(), "ballots", gomock.Any(), "This is field comment").
		Times(15).
		After(commentOnTableCall).
		Return(nil)

	// Act
	err := MakeComments(t.Context(), mockSC, model)

	// Assert
	require.NoError(t, err)
}

func TestMakeCommentsWithMultipleModels(t *testing.T) {
	type Ballot struct {
		//nolint
		tableName struct{} `pg:"ballots" comment:"This multiple table name comment"`
		Ballot    string   `json:"ballot" comment:"This is multiple field comment"`
	}

	mockCtrl, mockSC := initMocks(t)
	defer mockCtrl.Finish()

	models := []interface{}{Ballot{}, Ballot{}, Ballot{}}

	// Assert prepare
	mockSC.
		EXPECT().
		MakeTableComment(t.Context(), "ballots", "This multiple table name comment").
		Times(3).
		Return(nil)

	// Be aware: there is on issue with default order in checking
	// methods call: https://github.com/golang/mock/issues/653
	mockSC.
		EXPECT().
		MakeColumnComment(t.Context(), "ballots", "ballot", "This is multiple field comment").
		Times(3).
		Return(nil)

	// Act
	err := MakeComments(t.Context(), mockSC, models...)

	// Assert
	require.NoError(t, err)
}

func TestMakeCommentsWithMultipleModelsByPointers(t *testing.T) {
	type Ballot struct {
		//nolint
		tableName struct{} `pg:"ballots" comment:"This multiple table name comment"`
		Ballot    string   `json:"ballot" comment:"This is multiple field comment"`
	}

	mockCtrl, mockSC := initMocks(t)
	defer mockCtrl.Finish()

	models := []interface{}{&Ballot{}, &Ballot{}, &Ballot{}}

	// Assert prepare
	mockSC.
		EXPECT().
		MakeTableComment(t.Context(), "ballots", "This multiple table name comment").
		Times(3).
		Return(nil)

	// Be aware: there is on issue with default order in checking
	// methods call: https://github.com/golang/mock/issues/653
	mockSC.
		EXPECT().
		MakeColumnComment(t.Context(), "ballots", "ballot", "This is multiple field comment").
		Times(3).
		Return(nil)

	// Act
	err := MakeComments(t.Context(), mockSC, models...)

	// Assert
	require.NoError(t, err)
}

func TestMakeCommentsIgnoreFieldWithPgHyphen(t *testing.T) {
	type Ballot struct {
		//nolint
		tableName struct{} `pg:"ballots"`
		Ballot    string   `json:"ballot" pg:"-" comment:"This is field comment"`
	}

	mockCtrl, mockSC := initMocks(t)
	defer mockCtrl.Finish()

	model := Ballot{}

	// Assert prepare
	mockSC.
		EXPECT().
		MakeColumnComment(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Times(0)

	// Act
	err := MakeComments(t.Context(), mockSC, model)

	// Assert
	require.NoError(t, err)
}

func TestMakeCommentsIgnoreFieldsWithEmptyComment(t *testing.T) {
	type Ballot struct {
		//nolint
		tableName struct{} `pg:"ballots"`
		Ballot    string   `json:"ballot" comment:""`
	}

	mockCtrl, mockSC := initMocks(t)
	defer mockCtrl.Finish()

	model := Ballot{}

	// Assert prepare
	mockSC.
		EXPECT().
		MakeColumnComment(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Times(0)

	// Act
	err := MakeComments(t.Context(), mockSC, model)

	// Assert
	require.NoError(t, err)
}

func TestMakeCommentsIgnoreNilModel(t *testing.T) {
	mockCtrl, mockSC := initMocks(t)
	defer mockCtrl.Finish()

	// Assert prepare
	mockSC.
		EXPECT().
		MakeColumnComment(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Times(0)

	// Act
	err := MakeComments(t.Context(), mockSC, nil)

	// Assert
	require.NoError(t, err)
}

func TestMakeCommentsIgnoreNoModels(t *testing.T) {
	mockCtrl, mockSC := initMocks(t)
	defer mockCtrl.Finish()

	// Assert prepare
	mockSC.
		EXPECT().
		MakeColumnComment(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Times(0)

	// Act
	err := MakeComments(t.Context(), mockSC)

	// Assert
	require.NoError(t, err)
}

func TestMakeCommentsWithStructCompositionAndCorrectOrder(t *testing.T) {
	type Operation struct {
		CreatedAt int64  `json:"-" comment:"Date of creation in seconds since UNIX epoch."`
		UpdatedAt int64  `json:"-" comment:"Date of last update in seconds since UNIX epoch."`
		Network   string `json:"network" pg:",pk" comment:"Identifies belonging network."`
	}

	type Ballot struct {
		//nolint
		tableName struct{} `pg:"ballots" comment:"This table name comment"`
		Operation
		Ballot string `json:"ballot" comment:"This is field comment"`
	}

	mockCtrl, mockSC := initMocks(t)
	defer mockCtrl.Finish()

	model := Ballot{}

	// Assert prepare
	tableNameCall := mockSC.
		EXPECT().
		MakeTableComment(t.Context(), "ballots", "This table name comment").
		Times(1).
		Return(nil)

	createdAtCall := mockSC.
		EXPECT().
		MakeColumnComment(t.Context(), "ballots", "created_at", "Date of creation in seconds since UNIX epoch.").
		After(tableNameCall).
		Times(1).
		Return(nil)

	updatedAtCall := mockSC.
		EXPECT().
		MakeColumnComment(t.Context(), "ballots", "updated_at", "Date of last update in seconds since UNIX epoch.").
		After(createdAtCall).
		Times(1).
		Return(nil)

	networkCall := mockSC.
		EXPECT().
		MakeColumnComment(t.Context(), "ballots", "network", "Identifies belonging network.").
		After(updatedAtCall).
		Times(1).
		Return(nil)

	mockSC.
		EXPECT().
		MakeColumnComment(t.Context(), "ballots", "ballot", "This is field comment").
		After(networkCall).
		Times(1).
		Return(nil)

	// Act
	err := MakeComments(t.Context(), mockSC, model)

	// Assert
	require.NoError(t, err)
}

func TestMakeCommentsWithDeepStructComposition(t *testing.T) {
	type CreatedMetadata struct {
		CreatedAt int64 `json:"-" comment:"Date of creation in seconds since UNIX epoch."`
	}

	type UpdatedMetadata struct {
		CreatedMetadata
		UpdatedAt int64 `json:"-" comment:"Date of last update in seconds since UNIX epoch."`
	}

	type Operation struct {
		UpdatedMetadata
		Network string `json:"network" pg:",pk" comment:"Identifies belonging network."`
	}

	type Ballot struct {
		//nolint
		tableName struct{} `pg:"ballots" comment:"This table name comment"`
		Operation
		Ballot string `json:"ballot" comment:"This is field comment"`
	}

	mockCtrl, mockSC := initMocks(t)
	defer mockCtrl.Finish()

	model := Ballot{}

	// Assert prepare
	tableNameCall := mockSC.
		EXPECT().
		MakeTableComment(t.Context(), "ballots", "This table name comment").
		Times(1).
		Return(nil)

	createdAtCall := mockSC.
		EXPECT().
		MakeColumnComment(t.Context(), "ballots", "created_at", "Date of creation in seconds since UNIX epoch.").
		After(tableNameCall).
		Times(1).
		Return(nil)

	updatedAtCall := mockSC.
		EXPECT().
		MakeColumnComment(t.Context(), "ballots", "updated_at", "Date of last update in seconds since UNIX epoch.").
		After(createdAtCall).
		Times(1).
		Return(nil)

	networkCall := mockSC.
		EXPECT().
		MakeColumnComment(t.Context(), "ballots", "network", "Identifies belonging network.").
		After(updatedAtCall).
		Times(1).
		Return(nil)

	mockSC.
		EXPECT().
		MakeColumnComment(t.Context(), "ballots", "ballot", "This is field comment").
		After(networkCall).
		Times(1).
		Return(nil)

	// Act
	err := MakeComments(t.Context(), mockSC, model)

	// Assert
	require.NoError(t, err)
}

func TestMakeCommentsWithStructCompositionErrorOnEmbeddedTableName(t *testing.T) {
	type Operation struct {
		//nolint
		tableName struct{} `pg:"operation" comment:"This embedded type tableName comment."`
		CreatedAt int64    `json:"-" comment:"Date of creation in seconds since UNIX epoch."`
		UpdatedAt int64    `json:"-" comment:"Date of last update in seconds since UNIX epoch."`
		Network   string   `json:"network" pg:",pk" comment:"Identifies belonging network."`
	}

	type Ballot struct {
		//nolint
		tableName struct{} `pg:"ballots" comment:"This table name comment"`
		Operation
		Ballot string `json:"ballot" comment:"This is field comment"`
	}

	mockCtrl, mockSC := initMocks(t)
	defer mockCtrl.Finish()

	model := Ballot{}

	// Assert prepare
	mockSC.
		EXPECT().
		MakeTableComment(t.Context(), "ballots", "This table name comment").
		Times(1).
		Return(nil)

	// Act
	err := MakeComments(t.Context(), mockSC, model)

	// Assert
	require.Error(t, err, "Embedded type must not have tableName field.")
}

func TestMakeCommentsWithBunBaseModel(t *testing.T) {
	type Operation struct {
		bun.BaseModel `pg:"-" bun:"table:operation" comment:"This is bun comment."`
		CreatedAt     int64  `json:"-" comment:"Date of creation in seconds since UNIX epoch."`
		UpdatedAt     int64  `json:"-" comment:"Date of last update in seconds since UNIX epoch."`
		Network       string `json:"network" bun:",pk" comment:"Identifies belonging network."`
	}

	mockCtrl, mockSC := initMocks(t)
	defer mockCtrl.Finish()

	// Assert prepare
	tableNameCall := mockSC.
		EXPECT().
		MakeTableComment(t.Context(), "operation", "This is bun comment.").
		Times(1).
		Return(nil)

	createdAtCall := mockSC.
		EXPECT().
		MakeColumnComment(t.Context(), "operation", "created_at", "Date of creation in seconds since UNIX epoch.").
		After(tableNameCall).
		Times(1).
		Return(nil)

	updatedAtCall := mockSC.
		EXPECT().
		MakeColumnComment(t.Context(), "operation", "updated_at", "Date of last update in seconds since UNIX epoch.").
		After(createdAtCall).
		Times(1).
		Return(nil)

	mockSC.
		EXPECT().
		MakeColumnComment(t.Context(), "operation", "network", "Identifies belonging network.").
		After(updatedAtCall).
		Times(1).
		Return(nil)

	// Act
	err := MakeComments(t.Context(), mockSC, Operation{})

	// Assert
	require.NoError(t, err, "Bun model comments was failed")
}
