package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/cristalhq/jwt/v3"
	"github.com/reversersed/AuthService/pkg/middleware"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) GenerateAccessToken(ctx context.Context, guid string, ip string) (string, string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	builder := jwt.NewBuilder(s.signer)

	s.logger.Infof("generating token for %s: %s", ip, guid)
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        guid,
			Audience:  []string{guid, ip},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 60)), // 1 hour live
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		LastIP: ip,
	}
	token, err := builder.Build(claims)
	if err != nil {
		s.logger.Warnf("error while building token: %v", err)
		return "", "", middleware.InternalError(err.Error())
	}
	s.logger.Info("generating refresh token...")
	byteToken := make([]byte, 24)
	rand.Read(byteToken)
	refreshToken := base64.StdEncoding.EncodeToString(byteToken)

	cryptToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return "", "", middleware.InternalError(err.Error())
	}
	if err := s.storage.CreateNewRefreshPassword(ctx, guid, cryptToken, claims.IssuedAt.Time.UTC()); err != nil {
		return "", "", err
	}

	return token.String(), refreshToken, nil
}
func (s *service) ValidateUserToken(ctx context.Context, token string, refresh string, ip string) (*Claims, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	claimToken, err := jwt.ParseAndVerifyString(token, s.verifier)
	if err != nil {
		return nil, middleware.BadRequestError("token not verified: %v", err)
	}

	var claims Claims
	if err := json.Unmarshal(claimToken.RawClaims(), &claims); err != nil {
		return nil, middleware.InternalError(err.Error())
	}

	rowId, hash, err := s.storage.GetFreeRefreshToken(ctx, claims.ID, claims.IssuedAt.Time.UTC())
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword(hash, []byte(refresh)); err != nil {
		s.logger.Warnf("user %s(%s) tried to refresh token with incorrect refresh token: %v (%v)", claims.ID, ip, refresh, hash)
		return nil, middleware.NotFoundError("refresh token is incorrect")
	}
	if err := s.storage.RevokeRefreshToken(ctx, rowId); err != nil {
		return nil, err
	}
	if claims.LastIP != ip {
		s.logger.Infof("client %s tried to refresh token with different ip: %s, but last ip was %s", claims.ID, ip, claims.LastIP)
		s.emailService.SendEmailWarning(ip)
	}
	return &claims, nil
}
