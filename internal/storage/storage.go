package storage

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/reversersed/AuthService/pkg/middleware"
)

func (s *storage) CreateNewRefreshPassword(ctx context.Context, uuid string, refreshpassword []byte, creation time.Time) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		s.logger.Warnf("can't begin connection from pool: %v", err)
		tx.Rollback(ctx)
		return middleware.InternalError(err.Error())
	}
	tag, err := tx.Exec(ctx, "INSERT INTO sessions (clientid,refreshtoken,created) VALUES ($1,$2,$3)", uuid, string(refreshpassword), creation.Format("2006-01-02 15:04:05.000000000"))
	if err != nil {
		tx.Rollback(ctx)
		s.logger.Warnf("can't execute query: %v", err)
		return middleware.InternalError(err.Error())
	}
	s.logger.Infof("query %s was successful", tag.String())
	tx.Commit(ctx)
	return nil
}
func (s *storage) GetFreeRefreshToken(ctx context.Context, id string, createdTime time.Time) (string, []byte, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		s.logger.Warnf("can't begin connection from pool: %v", err)
		tx.Rollback(ctx)
		return "", nil, middleware.InternalError(err.Error())
	}
	model := struct {
		id           string
		refreshtoken string
	}{}

	s.logger.Infof("searching token for client %s with time %v", id, createdTime)
	err = tx.QueryRow(ctx, "SELECT id,refreshtoken FROM sessions WHERE clientid = $1 AND created = $2 AND refreshed IS NULL LIMIT 1", id, createdTime.Format("2006-01-02 15:04:05.000000000")).Scan(&model.id, &model.refreshtoken)
	if err == pgx.ErrNoRows {
		tx.Rollback(ctx)
		s.logger.Warnf("can't find token for client %s with timestamp %v", id, createdTime)
		return "", nil, middleware.NotFoundError("no token found: %v", err)
	} else if err != nil {
		tx.Rollback(ctx)
		return "", nil, middleware.InternalError(err.Error())
	}

	if err := tx.Commit(ctx); err != nil {
		tx.Rollback(ctx)
		return "", nil, middleware.InternalError(err.Error())
	}
	return model.id, []byte(model.refreshtoken), nil
}
func (s *storage) RevokeRefreshToken(ctx context.Context, rowId string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return nil
}
