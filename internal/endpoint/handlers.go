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
	var ip string
	if ip = c.ClientIP(); len(ip) == 0 {
		ip = c.RemoteIP()
	}

	token, refresh, err := e.service.GenerateAccessToken(request.Guid, ip)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, GetTokenResponse{Token: token, Refresh: refresh})
}

// @Summary Updates token with provided refresh token
// @Tags token
// @Produce json
// @Param body body RefreshTokenRequest true "Refreshed pair token-refresh"
// @Success 200 {object} RefreshTokenResponse
// @Failure 400 {object} middleware.customError "Received bad request"
// @Failure 404 {object} middleware.customError "Refresh token was incorrect or token not found"
// @Failure 500 {object} middleware.customError "Internal error occured"
// @Router /v1/token/refresh [post]
func (e *endpoint) RefreshToken(c *gin.Context) {
	var request RefreshTokenRequest
	if err := c.BindJSON(&request); err != nil {
		c.Error(middleware.BadRequestError(err.Error()))
		return
	}

	if err := e.validator.StructValidation(&request); err != nil {
		c.Error(err)
		return
	}
	var ip string
	if ip = c.ClientIP(); len(ip) == 0 {
		ip = c.RemoteIP()
	}
	claims, err := e.service.ValidateUserToken(request.Token, request.Refresh, ip)
	if err != nil {
		c.Error(err)
		return
	}
	token, refresh, err := e.service.GenerateAccessToken(claims.ID, claims.LastIP)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, RefreshTokenResponse{Token: token, Refresh: refresh})
}
