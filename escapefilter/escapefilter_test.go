package escapefilter

import (
	"fmt"
	"github.com/andreyvit/diff"
	"github.com/google/go-cmp/cmp"
	"strings"
	"testing"
)

func Test_processRune(t *testing.T) {
	lines := []string{"Hello World", "こんにちはABC世界"}

	tests := []struct {
		r        rune
		expected *Screen
	}{
		{
			r: '\u0007',
			expected: &Screen{
				lines: append([]string{}, lines...),
				row:   2,
				col:   18,
			},
		},
		{
			r: '\u0008',
			expected: &Screen{
				lines: append([]string{}, lines...),
				row:   2,
				col:   17,
			},
		},
		{
			r: '\u0009',
			expected: &Screen{
				lines: append([]string{}, lines...),
				row:   2,
				col:   25,
			},
		},
		{
			r: '\u000A',
			expected: &Screen{
				lines: append([]string{}, lines...),
				row:   3,
				col:   1,
			},
		},
		{
			r: '\u000B',
			expected: &Screen{
				lines: append([]string{}, lines...),
				row:   3,
				col:   18,
			},
		},
		{
			r: '\u000D',
			expected: &Screen{
				lines: append([]string{}, lines...),
				row:   2,
				col:   1,
			},
		},
		{
			r: '\u0020',
			expected: &Screen{
				lines: []string{lines[0], lines[1] + " "},
				row:   2,
				col:   19,
			},
		},
		{
			r: 'X',
			expected: &Screen{
				lines: []string{lines[0], lines[1] + "X"},
				row:   2,
				col:   19,
			},
		},
		{
			r: 'あ',
			expected: &Screen{
				lines: []string{lines[0], lines[1] + "あ"},
				row:   2,
				col:   20,
			},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("r=%U", tt.r), func(t *testing.T) {
			s := &Screen{
				lines: append([]string{}, lines...),
				row:   2,
				col:   18,
			}

			processRune(s, tt.r)

			opt := cmp.AllowUnexported(*tt.expected)
			if diff := cmp.Diff(tt.expected, s, opt); diff != "" {
				t.Errorf("Screen differs from expected\n%s", diff)
			}
		})
	}
}

func Test_New(t *testing.T) {
	s := New()

	if s.screen == nil {
		t.Errorf("screen should not be nil")
	}
}

func Test_EscapeFilter_Load(t *testing.T) {
	source := strings.Join([]string{
		"plain text\r",
		"\x1b[01;32mcolored text\x1b[00m\r",
		"this is erased\b\x1b[K\b\x1b[K\b\x1b[K\b\x1b[K\b\x1b[K\b\x1b[Koverwritten\r",
		"this is erased\r\x1b[13C\x1b[K\r\x1b[12C\x1b[K\r\x1b[11C\x1b[K\r\x1b[10C\x1b[K\r\x1b[9C\x1b[K\r\x1b[8C\x1b[Koverwritten\r",
		"\r",
		"\x1b[Aone line raised above\r",
		"\x1b[8;3Habsolute positioning\r",
		"plain text again",
		"",
		"",
	}, "\n")

	expected := strings.Join([]string{
		"plain text",
		"colored text",
		"this is overwritten",
		"this is overwritten",
		"one line raised above",
		"",
		"",
		"  absolute positioning",
		"plain text again",
		"",
		"",
	}, "\n")

	filter := New()
	filter.Load(strings.NewReader(source))

	if actual := filter.String(); actual != expected {
		t.Errorf("String() differs from expected\n%v", diff.LineDiff(expected, actual))
	}
}
