package service

import (
	"time"

	"github.com/cristalhq/jwt/v3"
)

//go:generate mockgen -source=init.go -destination=mocks/service.go

type logger interface {
	Info(...any)
	Warnf(string, ...any)
	Infof(string, ...any)
}
type storage interface {
	CreateNewRefreshPassword(string, []byte, time.Time) error
	GetFreeRefreshToken(string, time.Time) (string, []byte, error)
	RevokeRefreshToken(string) error
}
type emailService interface {
	SendEmailWarning(ip string)
}
type service struct {
	logger       logger
	signer       jwt.Signer
	verifier     jwt.Verifier
	storage      storage
	emailService emailService
}

func New(logger logger, storage storage, emailService emailService, secret string) (*service, error) {
	signer, err := jwt.NewSignerHS(jwt.HS512, []byte(secret))
	if err != nil {
		return nil, err
	}
	verifier, err := jwt.NewVerifierHS(jwt.HS512, []byte(secret))
	if err != nil {
		return nil, err
	}
	return &service{
		logger:       logger,
		signer:       signer,
		verifier:     verifier,
		storage:      storage,
		emailService: emailService,
	}, nil
}
