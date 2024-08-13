package service

import (
	"github.com/cristalhq/jwt/v3"
)

type logger interface {
	Info(...any)
	Warn(...any)
	Infof(string, ...any)
}
type storage interface {
}

type service struct {
	logger   logger
	signer   jwt.Signer
	verifier jwt.Verifier
	storage  storage
}

func New(logger logger, storage storage, secret string) (*service, error) {
	signer, err := jwt.NewSignerHS(jwt.HS512, []byte(secret))
	if err != nil {
		return nil, err
	}
	verifier, err := jwt.NewVerifierHS(jwt.HS512, []byte(secret))
	if err != nil {
		return nil, err
	}
	return &service{
		logger:   logger,
		signer:   signer,
		verifier: verifier,
		storage:  storage,
	}, nil
}
