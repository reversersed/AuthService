package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler(c *gin.Context) {
	c.Next()

	if e := c.Errors.Last(); e != nil {
		if IsCustomError(e.Err) {
			err := e.Err.(*customError)
			c.JSON(err.StatusCode, err)
		} else {
			c.JSON(http.StatusInternalServerError, InternalError(e.Error()))
		}
	}
}
