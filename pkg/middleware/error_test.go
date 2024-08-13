package middleware

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnwrapping(t *testing.T) {
	table := []struct {
		Name         string
		Error        *customError
		WrappedError error
	}{
		{"not found error", NotFoundError(""), ErrNotFound},
		{"internal error", InternalError(""), ErrInternal},
		{"bad request error", BadRequestError(""), ErrBadRequest},
		{"confict error", ConfictError(""), ErrConflict},
	}

	for _, c := range table {
		t.Run(c.Name, func(t *testing.T) {
			assert.Equal(t, c.Error.Unwrap(), c.WrappedError)
		})
	}
}
func TestErrors(t *testing.T) {
	table := []struct {
		Name      string
		ErrorFunc func(string, ...any) *customError
		Message   string
		Args      []any
	}{
		{"not found error with args", NotFoundError, "%s was not found as %d", []any{"test", 5}},
		{"not found error without args", NotFoundError, "message contains no args", []any{}},
		{"internal error", InternalError, "something wrong happened with %v", []any{"TokenHandler"}},
		{"bad request error", BadRequestError, "token: %f is not a string, %d is not a token", []any{5.1, 2}},
		{"confict error", ConfictError, "user %v already exists", []any{struct {
			name string
			id   int
		}{name: "user", id: 231}}},
	}

	for _, c := range table {
		t.Run(c.Name, func(t *testing.T) {
			err := c.ErrorFunc(c.Message, c.Args...)
			assert.EqualError(t, err, fmt.Sprintf("%v: %s", err.Unwrap(), fmt.Sprintf(c.Message, c.Args...)))
		})
	}
}
func TestIsCustomError(t *testing.T) {
	noError := errors.New("base error")
	Error := InternalError("internal message")

	assert.True(t, IsCustomError(Error))
	assert.False(t, IsCustomError(noError))
}
