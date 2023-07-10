package escapefilter

import (
	"bufio"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"io"
	"strings"
	"testing"
)

func Test_controlSequence_String(t *testing.T) {
	tests := []struct {
		cs       *controlSequence
		expected string
	}{
		{
			cs:       &controlSequence{final: "A"},
			expected: "\u001B[A",
		},
		{
			cs:       &controlSequence{param: "1;31", final: "m"},
			expected: "\u001B[1;31m",
		},
		{
			cs:       &controlSequence{param: "0", intermediate: "\"", final: "q"},
			expected: "\u001B[0\"q",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("controlSequence=%#v", tt.cs), func(t *testing.T) {
			if str := tt.cs.String(); str != tt.expected {
				t.Errorf("String() should return %q, got %q", tt.expected, str)
			}
		})
	}
}

func Test_readControlSequence(t *testing.T) {
	tests := []struct {
		str     string
		cs      *controlSequence
		next    rune
		isError bool
	}{
		// "\u001B[" is for the sake of clarity, supposed to have benn read before.
		// next: '\u0000' means Reader should be at EOF after read.
		{
			str:     "\u001B[Aabc",
			cs:      &controlSequence{final: "A"},
			next:    'a',
			isError: false,
		},
		{
			str:     "\u001B[1;31mdef",
			cs:      &controlSequence{param: "1;31", final: "m"},
			next:    'd',
			isError: false,
		},
		{
			str:     "\u001B[0\"qghi",
			cs:      &controlSequence{param: "0", intermediate: "\"", final: "q"},
			next:    'g',
			isError: false,
		},
		{
			str:     "\u001B[A",
			cs:      &controlSequence{final: "A"},
			next:    '\u0000',
			isError: false,
		},
		{
			str:     "\u001B[1;31m",
			cs:      &controlSequence{param: "1;31", final: "m"},
			next:    '\u0000',
			isError: false,
		},
		{
			str:     "\u001B[0\"q",
			cs:      &controlSequence{param: "0", intermediate: "\"", final: "q"},
			next:    '\u0000',
			isError: false,
		},
		{
			str:     "\u001B[",
			cs:      &controlSequence{},
			next:    '\u0000',
			isError: true,
		},
		{
			str:     "\u001B[123",
			cs:      &controlSequence{param: "123"},
			next:    '\u0000',
			isError: true,
		},
		{
			str:     "\u001B[0\"",
			cs:      &controlSequence{param: "0", intermediate: "\""},
			next:    '\u0000',
			isError: true,
		},
		{
			str:     "\u001B[1;31あ",
			cs:      &controlSequence{param: "1;31"},
			next:    'あ',
			isError: true,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("str=%q", tt.str), func(t *testing.T) {
			rd := bufio.NewReader(strings.NewReader(tt.str[2:]))
			cs, err := readControlSequence(rd)

			opt := cmp.AllowUnexported(*tt.cs)
			if diff := cmp.Diff(tt.cs, cs, opt); diff != "" {
				t.Errorf("readControlSequence() differs from expected\n%s", diff)
			}

			isError := err != nil
			switch {
			case tt.isError && !isError:
				t.Errorf("readControlSequence() should return error")
			case !tt.isError && isError:
				t.Errorf("readControlSequence() should return not error, got %#v", err)
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

func Test_processControlSequence(t *testing.T) {
	tests := []struct {
		lines    []string
		row      int
		col      int
		cs       *controlSequence
		expected *Screen
	}{
		{
			lines: []string{"Hello World", "こんにちはABC世界"},
			row:   2,
			col:   18,
			cs:    &controlSequence{param: "7", final: "D"},
			expected: &Screen{
				row:   2,
				col:   11,
				lines: []string{"Hello World", "こんにちはABC世界"},
			},
		},
		{
			lines: []string{"Hello World", "こんにちはABC世界"},
			row:   2,
			col:   11,
			cs:    &controlSequence{param: "0", final: "K"},
			expected: &Screen{
				row:   2,
				col:   11,
				lines: []string{"Hello World", "こんにちは"},
			},
		},
		{
			lines: []string{"Hello World", "こんにちは"},
			row:   2,
			col:   11,
			cs:    &controlSequence{param: "1;31", final: "m"},
			expected: &Screen{
				row:   2,
				col:   11,
				lines: []string{"Hello World", "こんにちは"},
			},
		},
		{
			lines: []string{"Hello World", "こんにちは"},
			row:   2,
			col:   11,
			cs:    &controlSequence{param: "1;31"},
			expected: &Screen{
				row:   2,
				col:   11,
				lines: []string{"Hello World", "こんにちは"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("lines=%q,row=%d,col=%d,cs=%q", tt.lines, tt.row, tt.col, tt.cs), func(t *testing.T) {
			s := &Screen{lines: tt.lines, row: tt.row, col: tt.col}

			err := processControlSequence(s, tt.cs)
			if err != nil {
				t.Errorf("processControlSequence() should not return error, got %#v", err)
			}

			opt := cmp.AllowUnexported(*tt.expected)
			if diff := cmp.Diff(tt.expected, s, opt); diff != "" {
				t.Errorf("Screen differs from expected\n%s", diff)
			}
		})
	}
}
