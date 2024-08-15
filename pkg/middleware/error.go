package middleware

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrNotFound   = errors.New("not found")
	ErrInternal   = errors.New("something wrong happened")
	ErrBadRequest = errors.New("received bad request")
	ErrConflict   = errors.New("conflict occured")
)

type customError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
	Err        error  `json:"-"`
}

func (c customError) Unwrap() error {
	return c.Err
}
func (c customError) Error() string {
	return fmt.Sprintf("%v: %s", c.Err, c.Message)
}
func IsCustomError(err error) bool {
	e := new(customError)
	return errors.As(err, &e)
}

func NotFoundError(message string, args ...any) *customError {
	return &customError{
		Message:    fmt.Sprintf(message, args...),
		StatusCode: http.StatusNotFound,
		Err:        ErrNotFound,
	}
}
func InternalError(message string, args ...any) *customError {
	return &customError{
		Message:    fmt.Sprintf(message, args...),
		StatusCode: http.StatusNotFound,
		Err:        ErrInternal,
	}
}
func BadRequestError(message string, args ...any) *customError {
	return &customError{
		Message:    fmt.Sprintf(message, args...),
		StatusCode: http.StatusBadRequest,
		Err:        ErrBadRequest,
	}
}
func ConfictError(message string, args ...any) *customError {
	return &customError{
		Message:    fmt.Sprintf(message, args...),
		StatusCode: http.StatusConflict,
		Err:        ErrConflict,
	}
}
