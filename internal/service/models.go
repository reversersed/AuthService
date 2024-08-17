package service

import (
	"github.com/cristalhq/jwt/v3"
)

type Claims struct {
	jwt.RegisteredClaims
	LastIP  string
	Created uint64
}
