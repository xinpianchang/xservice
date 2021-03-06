// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: buf/v1/calculator.proto

package v1

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
)

// Validate checks the field values on AddIntRequest with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *AddIntRequest) Validate() error {
	if m == nil {
		return nil
	}

	if val := m.GetA(); val < 0 || val >= 100 {
		return AddIntRequestValidationError{
			field:  "A",
			reason: "value must be inside range [0, 100)",
		}
	}

	// no validation rules for B

	return nil
}

// AddIntRequestValidationError is the validation error returned by
// AddIntRequest.Validate if the designated constraints aren't met.
type AddIntRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e AddIntRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e AddIntRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e AddIntRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e AddIntRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e AddIntRequestValidationError) ErrorName() string { return "AddIntRequestValidationError" }

// Error satisfies the builtin error interface
func (e AddIntRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sAddIntRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = AddIntRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = AddIntRequestValidationError{}

// Validate checks the field values on AddIntResponse with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *AddIntResponse) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for Result

	return nil
}

// AddIntResponseValidationError is the validation error returned by
// AddIntResponse.Validate if the designated constraints aren't met.
type AddIntResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e AddIntResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e AddIntResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e AddIntResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e AddIntResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e AddIntResponseValidationError) ErrorName() string { return "AddIntResponseValidationError" }

// Error satisfies the builtin error interface
func (e AddIntResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sAddIntResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = AddIntResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = AddIntResponseValidationError{}
