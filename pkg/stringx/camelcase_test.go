package stringx

import (
	"testing"
)

func TestCamelCase(t *testing.T) {
	tests := []struct {
		args string
		want string
	}{
		{"", ""},
		{"name", "Name"},
		{"name1", "Name1"},
		{"name_if", "NameIf"},
		{"_name_if", "XNameIf"},
	}
	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			if got := CamelCase(tt.args); got != tt.want {
				t.Errorf("CamelCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLowerCamelCase(t *testing.T) {
	tests := []struct {
		args string
		want string
	}{
		{"", ""},
		{"name", "name"},
		{"name1", "name1"},
		{"name_if", "nameIf"},
		{"_name_if", "xNameIf"},
	}
	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			if got := LowerCamelCase(tt.args); got != tt.want {
				t.Errorf("CamelCase() = %v, want %v", got, tt.want)
			}
		})
	}
}
