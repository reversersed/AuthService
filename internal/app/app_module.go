package app

import (
	"github.com/gin-gonic/gin"
	"github.com/reversersed/AuthService/internal/config"
	"github.com/reversersed/AuthService/pkg/logging/logrus"
)

type app struct {
	router   *gin.Engine
	cfg      *config.Config
	handlers []handler
	log      *logrus.Logger
}

type handler interface {
	RegisterRoute(*gin.RouterGroup)
	Close() error
}
