package validator

import (
	"fmt"
	"testing"

	"github.com/reversersed/AuthService/pkg/middleware"
	"github.com/stretchr/testify/assert"
)

func TestValidator(t *testing.T) {
	type testStruct struct {
		Required   string `validate:"required"`
		OtherField string `validate:"ip"`
		Uuid       string `validate:"uuid"`
	}
	var cases = []struct {
		name  string
		field *testStruct
		err   string
	}{
		{"validated struct", &testStruct{Required: "1", OtherField: "127.0.0.1"}, ""},
		{"required field testing", &testStruct{}, fmt.Sprintf("%v: Required: field is required", middleware.ErrBadRequest)},
		{"unknown field", &testStruct{Required: "1", OtherField: "1"}, fmt.Sprintf("%v: ip", middleware.ErrBadRequest)},
		{"uuid error", &testStruct{Required: "1", OtherField: "127.0.0.1", Uuid: "123"}, fmt.Sprintf("%v: Uuid: field must be a valid uuid", middleware.ErrBadRequest)},
		{"valid uuid", &testStruct{Required: "1", OtherField: "127.0.0.1", Uuid: "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"}, ""},
		{"nil struct", nil, "something wrong happened: validator: (nil *validator.testStruct)"},
	}

	valid := New()

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T) {
			err := valid.StructValidation(v.field)

			if len(v.err) == 0 {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, v.err)
			}
		})
	}
}
