package echox

import "github.com/labstack/echo/v4"

// UseDefaultQueryBinder change behavior of query binder to use default binder
// enable always bind query params
// cause https://github.com/labstack/echo/issues/1670 disable bind query params for POST/PUT methods
func UseDefaultQueryBinder(e *echo.Echo) {
	e.Binder = &defaultQueryBinder{
		binder:        e.Binder,
		defaultBinder: new(echo.DefaultBinder),
	}
}

type defaultQueryBinder struct {
	binder        echo.Binder
	defaultBinder *echo.DefaultBinder
}

func (b *defaultQueryBinder) Bind(i interface{}, c echo.Context) error {
	if err := b.defaultBinder.BindQueryParams(c, i); err != nil {
		return err
	}
	return b.binder.Bind(i, c)
}
