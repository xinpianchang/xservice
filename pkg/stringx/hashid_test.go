package stringx

import (
	"math"
	"testing"
)

func Test_simpleHashId(t *testing.T) {
	tests := []struct {
		name          string
		instance      *simpleHashId
		input         int64
		wantEncodeErr bool
		wantDecodeErr bool
	}{
		{
			name:          "simple hashid 1",
			instance:      NewSimpleHashId("123", 8),
			input:         1,
			wantEncodeErr: false,
			wantDecodeErr: false,
		},
		{
			name:          "simple hashid max",
			instance:      NewSimpleHashId("123", 8),
			input:         math.MaxInt64,
			wantEncodeErr: false,
			wantDecodeErr: false,
		},
		{
			name:          "simple hashid min",
			instance:      NewSimpleHashId("123", 8),
			input:         math.MinInt64,
			wantEncodeErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.instance.EncodeInt64(tt.input)
			if (err != nil) != tt.wantEncodeErr {
				t.Errorf("simpleHashId.EncodeInt64() error = %v, wantErr %v", err, tt.wantEncodeErr)
				return
			}

			if tt.wantEncodeErr {
				return
			}

			t.Logf("%s, input: %v, out: %v", tt.name, tt.input, got)

			out, err := tt.instance.DecodeInt64(got)
			if (err != nil) != tt.wantDecodeErr {
				t.Errorf("simpleHashId.DecodeInt64() error = %v, wantErr %v", err, tt.wantEncodeErr)
				return
			}

			if out != tt.input {
				t.Errorf("simpleHashId.DecodeInt64() = %v, want %v", out, tt.input)
				return
			}
		})
	}
}
