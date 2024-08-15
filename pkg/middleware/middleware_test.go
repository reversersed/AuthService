package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMiddleware(t *testing.T) {
	table := []struct {
		Name            string
		Endpoint        func(*gin.Context)
		ExceptedCode    int
		ExceptedMessage string
	}{
		{"undefined error", func(ctx *gin.Context) { ctx.Error(errors.New("smth")) }, http.StatusInternalServerError, "smth"},
		{"not found 404 code", func(ctx *gin.Context) { ctx.Error(NotFoundError("not found")) }, http.StatusNotFound, "not found"},
	}

	for _, v := range table {
		t.Run(v.Name, func(t *testing.T) {
			gin.SetMode(gin.ReleaseMode)
			engine := gin.Default()
			engine.Use(ErrorHandler)
			engine.GET("/", v.Endpoint)

			r := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			engine.ServeHTTP(w, r)
			assert.Equal(t, v.ExceptedCode, w.Result().StatusCode)
			err := new(customError)
			json.NewDecoder(w.Result().Body).Decode(err)
			assert.Equal(t, err.Message, v.ExceptedMessage)
		})
	}
}
