package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
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
	v.RegisterValidation("uuid", UuidValidation)
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
	case "uuid":
		return fmt.Sprintf("%s: field must be a valid uuid", err.Field())
	default:
		return err.Tag()
	}
}
func UuidValidation(fl validator.FieldLevel) bool {
	if len(fl.Field().String()) == 0 {
		return true
	}
	_, err := uuid.Parse(fl.Field().String())
	return err == nil
}
