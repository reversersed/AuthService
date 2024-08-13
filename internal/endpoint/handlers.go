package endpoint

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/reversersed/AuthService/pkg/middleware"
)

// @Summary Get access token with provided GUID
// @Tags token
// @Produce json
// @Param body body GetTokenRequest true "Request with user's guid"
// @Success 200 {object} GetTokenResponse
// @Failure 400 {object} middleware.customError "Received bad request"
// @Failure 500 {object} middleware.customError "Internal error occured"
// @Router /v1/token [post]
func (e *endpoint) GetAccessToken(c *gin.Context) {
	var request GetTokenRequest
	if err := c.BindJSON(&request); err != nil {
		c.Error(middleware.BadRequestError(err.Error()))
		return
	}

	if err := e.validator.StructValidation(&request); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, GetTokenResponse{Token: request.Guid, Refresh: "321"})
}
