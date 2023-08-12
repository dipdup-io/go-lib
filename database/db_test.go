package database

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/go-pg/pg/v10/orm"
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
	case "gorm":
		db := NewGorm()
		if err := db.Connect(ctx, cfg); err != nil {
			return nil, err
		}
		if err := db.DB().AutoMigrate(&State{}); err != nil {
			if err := db.Close(); err != nil {
				return nil, err
			}
			return nil, err
		}
		return db, nil
	case "pg-go":
		db := NewPgGo()
		if err := db.Connect(ctx, cfg); err != nil {
			return nil, err
		}
		if err := db.DB().WithContext(ctx).Model(&State{}).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		}); err != nil {
			if err := db.Close(); err != nil {
				return nil, err
			}
			return nil, err
		}
		return db, nil
	case "bun":
		db := NewBun()
		if err := db.Connect(ctx, cfg); err != nil {
			return nil, err
		}
		if _, err := db.DB().NewCreateTable().Model(&State{}).IfNotExists().Exec(ctx); err != nil {
			if err := db.Close(); err != nil {
				return nil, err
			}
			return nil, err
		}
		return db, nil
	default:
		return nil, errors.Errorf("unknown ORM: %s", typ)
	}
}

// TestSuite -
type TestSuite struct {
	suite.Suite
	psqlContainer *PostgreSQLContainer
	db            Database
	typ           string
}

// SetupSuite -
func (s *TestSuite) SetupSuite() {
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

	s.db, err = newDatabase(ctx, s.typ, config.Database{
		Kind:     config.DBKindPostgres,
		User:     s.psqlContainer.Config.User,
		Database: s.psqlContainer.Config.Database,
		Password: s.psqlContainer.Config.Password,
		Host:     s.psqlContainer.Config.Host,
		Port:     s.psqlContainer.MappedPort().Int(),
	})
	s.Require().NoError(err)
}

func (s *TestSuite) TearDownSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	s.Require().NoError(s.db.Close())
	s.Require().NoError(s.psqlContainer.Terminate(ctx))
}

func (s *TestSuite) TestStateCreate() {
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

func (s *TestSuite) TestStateUpdate() {
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

func (s *TestSuite) TestState() {
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

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	state, err := s.db.State(ctx, testIndex)
	s.Require().NoError(err)
	s.Require().Equal(state.Level, uint64(100))
	s.Require().Equal(state.Hash, "deadbeef")
	s.Require().Equal(state.UpdatedAt, 1691767900)
	s.Require().Equal(state.IndexName, testIndex)
}

func (s *TestSuite) TestDeleteState() {
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

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	newState := State{
		IndexName: testIndex,
	}
	s.Require().NoError(s.db.DeleteState(ctx, &newState))

	_, err = s.db.State(ctx, testIndex)
	s.Require().Error(err)
}

func (s *TestSuite) TestPing() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	s.Require().NoError(s.db.Ping(ctx))
}

func (s *TestSuite) TestMakeComments() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	s.Require().NoError(MakeComments(ctx, s.db, &State{}))
}

func TestSuite_Run(t *testing.T) {
	ts := new(TestSuite)
	for _, typ := range []string{"gorm", "pg-go", "bun"} {
		ts.typ = typ
	}
	suite.Run(t, ts)
}
