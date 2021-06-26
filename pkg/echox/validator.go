package echox

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

// ConfigValidator enable echo use ext validation framework
// ref: https://github.com/go-playground/validator
func ConfigValidator(e *echo.Echo) {
	e.Validator = &echoValidator{validator: validator.New()}
}

type echoValidator struct {
	validator *validator.Validate
}

func (t *echoValidator) Validate(i interface{}) error {
	return t.validator.Struct(i)
}
