package app

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/reversersed/AuthService/docs"
	"github.com/reversersed/AuthService/internal/config"
	"github.com/reversersed/AuthService/internal/endpoint"
	"github.com/reversersed/AuthService/internal/service"
	"github.com/reversersed/AuthService/internal/smtp"
	"github.com/reversersed/AuthService/internal/validator"
	"github.com/reversersed/AuthService/pkg/logging/logrus"
	"github.com/reversersed/AuthService/pkg/middleware"
	"github.com/reversersed/AuthService/pkg/shutdown"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	logger.Info("setting up config...")
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
	app.router.Use(gin.LoggerWithWriter(app.log.Writer()))
	app.router.Use(middleware.ErrorHandler)
	app.log.Info("router has been set up")

	app.log.Info("setting up smtp service...")
	smtp := smtp.New(cfg.Smtp, app.log)

	app.log.Info("setting up service...")
	service, err := service.New(app.log, nil, smtp, app.cfg.Server.SecretKey)
	if err != nil {
		return nil, err
	}

	app.log.Info("setting up endpoint...")
	app.handlers = append(app.handlers, endpoint.New(service, app.log, validator.New()))
	app.log.Info("endpoint set up")
	return app, nil
}
func (a *app) Run() error {

	generalRouter := a.router.Group("/api")
	for _, h := range a.handlers {
		h.RegisterRoute(generalRouter)
	}
	if a.cfg.Server.Environment == "debug" {
		a.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	go shutdown.Graceful(a)

	if err := a.router.Run(fmt.Sprintf("%s:%d", a.cfg.Server.Url, a.cfg.Server.Port)); err != nil {
		return err
	}
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
