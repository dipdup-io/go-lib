package database

import (
	"context"
	"testing"
	"time"

	"github.com/dipdup-net/go-lib/config"
	"github.com/stretchr/testify/suite"
	"github.com/uptrace/bun"
)

type logs struct {
	bun.BaseModel `bun:"logs"`

	Id        int       `pg:"id,pk"`
	LogString string    `pg:"log_string"`
	LogTime   time.Time `pg:"log_time"`
}

// PartitionTestSuite -
type PartitionTestSuite struct {
	suite.Suite
	psqlContainer *PostgreSQLContainer
	db            *Bun
	pm            RangePartitionManager
}

// SetupSuite -
func (s *PartitionTestSuite) SetupSuite() {
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

	s.db = NewBun()
	err = s.db.Connect(ctx, config.Database{
		Kind:     config.DBKindPostgres,
		User:     s.psqlContainer.Config.User,
		Database: s.psqlContainer.Config.Database,
		Password: s.psqlContainer.Config.Password,
		Host:     s.psqlContainer.Config.Host,
		Port:     s.psqlContainer.MappedPort().Int(),
	})
	s.Require().NoError(err)

	_, err = s.db.DB().NewCreateTable().Model(&logs{}).PartitionBy("RANGE(log_time)").IfNotExists().Exec(ctx)
	s.Require().NoError(err)

	s.pm = NewPartitionManager(s.db, PartitionByMonth)
}

// TearDownSuite -
func (s *PartitionTestSuite) TearDownSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	s.Require().NoError(s.db.Close())
	s.Require().NoError(s.psqlContainer.Terminate(ctx))
}

// TestParitioning -
func (s *PartitionTestSuite) TestParitioning() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	for _, l := range []logs{
		{
			Id:        1,
			LogTime:   time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			LogString: "test 1",
		},
		{
			Id:        2,
			LogTime:   time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
			LogString: "test 2",
		},
		{
			Id:        3,
			LogTime:   time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC),
			LogString: "test 3",
		},
	} {
		err := s.pm.CreatePartition(ctx, l.LogTime, "logs")
		s.Require().NoError(err)

		_, err = s.db.DB().NewInsert().Model(&l).Exec(ctx)
		s.Require().NoError(err)
	}
}

func TestSuitePartition_Run(t *testing.T) {
	suite.Run(t, new(PartitionTestSuite))
}
