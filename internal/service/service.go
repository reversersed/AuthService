package service

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/cristalhq/jwt/v3"
	"github.com/reversersed/AuthService/pkg/middleware"
)

func (s *service) GenerateAccessToken(guid string, ip string) (string, string, error) {
	builder := jwt.NewBuilder(s.signer)

	s.logger.Infof("generating token for %s: %s", ip, guid)
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        guid,
			Audience:  []string{guid, ip},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 60)), // 1 hour live
		},
		LastIP: ip,
	}
	token, err := builder.Build(claims)
	if err != nil {
		s.logger.Warn(err)
		return "", "", middleware.InternalError(err.Error())
	}
	s.logger.Info("generating refresh token...")
	byteToken := make([]byte, 24)
	rand.Read(byteToken)

	return token.String(), base64.StdEncoding.EncodeToString(byteToken), nil
}
func (s *service) ValidateUserToken(token string, ip string) (*Claims, error) {
	claimToken, err := jwt.ParseAndVerifyString(token, s.verifier)
	if err != nil {
		return nil, middleware.BadRequestError("token not verified: %v", err)
	}

	var claims Claims
	if err := json.Unmarshal(claimToken.RawClaims(), &claims); err != nil {
		return nil, middleware.InternalError(err.Error())
	}

	if claims.LastIP != ip {
		s.logger.Infof("client %s tried to refresh token with different ip: %s, but last ip was %s", claims.ID, ip, claims.LastIP)
		//sent email via smtp host...
		return nil, middleware.ConfictError("token's saved ip does not equal to user ip. sending verification email...")
	}
	return &claims, nil
}
