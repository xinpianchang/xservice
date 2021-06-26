package responsex

import (
	"bytes"
	"fmt"
)

type Error struct {
	Status     int
	HttpStatus int
	Message    string
	Internal   error
}

func NewError(status int, message string) *Error {
	return &Error{Status: status, Message: message}
}

func (e *Error) Error() string {
	var buf bytes.Buffer
	_, _ = fmt.Fprint(&buf, "status:", e.Status)
	if e.HttpStatus != 0 {
		_, _ = fmt.Fprint(&buf, ", httpStatus:", e.HttpStatus)
	}
	if e.Message != "" {
		_, _ = fmt.Fprint(&buf, ", message:", e.Message)
	}
	if e.Internal != nil {
		_, _ = fmt.Fprint(&buf, ", internal:", e.Internal)
	}
	return buf.String()
}

func (e *Error) SetInternal(internal error) *Error {
	e.Internal = internal
	return e
}

func (e *Error) SetHttpStatus(httpStatus int) *Error {
	e.HttpStatus = httpStatus
	return e
}
