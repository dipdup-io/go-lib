package database

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/dipdup-net/go-lib/config"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/stretchr/testify/suite"
)

const (
	testIndex = "test_index"
)

func newDatabase(ctx context.Context, typ string, cfg config.Database) (Database, error) {
	switch typ {
	case "bun":
		return NewBun(), nil
	default:
		return nil, errors.Errorf("unknown ORM: %s", typ)
	}
}

// DBTestSuite -
type DBTestSuite struct {
	suite.Suite
	psqlContainer *PostgreSQLContainer
	db            Database
	typ           string
}

// SetupSuite -
func (s *DBTestSuite) SetupSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxCancel()

	psqlContainer, err := NewPostgreSQLContainer(ctx, PostgreSQLContainerConfig{
		User:     "user",
		Password: "password",
		Database: "db_test",
		Port:     5432,
		Image:    "postgres:15",
	})
	s.Require().NoError(err)
	s.psqlContainer = psqlContainer

	cfg := config.Database{
		Kind:     config.DBKindPostgres,
		User:     s.psqlContainer.Config.User,
		Database: s.psqlContainer.Config.Database,
		Password: s.psqlContainer.Config.Password,
		Host:     s.psqlContainer.Config.Host,
		Port:     s.psqlContainer.MappedPort().Int(),
	}

	s.db, err = newDatabase(ctx, s.typ, cfg)
	s.Require().NoError(err)
	err = s.db.Connect(ctx, cfg)
	s.Require().NoError(err)
	err = s.db.CreateTable(ctx, &State{}, WithIfNotExists())
	s.Require().NoError(err)
}

func (s *DBTestSuite) TearDownSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	s.Require().NoError(s.db.Close())
	s.Require().NoError(s.psqlContainer.Terminate(ctx))
}

func (s *DBTestSuite) BeforeTest(suiteName, testName string) {
	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory("fixtures"),
	)
	s.Require().NoError(err)
	s.Require().NoError(fixtures.Load())
	s.Require().NoError(db.Close())
}

func (s *DBTestSuite) TestStateCreate() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	newIndex := testIndex + "2"
	newState := State{
		IndexName: newIndex,
		Level:     101,
		Hash:      "beefdead",
	}
	s.Require().NoError(s.db.CreateState(ctx, &newState))

	state, err := s.db.State(ctx, newIndex)
	s.Require().NoError(err)
	s.Require().Equal(state.Level, newState.Level)
	s.Require().Equal(state.Hash, newState.Hash)
	s.Require().Equal(state.UpdatedAt, newState.UpdatedAt)
}

func (s *DBTestSuite) TestStateUpdate() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	newState := State{
		IndexName: testIndex,
		Level:     101,
		Hash:      "beefdead",
	}
	s.Require().NoError(s.db.UpdateState(ctx, &newState))

	state, err := s.db.State(ctx, testIndex)
	s.Require().NoError(err)
	s.Require().Equal(state.Level, newState.Level)
	s.Require().Equal(state.Hash, newState.Hash)
	s.Require().Equal(state.UpdatedAt, newState.UpdatedAt)
}

func (s *DBTestSuite) TestState() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	state, err := s.db.State(ctx, testIndex)
	s.Require().NoError(err)
	s.Require().Equal(state.Level, uint64(100))
	s.Require().Equal(state.Hash, "deadbeef")
	s.Require().Equal(state.UpdatedAt, 1691767900)
	s.Require().Equal(state.IndexName, testIndex)
}

func (s *DBTestSuite) TestDeleteState() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	newState := State{
		IndexName: testIndex,
	}
	s.Require().NoError(s.db.DeleteState(ctx, &newState))

	_, err := s.db.State(ctx, testIndex)
	s.Require().Error(err)
}

func (s *DBTestSuite) TestPing() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	s.Require().NoError(s.db.Ping(ctx))
}

func (s *DBTestSuite) TestMakeComments() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	s.Require().NoError(MakeComments(ctx, s.db, &State{}))
}

func (s *DBTestSuite) TestExec() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	count, err := s.db.Exec(ctx, "delete from dipdup_state where index_name = ?", testIndex)
	s.Require().NoError(err)
	s.Require().EqualValues(1, count)
}

func TestSuite_Run(t *testing.T) {
	ts := new(DBTestSuite)
	for _, typ := range []string{"bun"} {
		ts.typ = typ
	}
	suite.Run(t, ts)
}
