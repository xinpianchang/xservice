package stringx

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	type args struct {
		list []string
		str  string
	}
	tests := []struct {
		args args
		want bool
	}{
		{
			args{[]string{"a", "b"}, "a"},
			true,
		},
		{
			args{[]string{"a", "b"}, "c"},
			false,
		},
		{
			args{[]string{"a", "b"}, "aa"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(path.Join(tt.args.list...), func(t *testing.T) {
			if got := Contains(tt.args.list, tt.args.str); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainsIgnoreCase(t *testing.T) {
	type args struct {
		list []string
		str  string
	}
	tests := []struct {
		args args
		want bool
	}{
		{
			args{[]string{"a", "b"}, "a"},
			true,
		},
		{
			args{[]string{"A", "a"}, "a"},
			true,
		},
		{
			args{[]string{"a", "b"}, "c"},
			false,
		},
		{
			args{[]string{"a", "b"}, "aa"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(path.Join(tt.args.list...), func(t *testing.T) {
			if got := ContainsIgnoreCase(tt.args.list, tt.args.str); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilter(t *testing.T) {
	cases := []struct {
		input   string
		ignores []rune
		expect  string
	}{
		{``, nil, ``},
		{`abcd`, nil, `abcd`},
		{`ab,cd,ef`, []rune{','}, `abcdef`},
		{`ab, cd,ef`, []rune{',', ' '}, `abcdef`},
		{`ab, cd, ef`, []rune{',', ' '}, `abcdef`},
		{`ab, cd, ef, `, []rune{',', ' '}, `abcdef`},
	}

	for _, each := range cases {
		t.Run(each.input, func(t *testing.T) {
			actual := Filter(each.input, func(r rune) bool {
				for _, x := range each.ignores {
					if x == r {
						return true
					}
				}
				return false
			})
			assert.Equal(t, each.expect, actual)
		})
	}
}

func TestHasEmpty(t *testing.T) {
	cases := []struct {
		args   []string
		expect bool
	}{
		{
			args:   []string{"a", "b", "c"},
			expect: true,
		},
		{
			args:   []string{"a", "", "c"},
			expect: false,
		},
		{
			args:   []string{"a"},
			expect: true,
		},
		{
			args:   []string{""},
			expect: false,
		},
		{
			args:   []string{},
			expect: true,
		},
	}

	for _, each := range cases {
		t.Run(path.Join(each.args...), func(t *testing.T) {
			assert.Equal(t, each.expect, NotEmpty(each.args...))
		})
	}
}

func TestNotEmpty(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NotEmpty(tt.args.args...); got != tt.want {
				t.Errorf("NotEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemove(t *testing.T) {
	cases := []struct {
		input  []string
		remove []string
		expect []string
	}{
		{
			input:  []string{"a", "b", "a", "c"},
			remove: []string{"a", "b"},
			expect: []string{"c"},
		},
		{
			input:  []string{"b", "c"},
			remove: []string{"a"},
			expect: []string{"b", "c"},
		},
		{
			input:  []string{"b", "a", "c"},
			remove: []string{"a"},
			expect: []string{"b", "c"},
		},
		{
			input:  []string{},
			remove: []string{"a"},
			expect: []string{},
		},
	}

	for _, each := range cases {
		t.Run(path.Join(each.input...), func(t *testing.T) {
			assert.ElementsMatch(t, each.expect, Remove(each.input, each.remove...))
		})
	}
}

func TestReverse(t *testing.T) {
	cases := []struct {
		input  string
		expect string
	}{
		{
			input:  "abcd",
			expect: "dcba",
		},
		{
			input:  "",
			expect: "",
		},
		{
			input:  "我爱中国",
			expect: "国中爱我",
		},
	}

	for _, each := range cases {
		t.Run(each.input, func(t *testing.T) {
			assert.Equal(t, each.expect, Reverse(each.input))
		})
	}
}

func TestSubstr(t *testing.T) {
	cases := []struct {
		input  string
		start  int
		stop   int
		err    error
		expect string
	}{
		{
			input:  "abcdefg",
			start:  1,
			stop:   4,
			expect: "bcd",
		},
		{
			input:  "我爱中国3000遍，even more",
			start:  1,
			stop:   9,
			expect: "爱中国3000遍",
		},
		{
			input:  "abcdefg",
			start:  -1,
			stop:   4,
			err:    ErrInvalidStartPosition,
			expect: "",
		},
		{
			input:  "abcdefg",
			start:  100,
			stop:   4,
			err:    ErrInvalidStartPosition,
			expect: "",
		},
		{
			input:  "abcdefg",
			start:  1,
			stop:   -1,
			err:    ErrInvalidStopPosition,
			expect: "",
		},
		{
			input:  "abcdefg",
			start:  1,
			stop:   100,
			err:    ErrInvalidStopPosition,
			expect: "",
		},
	}

	for _, each := range cases {
		t.Run(each.input, func(t *testing.T) {
			val, err := Substr(each.input, each.start, each.stop)
			assert.Equal(t, each.err, err)
			if err == nil {
				assert.Equal(t, each.expect, val)
			}
		})
	}
}

func TestTakeOne(t *testing.T) {
	cases := []struct {
		valid  string
		or     string
		expect string
	}{
		{"", "", ""},
		{"", "1", "1"},
		{"1", "", "1"},
		{"1", "2", "1"},
	}

	for _, each := range cases {
		t.Run(each.valid, func(t *testing.T) {
			actual := TakeOne(each.valid, each.or)
			assert.Equal(t, each.expect, actual)
		})
	}
}
