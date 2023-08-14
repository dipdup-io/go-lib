package database

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type (
	PostgreSQLContainer struct {
		testcontainers.Container
		Config     PostgreSQLContainerConfig
		mappedPort nat.Port
	}

	PostgreSQLContainerConfig struct {
		User     string
		Password string
		Database string
		Host     string
		Port     int
		Image    string
	}
)

// NewPostgreSQLContainer -
func NewPostgreSQLContainer(ctx context.Context, cfg PostgreSQLContainerConfig) (*PostgreSQLContainer, error) {
	if cfg.Port == 0 {
		cfg.Port = 5432
	}
	port := fmt.Sprintf("%d/tcp", cfg.Port)

	if cfg.Image == "" {
		cfg.Image = "postgres:15"
	}

	container, err := postgres.RunContainer(ctx,
		testcontainers.WithImage(cfg.Image),
		postgres.WithDatabase(cfg.Database),
		postgres.WithUsername(cfg.User),
		postgres.WithPassword(cfg.Password),
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting host for: %w", err)
	}
	cfg.Host = host

	mappedPort, err := container.MappedPort(ctx, nat.Port(port))
	if err != nil {
		return nil, fmt.Errorf("getting mapped port for %s: %w", port, err)
	}

	return &PostgreSQLContainer{
		Container:  container,
		Config:     cfg,
		mappedPort: mappedPort,
	}, nil
}

// GetDSN -
func (c PostgreSQLContainer) GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", c.Config.User, c.Config.Password, c.Config.Host, c.mappedPort.Port(), c.Config.Database)
}

// MappedPort -
func (c PostgreSQLContainer) MappedPort() nat.Port {
	return c.mappedPort
}
