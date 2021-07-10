package responsex

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Response represents RESTful response
type Response struct {
	HttpStatus int         `json:"-"`
	Status     int         `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
}

// New create new response
func New(status int, msg string, data interface{}) *Response {
	return &Response{
		HttpStatus: http.StatusOK,
		Status:     status,
		Message:    msg,
		Data:       data,
	}
}

// SetHttpStatus set http response status
func (t *Response) SetHttpStatus(httpStatus int) *Response {
	t.HttpStatus = httpStatus
	return t
}

// SetStatus set response body status code
func (t *Response) SetStatus(status int) *Response {
	t.Status = status
	return t
}

// SetMsg set response message
func (t *Response) SetMsg(msg string) *Response {
	t.Message = msg
	return t
}

// SetData set response data
func (t *Response) SetData(data interface{}) *Response {
	t.Data = data
	return t
}

// R response code
func R(c echo.Context, response *Response) error {
	return c.JSON(response.HttpStatus, response)
}

// Data shotcut response with data
func Data(c echo.Context, data interface{}) error {
	response := New(0, "OK", data)
	return c.JSON(response.HttpStatus, response)
}
