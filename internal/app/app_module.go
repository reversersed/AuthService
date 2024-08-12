package app

import "github.com/reversersed/AuthService/internal/config"

type app struct {
	cfg *config.Config
	log logger
}

type logger interface {
	Infof(string, ...any)
	Info(...any)
	Error(...any)
}
