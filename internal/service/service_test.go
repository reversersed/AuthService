package service

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	mock_service "github.com/reversersed/AuthService/internal/service/mocks"
	"github.com/reversersed/AuthService/pkg/middleware"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestGeneratingToken(t *testing.T) {
	table := []struct {
		Name          string
		Secret        string
		Uuid          string
		Ip            string
		Err           error
		MockBehaviour func(*mock_service.MockemailService, *mock_service.Mockstorage)
	}{
		{
			Name:   "empty secret",
			Secret: "",
			Err:    errors.New("jwt: key is nil"),
		},
		{
			Name:   "empty guid",
			Secret: "secret",
			Err:    middleware.BadRequestError("bad id"),
			Uuid:   "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
			Ip:     "127.0.0.1",
			MockBehaviour: func(ms *mock_service.MockemailService, m *mock_service.Mockstorage) {
				m.EXPECT().CreateNewRefreshPassword(gomock.Any(), "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", gomock.Any(), gomock.Any()).Return(middleware.BadRequestError("bad id"))
			},
		},
		{
			Name:   "successful generating",
			Secret: "secret",
			Uuid:   "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
			Ip:     "127.0.0.1",
			MockBehaviour: func(ms *mock_service.MockemailService, m *mock_service.Mockstorage) {
				m.EXPECT().CreateNewRefreshPassword(gomock.Any(), "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", gomock.Any(), gomock.Any()).Return(nil)
			},
		},
	}

	for _, v := range table {
		t.Run(v.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			logger := mock_service.NewMocklogger(ctrl)
			logger.EXPECT().Info(gomock.Any()).AnyTimes()
			logger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
			logger.EXPECT().Warnf(gomock.Any(), gomock.Any()).AnyTimes()
			storage := mock_service.NewMockstorage(ctrl)
			emailService := mock_service.NewMockemailService(ctrl)

			if v.MockBehaviour != nil {
				v.MockBehaviour(emailService, storage)
			}
			service, err := New(logger, storage, emailService, v.Secret)
			if v.Err != nil && len(v.Secret) == 0 {
				assert.Error(t, err)
				assert.EqualError(t, err, v.Err.Error())
				return
			} else {
				assert.NoError(t, err)
			}

			token, refresh, err := service.GenerateAccessToken(context.Background(), v.Uuid, v.Ip)
			if v.Err != nil && assert.Error(t, err) {
				assert.EqualError(t, err, v.Err.Error())
			} else if v.Err == nil && assert.NoError(t, err) {
				assert.NotEmpty(t, token)
				assert.NotEmpty(t, refresh)
			}
		})
	}
}

func TestTokenValidation(t *testing.T) {
	table := []struct {
		Name          string
		Secret        string
		Token         string
		Refresh       string
		Ip            string
		ExceptedID    string
		Err           error
		MockBehaviour func(*mock_service.MockemailService, *mock_service.Mockstorage)
	}{
		{
			Name:   "validation error",
			Secret: "secret",
			Token:  "notvalidtoken",
			Err:    middleware.BadRequestError("token not verified: jwt: token format is not valid"),
		},
		{
			Name:   "wrong algorithm token",
			Secret: "secret",
			Token:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			Err:    middleware.BadRequestError("token not verified: jwt: token is signed by another algorithm"),
		},
		{
			Name:   "wrong signature token",
			Secret: "secret",
			Token:  "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.VFb0qJ1LRg_4ujbZoRMXnVkUgiuKq5KxWqNdbKq_G9Vvz-S1zZa9LPxtHWKa64zDl2ofkT8F6jBt_K4riU-fPg",
			Err:    middleware.BadRequestError("token not verified: jwt: signature is not valid"),
		},
		{
			Name:    "refresh token error",
			Secret:  "secret",
			Token:   "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJqdGkiOiJhMGVlYmM5OS05YzBiLTRlZjgtYmI2ZC02YmI5YmQzODBhMTEiLCJhdWQiOlsiYTBlZWJjOTktOWMwYi00ZWY4LWJiNmQtNmJiOWJkMzgwYTExIiwiMTI3LjAuMC4xIl0sImV4cCI6MTcyMzkyMjEyNiwiTGFzdElQIjoiMTI3LjAuMC4xIiwiQ3JlYXRlZCI6MTcyMzkxODUyNjEyNn0.eDR5mwPC5OaUBw4ACBJ3B7Zb0-rLKvzYG_USwW-m3W2jPT-y8XvwaAq-TvO4xli6g6QMljnFPjbSpFrzajw46w",
			Refresh: "refreshtoken",
			MockBehaviour: func(ms *mock_service.MockemailService, m *mock_service.Mockstorage) {
				m.EXPECT().GetFreeRefreshToken(gomock.Any(), "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", uint64(1723918526126)).Return("", nil, middleware.NotFoundError("token not found"))
			},
			Err: middleware.NotFoundError("token not found"),
		},
		{
			Name:    "refresh hash validation error",
			Secret:  "secret",
			Token:   "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJqdGkiOiJhMGVlYmM5OS05YzBiLTRlZjgtYmI2ZC02YmI5YmQzODBhMTEiLCJhdWQiOlsiYTBlZWJjOTktOWMwYi00ZWY4LWJiNmQtNmJiOWJkMzgwYTExIiwiMTI3LjAuMC4xIl0sImV4cCI6MTcyMzkyMjEyNiwiTGFzdElQIjoiMTI3LjAuMC4xIiwiQ3JlYXRlZCI6MTcyMzkxODUyNjEyNn0.eDR5mwPC5OaUBw4ACBJ3B7Zb0-rLKvzYG_USwW-m3W2jPT-y8XvwaAq-TvO4xli6g6QMljnFPjbSpFrzajw46w",
			Refresh: "refreshtoken",
			MockBehaviour: func(ms *mock_service.MockemailService, m *mock_service.Mockstorage) {
				refresh, _ := bcrypt.GenerateFromPassword([]byte("wrongtoken"), bcrypt.DefaultCost)
				m.EXPECT().GetFreeRefreshToken(gomock.Any(), "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", uint64(1723918526126)).Return("rowid", refresh, nil)
			},
			Err: middleware.NotFoundError("refresh token is incorrect"),
		},
		{
			Name:    "revoking error",
			Secret:  "secret",
			Token:   "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJqdGkiOiJhMGVlYmM5OS05YzBiLTRlZjgtYmI2ZC02YmI5YmQzODBhMTEiLCJhdWQiOlsiYTBlZWJjOTktOWMwYi00ZWY4LWJiNmQtNmJiOWJkMzgwYTExIiwiMTI3LjAuMC4xIl0sImV4cCI6MTcyMzkyMjEyNiwiTGFzdElQIjoiMTI3LjAuMC4xIiwiQ3JlYXRlZCI6MTcyMzkxODUyNjEyNn0.eDR5mwPC5OaUBw4ACBJ3B7Zb0-rLKvzYG_USwW-m3W2jPT-y8XvwaAq-TvO4xli6g6QMljnFPjbSpFrzajw46w",
			Refresh: "refreshtoken",
			MockBehaviour: func(ms *mock_service.MockemailService, m *mock_service.Mockstorage) {
				refresh, _ := bcrypt.GenerateFromPassword([]byte("refreshtoken"), bcrypt.DefaultCost)
				m.EXPECT().GetFreeRefreshToken(gomock.Any(), "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", uint64(1723918526126)).Return("rowid", refresh, nil)
				m.EXPECT().RevokeRefreshToken(gomock.Any(), "rowid").Return(middleware.NotFoundError("row not found"))
			},
			Err: middleware.NotFoundError("row not found"),
		},
		{
			Name:    "different ip message sending",
			Secret:  "secret",
			Token:   "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJqdGkiOiJhMGVlYmM5OS05YzBiLTRlZjgtYmI2ZC02YmI5YmQzODBhMTEiLCJhdWQiOlsiYTBlZWJjOTktOWMwYi00ZWY4LWJiNmQtNmJiOWJkMzgwYTExIiwiMTI3LjAuMC4xIl0sImV4cCI6MTcyMzkyMjEyNiwiTGFzdElQIjoiMTI3LjAuMC4xIiwiQ3JlYXRlZCI6MTcyMzkxODUyNjEyNn0.eDR5mwPC5OaUBw4ACBJ3B7Zb0-rLKvzYG_USwW-m3W2jPT-y8XvwaAq-TvO4xli6g6QMljnFPjbSpFrzajw46w",
			Refresh: "refreshtoken",
			Ip:      "128.128.128.128",
			MockBehaviour: func(ms *mock_service.MockemailService, m *mock_service.Mockstorage) {
				refresh, _ := bcrypt.GenerateFromPassword([]byte("refreshtoken"), bcrypt.DefaultCost)
				m.EXPECT().GetFreeRefreshToken(gomock.Any(), "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", uint64(1723918526126)).Return("rowid", refresh, nil)
				m.EXPECT().RevokeRefreshToken(gomock.Any(), "rowid").Return(nil)
				ms.EXPECT().SendEmailWarning("128.128.128.128")
			},
			ExceptedID: "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		},
		{
			Name:    "right ip verification",
			Secret:  "secret",
			Token:   "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJqdGkiOiJhMGVlYmM5OS05YzBiLTRlZjgtYmI2ZC02YmI5YmQzODBhMTEiLCJhdWQiOlsiYTBlZWJjOTktOWMwYi00ZWY4LWJiNmQtNmJiOWJkMzgwYTExIiwiMTI3LjAuMC4xIl0sImV4cCI6MTcyMzkyMjEyNiwiTGFzdElQIjoiMTI3LjAuMC4xIiwiQ3JlYXRlZCI6MTcyMzkxODUyNjEyNn0.eDR5mwPC5OaUBw4ACBJ3B7Zb0-rLKvzYG_USwW-m3W2jPT-y8XvwaAq-TvO4xli6g6QMljnFPjbSpFrzajw46w",
			Refresh: "refreshtoken",
			Ip:      "127.0.0.1",
			MockBehaviour: func(ms *mock_service.MockemailService, m *mock_service.Mockstorage) {
				refresh, _ := bcrypt.GenerateFromPassword([]byte("refreshtoken"), bcrypt.DefaultCost)
				m.EXPECT().GetFreeRefreshToken(gomock.Any(), "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", uint64(1723918526126)).Return("rowid", refresh, nil)
				m.EXPECT().RevokeRefreshToken(gomock.Any(), "rowid").Return(nil)
			},
			ExceptedID: "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		},
	}

	for _, v := range table {
		t.Run(v.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			logger := mock_service.NewMocklogger(ctrl)
			logger.EXPECT().Info(gomock.Any()).AnyTimes()
			logger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
			logger.EXPECT().Warnf(gomock.Any(), gomock.Any()).AnyTimes()
			storage := mock_service.NewMockstorage(ctrl)
			emailService := mock_service.NewMockemailService(ctrl)

			if v.MockBehaviour != nil {
				v.MockBehaviour(emailService, storage)
			}
			service, err := New(logger, storage, emailService, v.Secret)
			if v.Err != nil && len(v.Secret) == 0 {
				assert.Error(t, err)
				assert.EqualError(t, err, v.Err.Error())
				return
			} else {
				assert.NoError(t, err)
			}

			claims, err := service.ValidateUserToken(context.Background(), v.Token, v.Refresh, v.Ip)
			if v.Err != nil && assert.Error(t, err) {
				assert.EqualError(t, err, v.Err.Error())
			} else if v.Err == nil && assert.NoError(t, err) {
				assert.Equal(t, v.ExceptedID, claims.ID)
			}
		})
	}
}
