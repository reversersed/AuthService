package app

import (
	"github.com/gin-gonic/gin"
	"github.com/reversersed/AuthService/internal/config"
)

type app struct {
	router   *gin.Engine
	cfg      *config.Config
	handlers []handler
	log      logger
}

type logger interface {
	Infof(string, ...any)
	Info(...any)
	Error(...any)
}
type handler interface {
	RegisterRoute(*gin.RouterGroup)
	Close() error
}
