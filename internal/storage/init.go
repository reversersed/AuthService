package storage

import "github.com/jackc/pgx/v5/pgxpool"

//go:generate mockgen -source=init.go -destination=mocks/storage.go

type logger interface {
	Info(...any)
	Infof(string, ...any)
	Warnf(string, ...any)
}
type storage struct {
	pool   *pgxpool.Pool
	logger logger
}

func New(pool *pgxpool.Pool, logger logger) *storage {
	return &storage{
		pool:   pool,
		logger: logger,
	}
}
