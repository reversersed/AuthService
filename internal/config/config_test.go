package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/reversersed/AuthService/internal/smtp"
	"github.com/reversersed/AuthService/pkg/postgres"
	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	dir := t.TempDir()
	file, _ := os.OpenFile(fmt.Sprintf("%s/.env", dir), os.O_CREATE|os.O_TRUNC, os.ModeAppend)

	file.WriteString(`
SERVICE_ENVIRONMENT = debug
SERVICE_HOST_URL = localhost
SERVICE_HOST_PORT = 1001
JWT_SECRET_KEY = jwtsecret

SMTP_HOST = smtphost
SMTP_PORT = 587
SMTP_USER = smtpuser
SMTP_PASSWORD = smtppass

POSTGRES_HOST = db
POSTGRES_PORT = 1000
POSTGRES_PASSWORD = dbpass
POSTGRES_USER = root
POSTGRES_DB = base`)

	cfg, e := Load(fmt.Sprintf("%s/.env", dir))

	file.Close()

	assert.NoError(t, e)
	if assert.NotNil(t, cfg) {
		excepted := &Config{
			Database: &postgres.DatabaseConfig{
				Host:     "db",
				Port:     1000,
				Password: "dbpass",
				User:     "root",
				Database: "base",
			},
			Server: &ServerConfig{
				Environment: "debug",
				Url:         "localhost",
				Port:        1001,
				SecretKey:   "jwtsecret",
			},
			Smtp: &smtp.SmtpConfig{
				Host:     "smtphost",
				Port:     587,
				User:     "smtpuser",
				Password: "smtppass",
			},
		}

		assert.EqualValues(t, excepted, cfg)
	}
}
