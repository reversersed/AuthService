package postgres

import (
	"context"
	"flag"
	"log"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	mock_postgres "github.com/reversersed/AuthService/pkg/postgres/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var cfg *DatabaseConfig

func TestMain(m *testing.M) {
	flag.Parse()
	if testing.Short() {
		log.Println("=== integration tests are not running in short mode")
		return
	}

	ctx := context.Background()
	container, err := postgres.Run(ctx,
		"postgres",
		postgres.WithDatabase("testbase"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpassword"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		log.Fatalf("can't run the container: %v", err)
		return
	}

	host, err := container.Host(ctx)
	if err != nil {
		log.Fatalf("can't get container IP: %v", err)
		return
	}
	port, err := container.MappedPort(ctx, "5432/tcp")
	if err != nil {
		log.Fatalf("can't get container port: %v", err)
		return
	}

	cfg = &DatabaseConfig{
		Host:          host,
		Port:          port.Int(),
		User:          "testuser",
		Password:      "testpassword",
		Database:      "testbase",
		MigrationPath: "../../migrations",
	}
	code := m.Run()

	if err := container.Terminate(ctx); err != nil {
		log.Fatalf("can't terminate container: %v", err)
		return
	}

	os.Exit(code)
}
func TestNewPoolSuccess(t *testing.T) {
	if !assert.NotNil(t, cfg) {
		return
	}
	ctrl := gomock.NewController(t)
	logger := mock_postgres.NewMocklogger(ctrl)
	logger.EXPECT().Info(gomock.Any()).AnyTimes()
	logger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
	pool, err := NewConnectionPool(cfg, logger)
	if assert.NoError(t, err) {
		assert.NotNil(t, pool)
	}
}
func TestNewPoolWithWrongCredentials(t *testing.T) {
	if !assert.NotNil(t, cfg) {
		return
	}

	t.Run("wrong password", func(t *testing.T) {
		cfg := *cfg
		cfg.Password = "wrongpassword"
		ctrl := gomock.NewController(t)
		logger := mock_postgres.NewMocklogger(ctrl)
		logger.EXPECT().Info(gomock.Any()).AnyTimes()
		logger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

		_, err := NewConnectionPool(&cfg, logger)
		assert.Error(t, err)
	})
	t.Run("wrong username", func(t *testing.T) {
		cfg := *cfg
		ctrl := gomock.NewController(t)
		logger := mock_postgres.NewMocklogger(ctrl)
		logger.EXPECT().Info(gomock.Any()).AnyTimes()
		logger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

		cfg.User = "notexistinguser"

		_, err := NewConnectionPool(&cfg, logger)
		assert.Error(t, err)
	})
	t.Run("wrong migrations path", func(t *testing.T) {
		cfg := *cfg
		ctrl := gomock.NewController(t)
		logger := mock_postgres.NewMocklogger(ctrl)
		logger.EXPECT().Info(gomock.Any()).AnyTimes()
		logger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

		cfg.MigrationPath = "/noFolderForMigrations"

		_, err := NewConnectionPool(&cfg, logger)
		assert.Error(t, err)
	})
}
