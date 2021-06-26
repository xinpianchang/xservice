package responsex

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Response struct {
	HttpStatus int         `json:"-"`
	Status     int         `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
}

func New(status int, msg string, data interface{}) *Response {
	return &Response{
		HttpStatus: http.StatusOK,
		Status:     status,
		Message:    msg,
		Data:       data,
	}
}

func (t *Response) SetHttpStatus(httpStatus int) *Response {
	t.HttpStatus = httpStatus
	return t
}

func (t *Response) SetStatus(status int) *Response {
	t.Status = status
	return t
}

func (t *Response) SetMsg(msg string) *Response {
	t.Message = msg
	return t
}

func (t *Response) SetData(data interface{}) *Response {
	t.Data = data
	return t
}

// R response code
func R(c echo.Context, response *Response) error {
	return c.JSON(response.HttpStatus, response)
}

func Data(c echo.Context, data interface{}) error {
	response := New(0, "OK", data)
	return c.JSON(response.HttpStatus, response)
}
