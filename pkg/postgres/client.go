package postgres

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

type logger interface {
	Info(...any)
	Infof(string, ...any)
}
type DatabaseConfig struct {
	Host     string `env:"POSTGRES_HOST" env-required:"true" env-description:"Postgres hosting address"`
	Port     int    `env:"POSTGRES_PORT" env-required:"true" env-description:"Portgres hosting port"`
	User     string `env:"POSTGRES_USER" env-required:"true" env-description:"Postgres username"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true" env-description:"Postres user password to connect"`
	Database string `env:"POSTGRES_DB" env-required:"true" env-description:"Database name"`
}

func NewConnectionPool(cfg *DatabaseConfig, logger logger) (*pgxpool.Pool, error) {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	pool, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		return nil, err
	}
	pool.Config().AfterConnect = func(ctx context.Context, c *pgx.Conn) error {
		byteToken := make([]byte, 24)
		rand.Read(byteToken)
		trace := base64.StdEncoding.EncodeToString(byteToken)

		logger.Info("database establishing new connection... trace: ", trace)
		<-c.PgConn().CleanupDone()
		logger.Info("database connection cleaned up, trace: ", trace)
		return nil
	}
	err = pool.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	logger.Info("retrieving db from pool to migrate...")
	instance, err := postgres.WithInstance(stdlib.OpenDBFromPool(pool), &postgres.Config{DatabaseName: cfg.Database})
	if err != nil {
		return nil, err
	}

	migrate, err := migrate.NewWithDatabaseInstance("file:///migrations", cfg.Database, instance)
	if err != nil {
		return nil, err
	}

	logger.Info("starting up migration...")
	err = migrate.Up()
	if err != nil {
		return nil, err
	}

	source, err := migrate.Close()
	if source != nil || err != nil {
		return nil, fmt.Errorf("source: %v, database: %v", source, err)
	}
	version, dirty, err := migrate.Version()
	if err != nil {
		return nil, fmt.Errorf("no migrations were applied: %v", err)
	}
	logger.Infof("migrations done, current version: %d, database dirty: %v", version, dirty)
	return pool, nil
}
