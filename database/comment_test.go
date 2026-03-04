package database

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
)

func initMocks(t *testing.T) (*gomock.Controller, *MockSchemeCommenter) {
	mockCtrl := gomock.NewController(t)
	mockSchemeCommenter := NewMockSchemeCommenter(mockCtrl)

	return mockCtrl, mockSchemeCommenter
}

func TestMakeCommentsWithTableName(t *testing.T) {
	type Ballot struct {
		//nolint
		tableName struct{} `comment:"Ballot table" pg:"ballots"`
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
		Ballot    string   `comment:"This is field comment" json:"ballot"`
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
		tableName       struct{}    `comment:"Ballot table"          pg:"ballots"`
		CreatedAt       int64       `comment:"This is field comment" json:"-"`
		UpdatedAt       int64       `comment:"This is field comment" json:"-"`
		Network         string      `comment:"This is field comment" json:"network"          pg:",pk"`
		Hash            string      `comment:"This is field comment" json:"hash"             pg:",pk"`
		Branch          string      `comment:"This is field comment" json:"branch"`
		Status          string      `comment:"This is field comment" json:"status"`
		Kind            string      `comment:"This is field comment" json:"kind"`
		Signature       string      `comment:"This is field comment" json:"signature"`
		Protocol        string      `comment:"This is field comment" json:"protocol"`
		Level           uint64      `comment:"This is field comment" json:"level"`
		Errors          interface{} `comment:"This is field comment" json:"errors,omitempty" pg:"type:jsonb"`
		ExpirationLevel *uint64     `comment:"This is field comment" json:"expiration_level"`
		Raw             interface{} `comment:"This is field comment" json:"raw,omitempty"    pg:"type:jsonb"`
		Ballot          string      `comment:"This is field comment" json:"ballot"`
		Period          int64       `comment:"This is field comment" json:"period"`
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
		tableName struct{} `comment:"This multiple table name comment" pg:"ballots"`
		Ballot    string   `comment:"This is multiple field comment"   json:"ballot"`
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

	// Be aware: there is an issue with default order in checking
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
		tableName struct{} `comment:"This multiple table name comment" pg:"ballots"`
		Ballot    string   `comment:"This is multiple field comment"   json:"ballot"`
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

	// Be aware: there is an issue with default order in checking
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
		Ballot    string   `comment:"This is field comment" json:"ballot" pg:"-"`
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
		Ballot    string   `comment:""   json:"ballot"`
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
		CreatedAt int64  `comment:"Date of creation in seconds since UNIX epoch."    json:"-"`
		UpdatedAt int64  `comment:"Date of last update in seconds since UNIX epoch." json:"-"`
		Network   string `comment:"Identifies belonging network."                    json:"network" pg:",pk"`
	}

	type Ballot struct {
		//nolint
		tableName struct{} `comment:"This table name comment" pg:"ballots"`
		Operation
		Ballot string `comment:"This is field comment" json:"ballot"`
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
		CreatedAt int64 `comment:"Date of creation in seconds since UNIX epoch." json:"-"`
	}

	type UpdatedMetadata struct {
		CreatedMetadata
		UpdatedAt int64 `comment:"Date of last update in seconds since UNIX epoch." json:"-"`
	}

	type Operation struct {
		UpdatedMetadata
		Network string `comment:"Identifies belonging network." json:"network" pg:",pk"`
	}

	type Ballot struct {
		//nolint
		tableName struct{} `comment:"This table name comment" pg:"ballots"`
		Operation
		Ballot string `comment:"This is field comment" json:"ballot"`
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
		tableName struct{} `comment:"This embedded type tableName comment."            pg:"operation"`
		CreatedAt int64    `comment:"Date of creation in seconds since UNIX epoch."    json:"-"`
		UpdatedAt int64    `comment:"Date of last update in seconds since UNIX epoch." json:"-"`
		Network   string   `comment:"Identifies belonging network."                    json:"network" pg:",pk"`
	}

	type Ballot struct {
		//nolint
		tableName struct{} `comment:"This table name comment" pg:"ballots"`
		Operation
		Ballot string `comment:"This is field comment" json:"ballot"`
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
		bun.BaseModel `bun:"table:operation"                                      comment:"This is bun comment."          pg:"-"`
		CreatedAt     int64  `comment:"Date of creation in seconds since UNIX epoch."    json:"-"`
		UpdatedAt     int64  `comment:"Date of last update in seconds since UNIX epoch." json:"-"`
		Network       string `bun:",pk"                                                  comment:"Identifies belonging network." json:"network"`
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
