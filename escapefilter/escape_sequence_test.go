package escapefilter

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"testing"
	"unicode/utf8"
)

func Test_readEscapeSequence(t *testing.T) {
	tests := []struct {
		str     string
		r       rune
		next    rune
		isError bool
	}{
		// "\u001B[" is for the sake of clarity, supposed to have benn read before.
		// next: '\u0000' means Reader should be at EOF after read.
		{
			str:     "\u001BNabc",
			r:       'N',
			next:    'a',
			isError: false,
		},
		{
			str:     "\u001B]abc",
			r:       ']',
			next:    'a',
			isError: false,
		},
		{
			str:     "\u001BN",
			r:       'N',
			next:    '\u0000',
			isError: false,
		},
		{
			str:     "\u001B[",
			r:       '[',
			next:    '\u0000',
			isError: false,
		},
		{
			str:     "\u001B",
			r:       utf8.RuneError,
			next:    '\u0000',
			isError: true,
		},
		{
			str:     "\u001B0",
			r:       utf8.RuneError,
			next:    '0',
			isError: true,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("str=%q", tt.str), func(t *testing.T) {
			rd := bufio.NewReader(strings.NewReader(tt.str[1:]))
			r, err := readEscapeSequence(rd)

			if r != tt.r  {
				t.Errorf("readEscapeSequence() should return %q, got %q", tt.r, r)
			}

			isError := err != nil
			switch {
			case tt.isError && !isError:
				t.Errorf("readEscapeSequence() should return error")
			case !tt.isError && isError:
				t.Errorf("readEscapeSequence() should return not error, got %#v", err)
			}

			r, s, err := rd.ReadRune()
			if err != nil && err != io.EOF {
				t.Fatalf("unexpected error occurred: %v", err)
			}

			eof := s == 0 && err == io.EOF
			switch {
			case tt.next == '\u0000' && !eof:
				t.Errorf("next rune should be EOF, got '%q'", r)
			case tt.next != '\u0000' && eof:
				t.Errorf("next rune should be '%q', got EOF", tt.next)
			case tt.next != '\u0000' && r != tt.next:
				t.Errorf("next rune should be '%q', got '%q'", tt.next, r)
			}
		})
	}
}
