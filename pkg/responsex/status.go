package responsex

import (
	"fmt"
)

type StatusCode int

var statusMap = make(map[StatusCode]string)

func Status(status StatusCode, message ...string) *Error {
	var m string
	if message != nil {
		item := make([]interface{}, 0, len(message))
		for _, it := range message {
			item = append(item, it)
		}
		m = fmt.Sprint(item...)
	} else {
		if v, ok := statusMap[status]; ok {
			m = v
		} else {
			m = fmt.Sprintf("status:%d", status)
		}
	}
	return NewError(int(status), m)
}

func StatusMessage(status StatusCode, args ...interface{}) *Error {
	var m string
	if v, ok := statusMap[status]; ok {
		m = v
	}
	if len(args) > 0 {
		m = fmt.Sprintf(m, args...)
	}
	return NewError(int(status), m)
}

func SetStatusMap(m map[StatusCode]string) {
	statusMap = m
}
