package app

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/reversersed/AuthService/internal/config"
	"github.com/reversersed/AuthService/pkg/logging/logrus"
	"github.com/reversersed/AuthService/pkg/shutdown"
)

// @title API
// @version 1.0

// @host localhost:9000
// @BasePath /api/

// @scheme http
// @accept json

// @securityDefinitions.apiKey ApiKeyAuth
// @in Headers
// @name Authorization
func New() (*app, error) {
	logger, err := logrus.GetLogger()
	if err != nil {
		return nil, err
	}

	cfg, err := config.Load()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	app := &app{
		cfg:    cfg,
		log:    logger,
		router: gin.New(),
	}

	app.log.Info("setting up gin router...")
	gin.SetMode(app.cfg.Server.Environment)
	app.router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PATCH", "DELETE"},
	}))
	app.router.RemoteIPHeaders = []string{"x-forwarded-for", "X-Forwarded-For", "X-FORWARDED-FOR"}
	app.router.Use(gin.LoggerWithWriter(logger.Writer()))
	//app.router.Use(middleware.ErrorHandler)
	app.log.Info("router has been set up")
	return app, nil
}
func (a *app) Run() error {

	generalRouter := a.router.Group("/api")
	for _, h := range a.handlers {
		h.RegisterRoute(generalRouter)
	}

	go shutdown.Graceful(a)
	return nil
}
func (a *app) Close() error {
	for _, h := range a.handlers {
		if err := h.Close(); err != nil {
			return err
		}
	}
	return nil
}
