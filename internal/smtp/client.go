package smtp

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
)

//go:generate mockgen -source=client.go -destination=mocks/client.go

type logger interface {
	Errorf(string, ...any)
	Warnf(string, ...any)
	Infof(string, ...any)
}
type SmtpConfig struct {
	Host     string `env:"SMTP_HOST" env-description:"SMTP hosting" env-default:"-"`
	Port     int    `env:"SMTP_PORT" env-description:"SMTP port" env-default:"0"`
	User     string `env:"SMTP_USER" env-description:"SMTP user"`
	Password string `env:"SMTP_PASSWORD" env-description:"SMTP password"`
}

type service struct {
	config *SmtpConfig
	logger logger
}

func New(cfg *SmtpConfig, logger logger) *service {
	return &service{
		config: cfg,
		logger: logger,
	}
}
func (s *service) SendEmailWarning(ip string) {
	if s.config.Host == "-" {
		s.logger.Warnf("service configuration has not set. canceling email sending...")
		return
	}
	auth := smtp.PlainAuth("", s.config.User, s.config.Password, s.config.Host)
	if len(s.config.Password) == 0 {
		auth = nil
	}
	t, err := template.ParseFiles("templates/email.warning.html")
	if err != nil {
		s.logger.Errorf("can't find or parse html template: %s", err)
		return
	}
	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	_, err = body.Write([]byte(fmt.Sprintf("From: %s \r\nTo: %s \r\nSubject: Выполнен вход в аккаунт \n%s\n\n", s.config.User, s.config.User, mimeHeaders))) // To: user.Email
	if err != nil {
		s.logger.Errorf("can't create email header: %s", err)
		return
	}

	err = t.Execute(&body, struct {
		UserLogin string
		IP        string
	}{
		UserLogin: "client name", // user.Login
		IP:        ip,
	})
	if err != nil {
		s.logger.Errorf("can't create email body: %s", err)
		return
	}

	err = smtp.SendMail(fmt.Sprintf("%s:%d", s.config.Host, s.config.Port), auth, s.config.User, []string{s.config.User}, body.Bytes())
	if err != nil {
		s.logger.Errorf("can't send email: %s", err)
		return
	}
	s.logger.Infof("sent email confirmation message to %s from %s", s.config.User, s.config.User)
}
