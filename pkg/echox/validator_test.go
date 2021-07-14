package echox

import (
	"testing"
	"time"

	"github.com/go-playground/validator"
)

// test deep validate
func TestValidator_validate(t *testing.T) {
	type Inner struct {
		InnerName string `validate:"required"`
	}

	type inner struct {
		InnerName string `validate:"required"`
	}

	tests := []struct {
		name    string
		obj     interface{}
		wantErr bool
	}{
		{
			name: "simple required 1",
			obj: &struct {
				Name string `validate:"required"`
			}{Name: "name"},
			wantErr: false,
		},

		{
			name: "simple required 2",
			obj: &struct {
				Name string `validate:"required"`
			}{},
			wantErr: true,
		},

		{
			name: "Inner required 1",
			obj: &struct {
				Name string `validate:"required"`
				Inner
			}{"name", Inner{InnerName: "Inner"}},
			wantErr: false,
		},

		{
			name: "Inner required 2",
			obj: &struct {
				Name string `validate:"required"`
				Inner
			}{"", Inner{InnerName: "Inner"}},
			wantErr: true,
		},

		{
			name: "Inner required 3",
			obj: &struct {
				Inner
			}{},
			wantErr: true,
		},

		{
			name: "Inner required 4",
			obj: &struct {
				Inner Inner
			}{},
			wantErr: true,
		},

		{
			name: "Inner required 5",
			obj: &struct {
				Inner
			}{Inner{InnerName: "ok"}},
			wantErr: false,
		},

		{
			name: "Inner required 6",
			obj: &struct {
				Inner Inner
			}{Inner: Inner{InnerName: "ok"}},
			wantErr: false,
		},

		{
			name: "unexported inner",
			obj: &struct {
				inner Inner
			}{inner: Inner{InnerName: "ok"}},
			wantErr: false,
		},

		{
			name: "unexported inner 2",
			obj: &struct {
				inner
			}{inner: inner{InnerName: ""}},
			wantErr: true,
		},

		{
			name: "point innter 1",
			obj: &struct {
				Inner *Inner `json:"inner" validate:"required"`
			}{},
			wantErr: true,
		},

		{
			name: "point innter 2",
			obj: &struct {
				Inner *Inner `json:"inner" validate:"required"`
			}{Inner: &Inner{InnerName: "OK"}},
			wantErr: false,
		},

		{
			name: "time struct 1",
			obj: struct {
				Date time.Time `validate:"required"`
			}{Date: time.Now()},
			wantErr: false,
		},

		{
			name: "time struct 2",
			obj: struct {
				Date time.Time `validate:"required"`
			}{Date: time.Time{}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		v := &EchoValidator{Validator: validator.New()}
		if err := v.Validate(tt.obj); (err != nil) != tt.wantErr {
			t.Errorf("%q. Validator.validate() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}
