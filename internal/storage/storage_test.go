package storage

import (
	"context"
	"flag"
	"log"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgxpool"
	mock_storage "github.com/reversersed/AuthService/internal/storage/mocks"
	"github.com/reversersed/AuthService/pkg/middleware"
	"github.com/reversersed/AuthService/pkg/postgres"
	"github.com/stretchr/testify/assert"
	pgContainer "github.com/testcontainers/testcontainers-go/modules/postgres"
	"golang.org/x/crypto/bcrypt"
)

var pool *pgxpool.Pool

func TestMain(m *testing.M) {
	flag.Parse()
	if testing.Short() {
		log.Println("=== integration tests are not running in short mode")
		return
	}

	ctx := context.Background()
	container, err := pgContainer.Run(ctx,
		"postgres",
		pgContainer.WithDatabase("testbase"),
		pgContainer.WithUsername("testuser"),
		pgContainer.WithPassword("testpassword"),
	)
	if err != nil {
		log.Fatalf("can't run the container: %v", err)
		return
	}

	host, err := container.ContainerIP(ctx)
	if err != nil {
		log.Fatalf("can't get container IP: %v", err)
		return
	}
	port, err := container.MappedPort(ctx, "5432/tcp")
	if err != nil {
		log.Fatalf("can't get container port: %v", err)
		return
	}

	cfg := &postgres.DatabaseConfig{
		Host:     host,
		Port:     port.Int(),
		User:     "testuser",
		Password: "testpassword",
		Database: "testbase",
	}
	pool, err = postgres.NewConnectionPool(cfg, nil)
	if err != nil {
		log.Fatalf("can't create connection pool: %v", err)
		return
	}
	defer pool.Close()
	code := m.Run()

	if err := container.Terminate(ctx); err != nil {
		log.Fatalf("can't terminate container: %v", err)
		return
	}

	os.Exit(code)
}
func TestDefaultRoute(t *testing.T) {
	if !assert.NotNil(t, pool) {
		return
	}
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	logger := mock_storage.NewMocklogger(ctrl)
	logger.EXPECT().Info(gomock.Any())
	storage := New(pool, logger)

	var (
		userId  = "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"
		refresh = []byte("refreshtesttoken")
		created = uint64(time.Now().UTC().UnixMilli())
	)

	err := storage.CreateNewRefreshPassword(ctx, userId, refresh, created)
	assert.NoError(t, err)

	row, hash, err := storage.GetFreeRefreshToken(ctx, userId, created)
	assert.NoError(t, err)

	err = bcrypt.CompareHashAndPassword(hash, refresh)
	assert.NoError(t, err)

	err = storage.RevokeRefreshToken(ctx, row)
	assert.NoError(t, err)

	_, _, err = storage.GetFreeRefreshToken(ctx, userId, created)
	if assert.Error(t, err) {
		assert.EqualError(t, err, middleware.NotFoundError("no token found: no rows found").Error())
	}
}
