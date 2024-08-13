package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator"
	"github.com/reversersed/AuthService/pkg/middleware"
)

type ValidationErrors validator.ValidationErrors
type Validator struct {
	*validator.Validate
}

func New() *Validator {
	v := validator.New()

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return fld.Name
		}
		return name
	})
	return &Validator{v}
}
func (v *Validator) StructValidation(data any) error {
	result := v.Validate.Struct(data)

	if result == nil {
		return nil
	}
	if er, ok := result.(*validator.InvalidValidationError); ok {
		return middleware.InternalError(er.Error())
	}
	for _, i := range result.(validator.ValidationErrors) {
		return middleware.BadRequestError(errorToStringByTag(i))
	}
	return nil
}
func errorToStringByTag(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s: field is required", err.Field())
	default:
		return err.Tag()
	}
}
