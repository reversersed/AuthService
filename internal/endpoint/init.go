package endpoint

import "github.com/gin-gonic/gin"

type service interface {
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
		g.POST("/refresh")
	}
	e.logger.Info("endpoint routes registered")
}
