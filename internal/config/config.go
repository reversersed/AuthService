package config

import (
	"fmt"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/reversersed/AuthService/internal/smtp"
)

type ServerConfig struct {
	Environment string `env:"SERVICE_ENVIRONMENT" env-description:"Service environment" env-default:"debug"`
	Url         string `env:"SERVICE_HOST_URL" env-required:"true" env-description:"Server listening address"`
	Port        int    `env:"SERVICE_HOST_PORT" env-required:"true" env-description:"Server listening port"`
	SecretKey   string `env:"JWT_SECRET_KEY" env-required:"true" env-description:"Secret key for JWT authentication"`
}
type Config struct {
	Server *ServerConfig
	Smtp   *smtp.SmtpConfig
}

var cfg *Config
var once sync.Once

func Load() (*Config, error) {
	var e error
	once.Do(func() {
		server := &ServerConfig{}
		smtp := &smtp.SmtpConfig{}

		if err := cleanenv.ReadConfig("config/.env", server); err != nil {
			desc, _ := cleanenv.GetDescription(cfg, nil)

			e = fmt.Errorf("%v: %s", err, desc)
			return
		}
		if err := cleanenv.ReadConfig("config/.env", smtp); err != nil {
			desc, _ := cleanenv.GetDescription(cfg, nil)

			e = fmt.Errorf("%v: %s", err, desc)
			return
		}

		cfg = &Config{
			Server: server,
			Smtp:   smtp,
		}
	})
	if e != nil {
		return nil, e
	}
	return cfg, nil
}
