package echox

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

// ConfigValidator enable echo use ext validation framework
// ref: https://github.com/go-playground/validator
func ConfigValidator(e *echo.Echo) {
	v := &echoValidator{validator: validator.New(), binder: e.Binder}
	e.Validator = v
	e.Binder = v
}

type echoValidator struct {
	validator *validator.Validate
	binder    echo.Binder
}

func (t *echoValidator) Validate(i interface{}) error {
	return t.validator.Struct(i)
}

func (t *echoValidator) Bind(i interface{}, c echo.Context) error {
	err := t.binder.Bind(i, c)
	if err != nil {
		return err
	}
	return t.Validate(i)
}
