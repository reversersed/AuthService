package endpoint

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mock_endpoint "github.com/reversersed/AuthService/internal/endpoint/mocks"
	Service "github.com/reversersed/AuthService/internal/service"
	"github.com/reversersed/AuthService/pkg/middleware"
	"github.com/stretchr/testify/assert"
)

func TestGetAccessToken(t *testing.T) {
	table := []struct {
		Name           string
		MockBehaviour  func(*mock_endpoint.Mockservice, *mock_endpoint.Mocklogger, *mock_endpoint.Mockvalidator)
		ExceptedStatus int
		IpHeader       string
		QueryId        string
		ExceptedBody   string
	}{
		{
			Name: "empty request",
			MockBehaviour: func(m1 *mock_endpoint.Mockservice, m2 *mock_endpoint.Mocklogger, m3 *mock_endpoint.Mockvalidator) {
				m2.EXPECT().Info(gomock.Any()).AnyTimes()
			},
			ExceptedStatus: http.StatusBadRequest,
		},
		{
			Name: "wrong id type",
			MockBehaviour: func(m1 *mock_endpoint.Mockservice, m2 *mock_endpoint.Mocklogger, m3 *mock_endpoint.Mockvalidator) {
				m2.EXPECT().Info(gomock.Any()).AnyTimes()
				m3.EXPECT().StructValidation(gomock.Any()).Return(middleware.BadRequestError("wrong id request"))
			},
			ExceptedStatus: http.StatusBadRequest,
			QueryId:        "123",
		},
		{
			Name: "error generating access token",
			MockBehaviour: func(m1 *mock_endpoint.Mockservice, m2 *mock_endpoint.Mocklogger, m3 *mock_endpoint.Mockvalidator) {
				m2.EXPECT().Info(gomock.Any()).AnyTimes()
				m3.EXPECT().StructValidation(gomock.Any()).Return(nil)
				m1.EXPECT().GenerateAccessToken(gomock.Any(), "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", "127.0.0.1").Return("", "", middleware.NotFoundError("id not found"))
			},
			ExceptedStatus: http.StatusNotFound,
			QueryId:        "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
			IpHeader:       "127.0.0.1",
		},
		{
			Name: "success response",
			MockBehaviour: func(m1 *mock_endpoint.Mockservice, m2 *mock_endpoint.Mocklogger, m3 *mock_endpoint.Mockvalidator) {
				m2.EXPECT().Info(gomock.Any()).AnyTimes()
				m3.EXPECT().StructValidation(gomock.Any()).Return(nil)
				m1.EXPECT().GenerateAccessToken(gomock.Any(), "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", "127.0.0.1").Return("token", "refresh", nil)
			},
			QueryId:        "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
			IpHeader:       "127.0.0.1",
			ExceptedStatus: http.StatusOK,
			ExceptedBody:   "{\"token\":\"token\",\"refresh\":\"refresh\"}",
		},
	}

	for _, v := range table {
		t.Run(v.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			logger := mock_endpoint.NewMocklogger(ctrl)
			service := mock_endpoint.NewMockservice(ctrl)
			validator := mock_endpoint.NewMockvalidator(ctrl)

			h := New(service, logger, validator)
			v.MockBehaviour(service, logger, validator)

			e := gin.Default()
			e.Use(middleware.ErrorHandler)
			h.RegisterRoute(e.Group("/api/"))

			w := httptest.NewRecorder()
			var r *http.Request
			if len(v.QueryId) == 0 {
				r = httptest.NewRequest(http.MethodPost, "/api/v1/token", nil)
			} else {
				body, _ := json.Marshal(&GetTokenRequest{Guid: v.QueryId})
				r = httptest.NewRequest(http.MethodPost, "/api/v1/token", bytes.NewBuffer(body))
			}
			r.Header.Add("X-FORWARDED-FOR", v.IpHeader)
			e.ServeHTTP(w, r)

			assert.Equal(t, v.ExceptedStatus, w.Result().StatusCode)
			if len(v.ExceptedBody) != 0 {
				b, _ := io.ReadAll(w.Result().Body)
				assert.Equal(t, v.ExceptedBody, string(b))
			}
			err := h.Close()
			assert.NoError(t, err)
		})
	}
}
func TestRefreshToken(t *testing.T) {
	table := []struct {
		Name           string
		MockBehaviour  func(*mock_endpoint.Mockservice, *mock_endpoint.Mocklogger, *mock_endpoint.Mockvalidator)
		ExceptedStatus int
		IpHeader       string
		QueryToken     string
		QueryRefresh   string
		ExceptedBody   string
	}{
		{
			Name: "empty request",
			MockBehaviour: func(m1 *mock_endpoint.Mockservice, m2 *mock_endpoint.Mocklogger, m3 *mock_endpoint.Mockvalidator) {
				m2.EXPECT().Info(gomock.Any()).AnyTimes()
			},
			ExceptedStatus: http.StatusBadRequest,
		},
		{
			Name: "validation error",
			MockBehaviour: func(m1 *mock_endpoint.Mockservice, m2 *mock_endpoint.Mocklogger, m3 *mock_endpoint.Mockvalidator) {
				m2.EXPECT().Info(gomock.Any()).AnyTimes()
				m3.EXPECT().StructValidation(gomock.Any()).Return(middleware.BadRequestError("wrong token format"))
			},
			ExceptedStatus: http.StatusBadRequest,
			QueryToken:     "123",
			QueryRefresh:   "321",
		},
		{
			Name: "token validation error",
			MockBehaviour: func(m1 *mock_endpoint.Mockservice, m2 *mock_endpoint.Mocklogger, m3 *mock_endpoint.Mockvalidator) {
				m2.EXPECT().Info(gomock.Any()).AnyTimes()
				m3.EXPECT().StructValidation(gomock.Any()).Return(nil)
				m1.EXPECT().ValidateUserToken(gomock.Any(), "token", "refresh", "127.0.0.1").Return(nil, middleware.BadRequestError("cant validate user token"))
			},
			ExceptedStatus: http.StatusBadRequest,
			QueryToken:     "token",
			QueryRefresh:   "refresh",
			IpHeader:       "127.0.0.1",
		},
		{
			Name: "token regenerating error",
			MockBehaviour: func(m1 *mock_endpoint.Mockservice, m2 *mock_endpoint.Mocklogger, m3 *mock_endpoint.Mockvalidator) {
				m2.EXPECT().Info(gomock.Any()).AnyTimes()
				m3.EXPECT().StructValidation(gomock.Any()).Return(nil)
				m1.EXPECT().ValidateUserToken(gomock.Any(), "token", "refresh", "127.0.0.1").Return(&Service.Claims{LastIP: "127.0.0.1"}, nil)
				m1.EXPECT().GenerateAccessToken(gomock.Any(), "", "127.0.0.1").Return("", "", middleware.BadRequestError("wrong id provided"))
			},
			ExceptedStatus: http.StatusBadRequest,
			QueryToken:     "token",
			QueryRefresh:   "refresh",
			IpHeader:       "127.0.0.1",
		},
		{
			Name: "success update",
			MockBehaviour: func(m1 *mock_endpoint.Mockservice, m2 *mock_endpoint.Mocklogger, m3 *mock_endpoint.Mockvalidator) {
				m2.EXPECT().Info(gomock.Any()).AnyTimes()
				m3.EXPECT().StructValidation(gomock.Any()).Return(nil)
				m1.EXPECT().ValidateUserToken(gomock.Any(), "token", "refresh", "127.0.0.1").Return(&Service.Claims{LastIP: "127.0.0.1"}, nil)
				m1.EXPECT().GenerateAccessToken(gomock.Any(), "", "127.0.0.1").Return("token", "refresh", nil)
			},
			ExceptedStatus: http.StatusOK,
			QueryToken:     "token",
			QueryRefresh:   "refresh",
			IpHeader:       "127.0.0.1",
			ExceptedBody:   "{\"token\":\"token\",\"refresh\":\"refresh\"}",
		},
	}

	for _, v := range table {
		t.Run(v.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			logger := mock_endpoint.NewMocklogger(ctrl)
			service := mock_endpoint.NewMockservice(ctrl)
			validator := mock_endpoint.NewMockvalidator(ctrl)

			h := New(service, logger, validator)
			v.MockBehaviour(service, logger, validator)

			e := gin.Default()
			e.Use(middleware.ErrorHandler)
			h.RegisterRoute(e.Group("/api/"))

			w := httptest.NewRecorder()
			var r *http.Request
			if len(v.QueryToken) == 0 {
				r = httptest.NewRequest(http.MethodPost, "/api/v1/token/refresh", nil)
			} else {
				body, _ := json.Marshal(&RefreshTokenRequest{Token: v.QueryToken, Refresh: v.QueryRefresh})
				r = httptest.NewRequest(http.MethodPost, "/api/v1/token/refresh", bytes.NewBuffer(body))
			}
			r.Header.Add("X-FORWARDED-FOR", v.IpHeader)
			e.ServeHTTP(w, r)

			assert.Equal(t, v.ExceptedStatus, w.Result().StatusCode)
			if len(v.ExceptedBody) != 0 {
				b, _ := io.ReadAll(w.Result().Body)
				assert.Equal(t, v.ExceptedBody, string(b))
			}
			err := h.Close()
			assert.NoError(t, err)
		})
	}
}
