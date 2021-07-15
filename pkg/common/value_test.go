package common

import (
	"reflect"
	"testing"
)

func TestIsEmptyValue(t *testing.T) {
	type args struct {
		v reflect.Value
	}
	type Foo struct{}

	var bar1 interface{}
	var m1 = make(map[string]interface{})
	var m2 = map[string]interface{}{"a": 1}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{"empty int", args{reflect.ValueOf(int(0))}, true},
		{"none empty int", args{reflect.ValueOf(int(1))}, false},
		{"uint 1", args{reflect.ValueOf(uint(0))}, true},
		{"uint 2", args{reflect.ValueOf(uint(1))}, false},
		{"empty string", args{reflect.ValueOf("")}, true},
		{"none empty string", args{reflect.ValueOf("abc")}, false},
		{"empty struct", args{reflect.ValueOf(Foo{})}, false},
		{"empty struct ptr", args{reflect.ValueOf(&Foo{})}, false},
		{"empty map1", args{reflect.ValueOf(m1)}, true},
		{"empty map2", args{reflect.ValueOf(m2)}, false},
		{"bool 1", args{reflect.ValueOf(true)}, false},
		{"bool 2", args{reflect.ValueOf(false)}, true},
		{"float 1", args{reflect.ValueOf(float32(1))}, false},
		{"float 1", args{reflect.ValueOf(float32(0))}, true},
		{"empty nil", args{reflect.ValueOf(nil)}, false},
		{"empty interface{}", args{reflect.ValueOf(bar1)}, false},
		{"empty interface{} ptr", args{reflect.ValueOf(&bar1)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEmptyValue(tt.args.v); got != tt.want {
				t.Errorf("IsEmptyValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
