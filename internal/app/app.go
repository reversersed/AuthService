package app

import (
	"github.com/reversersed/AuthService/internal/config"
	"github.com/reversersed/AuthService/pkg/logging/logrus"
)

func New() (*app, error) {
	var logger logger
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
		cfg: cfg,
		log: logger,
	}
	return app, nil
}
func (a *app) Run() error {
	return nil
}
