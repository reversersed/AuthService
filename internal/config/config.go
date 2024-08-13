package config

import (
	"fmt"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Environment string `env:"SERVICE_ENVIRONMENT" env-required:"true" env-description="Service environment" env-default="debug"`
	Url         string `env:"SERVICE_HOST_URL" env-required:"true" env-description="Server listening address"`
	Port        int    `env:"SERVICE_HOST_PORT" env-required:"true" env-description="Server listening port"`
	SecretKey   string `env:"JWT_SECRET_KEY" env-required:"true" env-description="Secret key for JWT authentication"`
}

var cfg *Config
var once sync.Once

func Load() (*Config, error) {
	var e error
	once.Do(func() {
		cfg = &Config{}

		if err := cleanenv.ReadConfig("config/.env", cfg); err != nil {
			desc, _ := cleanenv.GetDescription(cfg, nil)

			e = fmt.Errorf("%v: %s", err, desc)
			return
		}
	})
	if e != nil {
		return nil, e
	}
	return cfg, nil
}
