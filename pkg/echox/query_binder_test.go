package echox

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_query_binder_query_param(t *testing.T) {
	type Foo struct {
		Id   int    `json:"id"`
		Name string `json:"name" query:"name"`
	}

	build := func() (*echo.Echo, *httptest.ResponseRecorder, echo.Context) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/api/endpoint?name=foo", strings.NewReader(`{"id": 1}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		return e, rec, c
	}

	handler := func(c echo.Context) error {
		var foo Foo
		if err := c.Bind(&foo); err != nil {
			return err
		}
		return c.JSON(http.StatusOK, foo)
	}

	_, rec, c := build()
	err := handler(c)
	require.NoError(t, err)
	body := rec.Body.String()
	assert.Contains(t, body, `"id":1`)
	assert.Contains(t, body, `"name":""`)

	ex, rec, c := build()
	UseDefaultQueryBinder(ex)
	err = handler(c)
	require.NoError(t, err)
	body = rec.Body.String()
	assert.Contains(t, body, `"id":1`)
	assert.Contains(t, body, `"name":"foo"`)
}
