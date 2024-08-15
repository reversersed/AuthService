package postgres

import (
	"context"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DatabaseConfig struct {
	Host     string `env:"POSTGRES_HOST" env-required:"true" env-description:"Postgres hosting address"`
	Port     int    `env:"POSTGRES_PORT" env-required:"true" env-description:"Portgres hosting port"`
	User     string `env:"POSTGRES_USER" env-required:"true" env-description:"Postgres username"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true" env-description:"Postres user password to connect"`
	Database string `env:"POSTGRES_DB" env-required:"true" env-description:"Database name"`
}

func NewConnectionPool(cfg *DatabaseConfig) (*pgxpool.Pool, error) {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	pool, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	migrate, err := migrate.New("file:///migrations", connectionString)
	if err != nil {
		return nil, err
	}
	err = migrate.Up()
	if err != nil {
		return nil, err
	}
	return pool, nil
}
