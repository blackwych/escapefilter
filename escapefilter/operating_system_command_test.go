package escapefilter

import (
	"bufio"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"io"
	"strings"
	"testing"
)

func Test_operatingSystemCommand_String(t *testing.T) {
	tests := []struct {
		osc      *operatingSystemCommand
		expected string
	}{
		{
			osc:      &operatingSystemCommand{command: "0", param: "Hello World", final: "\u0007"},
			expected: "\u001B]0;Hello World\u0007",
		},
		{
			osc:      &operatingSystemCommand{command: "0", param: "Hello World", final: "\u001B\\"},
			expected: "\u001B]0;Hello World\u001B\\",
		},
		{
			osc:      &operatingSystemCommand{command: "4", param: "rgb:127/127/127", final: "\u0007"},
			expected: "\u001B]4;rgb:127/127/127\u0007",
		},
		{
			osc:      &operatingSystemCommand{command: "4", param: "rgb:127/127/127", final: "\u001B\\"},
			expected: "\u001B]4;rgb:127/127/127\u001B\\",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("operatingSystemCommand=%#v", tt.osc), func(t *testing.T) {
			if str := tt.osc.String(); str != tt.expected {
				t.Errorf("String() should return %q, got %q", tt.expected, str)
			}
		})
	}
}

func Test_readOperatingSystemCommand(t *testing.T) {
	tests := []struct {
		str     string
		osc     *operatingSystemCommand
		next    rune
		isError bool
	}{
		// "\u001B]" is for the sake of clarity, supposed to have benn read before.
		// next: '\u0000' means Reader should be at EOF after read.
		{
			str:     "\u001B]0;Hello World\u0007abc",
			osc:     &operatingSystemCommand{command: "0", param: "Hello World", final: "\u0007"},
			next:    'a',
			isError: false,
		},
		{
			str:     "\u001B]0;Hello World\u001B\\def",
			osc:     &operatingSystemCommand{command: "0", param: "Hello World", final: "\u001B\\"},
			next:    'd',
			isError: false,
		},
		{
			str:     "\u001B]4;rgb:127/127/127\u0007ghi",
			osc:     &operatingSystemCommand{command: "4", param: "rgb:127/127/127", final: "\u0007"},
			next:    'g',
			isError: false,
		},
		{
			str:     "\u001B]4;rgb:127/127/127\u001B\\jkl",
			osc:     &operatingSystemCommand{command: "4", param: "rgb:127/127/127", final: "\u001B\\"},
			next:    'j',
			isError: false,
		},
		{
			str:     "\u001B]0;Hello\u001B World\u001B\\mno",
			osc:     &operatingSystemCommand{command: "0", param: "Hello\u001B World", final: "\u001B\\"},
			next:    'm',
			isError: false,
		},
		{
			str:     "\u001B]0;Hello World\u0007",
			osc:     &operatingSystemCommand{command: "0", param: "Hello World", final: "\u0007"},
			next:    '\u0000',
			isError: false,
		},
		{
			str:     "\u001B]0;Hello World\u001B\\",
			osc:     &operatingSystemCommand{command: "0", param: "Hello World", final: "\u001B\\"},
			next:    '\u0000',
			isError: false,
		},
		{
			str:     "\u001B]4;rgb:127/127/127\u0007",
			osc:      &operatingSystemCommand{command: "4", param: "rgb:127/127/127", final: "\u0007"},
			next:    '\u0000',
			isError: false,
		},
		{
			str:     "\u001B]4;rgb:127/127/127\u001B\\",
			osc:      &operatingSystemCommand{command: "4", param: "rgb:127/127/127", final: "\u001B\\"},
			next:    '\u0000',
			isError: false,
		},
		{
			str:     "\u001B]0;Hello World",
			osc:     &operatingSystemCommand{command: "0", param: "Hello World"},
			next:    '\u0000',
			isError: true,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("str=%q", tt.str), func(t *testing.T) {
			rd := bufio.NewReader(strings.NewReader(tt.str[2:]))
			osc, err := readOperatingSystemCommand(rd)

			opt := cmp.AllowUnexported(*tt.osc)
			if diff := cmp.Diff(tt.osc, osc, opt); diff != "" {
				t.Errorf("readOperatingSystemCommand() differs from expected\n%s", diff)
			}

			isError := err != nil
			switch {
			case tt.isError && !isError:
				t.Errorf("readOperatingSystemCommand() should return error")
			case !tt.isError && isError:
				t.Errorf("readOperatingSystemCommand() should return not error, got %#v", err)
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
