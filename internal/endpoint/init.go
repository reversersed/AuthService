package endpoint

import (
	"context"

	"github.com/gin-gonic/gin"
	Service "github.com/reversersed/AuthService/internal/service"
)

//go:generate mockgen -source=init.go -destination=mocks/mock.go

type service interface {
	GenerateAccessToken(context.Context, string, string) (string, string, error)
	ValidateUserToken(context.Context, string, string, string) (*Service.Claims, error)
}
type logger interface {
	Info(...any)
}
type validator interface {
	StructValidation(data any) error
}
type endpoint struct {
	logger    logger
	service   service
	validator validator
}

func New(service service, logger logger, validator validator) *endpoint {
	return &endpoint{
		service:   service,
		logger:    logger,
		validator: validator,
	}
}
func (*endpoint) Close() error {
	return nil
}
func (e *endpoint) RegisterRoute(r *gin.RouterGroup) {
	g := r.Group("/v1/token")
	{
		g.POST("", e.GetAccessToken)
		g.POST("/refresh", e.RefreshToken)
	}
	e.logger.Info("endpoint routes registered")
}
