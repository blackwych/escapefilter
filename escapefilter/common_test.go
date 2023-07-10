package escapefilter

import (
	"fmt"
	"testing"
)

func Test_parseInt_Normal(t *testing.T) {
	tests := []struct {
		str string
		def int
		n   int
	}{
		{str: "", def: 1, n: 1},
		{str: "123", def: 1, n: 123},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("str=%s,def=%d", tt.str, tt.def), func(t *testing.T) {
			n, err := parseInt(tt.str, tt.def)
			if err != nil {
				t.Fatalf("parseInt() should not return error, got %v", err)
			}

			if n != tt.n {
				t.Errorf("parseInt() should return %d, got %d", tt.n, n)
			}
		})
	}
}

func Test_parseInt_Error(t *testing.T) {
	tests := []struct {
		str string
		def int
	}{
		{str: "123,456", def: 1},
		{str: "123;456", def: 1},
		{str: "abc", def: 1},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("str=%s,def=%d", tt.str, tt.def), func(t *testing.T) {
			_, err := parseInt(tt.str, tt.def)

			if err == nil {
				t.Errorf("parseInt() should return error")
			}
		})
	}
}
