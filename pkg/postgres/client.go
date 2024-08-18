package postgres

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

//go:generate mockgen -source=client.go -destination=mocks/client.go

type logger interface {
	Info(...any)
	Infof(string, ...any)
}
type DatabaseConfig struct {
	Host          string `env:"POSTGRES_HOST" env-required:"true" env-description:"Postgres hosting address"`
	Port          int    `env:"POSTGRES_PORT" env-required:"true" env-description:"Portgres hosting port"`
	User          string `env:"POSTGRES_USER" env-required:"true" env-description:"Postgres username"`
	Password      string `env:"POSTGRES_PASSWORD" env-required:"true" env-description:"Postres user password to connect"`
	Database      string `env:"POSTGRES_DB" env-required:"true" env-description:"Database name"`
	MigrationPath string `env:"POSTGRES_MIGRATIONS_PATH" env-description:"Path to migrations folder" env-default:"/migrations"`
}

func NewConnectionPool(cfg *DatabaseConfig, logger logger) (*pgxpool.Pool, error) {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	pool, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		return nil, err
	}
	pool.Config().MaxConnLifetime = 1 * time.Hour
	pool.Config().MaxConnIdleTime = 30 * time.Minute
	pool.Config().MaxConns = 100
	pool.Config().AfterConnect = func(ctx context.Context, c *pgx.Conn) error {
		byteToken := make([]byte, 24)
		rand.Read(byteToken)
		trace := base64.StdEncoding.EncodeToString(byteToken)

		if logger != nil {
			logger.Info("database establishing new connection... trace: ", trace)
		}
		<-c.PgConn().CleanupDone()
		if logger != nil {
			logger.Info("database connection cleaned up, trace: ", trace)
		}
		return nil
	}
	err = pool.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	if logger != nil {
		logger.Info("retrieving db from pool to migrate...")
	}
	instance, err := postgres.WithInstance(stdlib.OpenDBFromPool(pool), &postgres.Config{DatabaseName: cfg.Database})
	if err != nil {
		return nil, err
	}
	defer instance.Close()
	migrate, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", cfg.MigrationPath), cfg.Database, instance)
	if err != nil {
		return nil, fmt.Errorf("error opening migrations: %v", err)
	}

	if logger != nil {
		logger.Info("starting up migration...")
	}
	err = migrate.Up()
	if err != nil && err.Error() != "no change" {
		return nil, err
	}

	version, dirty, err := migrate.Version()
	if err != nil {
		return nil, fmt.Errorf("no migrations were applied: %v", err)
	}

	source, err := migrate.Close()
	if source != nil || err != nil {
		return nil, fmt.Errorf("source: %v, database: %v", source, err)
	}
	if logger != nil {
		logger.Infof("migrations done, current version: %d, database dirty: %v", version, dirty)
	}
	return pool, nil
}
