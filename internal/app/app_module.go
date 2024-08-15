package app

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/reversersed/AuthService/internal/config"
	"github.com/reversersed/AuthService/pkg/logging/logrus"
)

type app struct {
	router   *gin.Engine
	cfg      *config.Config
	handlers []handler
	log      *logrus.Logger
	dbPool   *pgxpool.Pool
}

type handler interface {
	RegisterRoute(*gin.RouterGroup)
	Close() error
}
