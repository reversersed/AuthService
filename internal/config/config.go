package config

import (
	"fmt"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type ServerConfig struct {
	Environment string `env:"SERVICE_ENVIRONMENT" env-required:"true" env-description="Service environment" env-default="debug"`
	Url         string `env:"SERVICE_HOST_URL" env-required:"true" env-description="Server listening address"`
	Port        int    `env:"SERVICE_HOST_PORT" env-required:"true" env-description="Server listening port"`
}
type Config struct {
	Server *ServerConfig
}

var cfg *Config
var once sync.Once

func Load() (*Config, error) {
	var e error
	once.Do(func() {
		server := &ServerConfig{}

		if err := cleanenv.ReadConfig("config/.env", server); err != nil {
			desc, _ := cleanenv.GetDescription(server, nil)

			e = fmt.Errorf("%v: %s", err, desc)
			return
		}

		cfg = &Config{
			Server: server,
		}
	})
	if e != nil {
		return nil, e
	}
	return cfg, nil
}
