package echox

import (
	"reflect"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

// ConfigValidator enable echo use ext validation framework
// ref: https://github.com/go-playground/validator
func ConfigValidator(e *echo.Echo) {
	v := &EchoValidator{Validator: validator.New(), binder: e.Binder}
	e.Validator = v
	e.Binder = v
}

type EchoValidator struct {
	Validator *validator.Validate
	binder    echo.Binder
}

func (t *EchoValidator) Validate(i interface{}) error {
	return t.Validator.Struct(i)
}

func (t *EchoValidator) Bind(i interface{}, c echo.Context) error {
	err := t.binder.Bind(i, c)
	if err != nil {
		return err
	}
	tt := reflect.TypeOf(i)
	for {
		if tt.Kind() == reflect.Ptr {
			tt = tt.Elem()
			continue
		}

		break
	}

	if tt.Kind() == reflect.Struct {
		return t.Validate(i)
	}

	return nil
}
