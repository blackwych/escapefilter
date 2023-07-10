package escapefilter

import (
	"fmt"
	"github.com/andreyvit/diff"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func Test_NewScreen(t *testing.T) {
	s := NewScreen()

	if len := len(s.lines); len != 0 {
		t.Errorf("lines should be empty, got length %d", len)
	}

	if s.row != 1 {
		t.Errorf("row should be %d, got %d", 1, s.row)
	}

	if s.col != 1 {
		t.Errorf("col should be %d, got %d", 1, s.col)
	}
}

func Test_putRune(t *testing.T) {
	tests := []struct {
		line    string
		col     int
		r       rune
		newLine string
		newCol  int
	}{
		{
			line:    "Hello World",
			col:     12,
			r:       'X',
			newLine: "Hello WorldX",
			newCol:  13,
		},
		{
			line:    "Hello World",
			col:     5,
			r:       'X',
			newLine: "HellX World",
			newCol:  6,
		},
		{
			line:    "Hello World",
			col:     15,
			r:       'X',
			newLine: "Hello World   X",
			newCol:  16,
		},
		{
			line:    "こんにちはABC世界",
			col:     18,
			r:       'あ',
			newLine: "こんにちはABC世界あ",
			newCol:  20,
		},
		{
			line:    "こんにちはABC世界",
			col:     5,
			r:       'あ',
			newLine: "こんあちはABC世界",
			newCol:  7,
		},
		{
			line:    "こんにちはABC世界",
			col:     6,
			r:       'あ',
			newLine: "こん あ はABC世界",
			newCol:  8,
		},
		{
			line:    "こんにちはABC世界",
			col:     11,
			r:       'あ',
			newLine: "こんにちはあC世界",
			newCol:  13,
		},
		{
			line:    "こんにちはABC世界",
			col:     13,
			r:       'あ',
			newLine: "こんにちはABあ 界",
			newCol:  15,
		},
		{
			line:    "こんにちはABC世界",
			col:     14,
			r:       'あ',
			newLine: "こんにちはABCあ界",
			newCol:  16,
		},
		{
			line:    "こんにちはABC世界",
			col:     20,
			r:       'あ',
			newLine: "こんにちはABC世界  あ",
			newCol:  22,
		},
		{
			line:    "こんにちはABC世界",
			col:     18,
			r:       '\u0000',
			newLine: "こんにちはABC世界",
			newCol:  18,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("line=%s,col=%d,r=%U", tt.line, tt.col, tt.r), func(t *testing.T) {
			newLine, newCol := putRune(tt.line, tt.col, tt.r)

			if newLine != tt.newLine {
				t.Errorf("newLine differs from expected\n%v", diff.LineDiff(tt.newLine, newLine))
			}

			if newCol != tt.newCol {
				t.Errorf("newCol should be %d, got %d", tt.newCol, newCol)
			}
		})
	}
}

func Test_Screen_PutRune(t *testing.T) {
	lines := []string{"Hello World", "こんにちはABC世界"}

	tests := []struct {
		row      int
		col      int
		r        rune
		expected *Screen
	}{
		{
			row: 1,
			col: 5,
			r:   'X',
			expected: &Screen{
				lines: []string{"HellX World", "こんにちはABC世界"},
				row:   1,
				col:   6,
			},
		},
		{
			row: 1,
			col: 15,
			r:   'X',
			expected: &Screen{
				lines: []string{"Hello World   X", "こんにちはABC世界"},
				row:   1,
				col:   16,
			},
		},
		{
			row: 2,
			col: 6,
			r:   'あ',
			expected: &Screen{
				lines: []string{"Hello World", "こん あ はABC世界"},
				row:   2,
				col:   8,
			},
		},
		{
			row: 2,
			col: 18,
			r:   'あ',
			expected: &Screen{
				lines: []string{"Hello World", "こんにちはABC世界あ"},
				row:   2,
				col:   20,
			},
		},
		{
			row: 3,
			col: 1,
			r:   'あ',
			expected: &Screen{
				lines: []string{"Hello World", "こんにちはABC世界", "あ"},
				row:   3,
				col:   3,
			},
		},
		{
			row: 4,
			col: 2,
			r:   'あ',
			expected: &Screen{
				lines: []string{"Hello World", "こんにちはABC世界", "", " あ"},
				row:   4,
				col:   4,
			},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("row=%d,col=%d,r=%q", tt.row, tt.col, tt.r), func(t *testing.T) {
			s := &Screen{lines: append([]string{}, lines...), row: tt.row, col: tt.col}
			s.PutRune(tt.r)

			opt := cmp.AllowUnexported(*tt.expected)
			if diff := cmp.Diff(tt.expected, s, opt); diff != "" {
				t.Errorf("PutRune() differs from expected\n%s", diff)
			}
		})
	}
}

func Test_Screen_MoveCursor(t *testing.T) {
	tests := []struct {
		row int
		col int
	}{
		{row: 1, col: 1},
		{row: 1, col: 1},
		{row: 2, col: 4},
		{row: 2, col: 4},
		{row: 3, col: 12},
		{row: 4, col: 50},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("row=%d,col=%d", tt.row, tt.col), func(t *testing.T) {
			s := NewScreen()
			s.MoveCursor(tt.row, tt.col)

			if row := s.Row(); row != tt.row {
				t.Errorf("Row() should return %d, got %d", tt.row, row)
			}

			if col := s.Col(); col != tt.col {
				t.Errorf("Col() should return %d, got %d", tt.col, col)
			}
		})
	}
}

func Test_Screen_PrevTabStop(t *testing.T) {
	tests := []struct {
		row int
		col int
		n   int
		ts  int
	}{
		{row: 1, col: 1, n: 1, ts: 1},
		{row: 1, col: 20, n: 2, ts: 9},
		{row: 3, col: 12, n: 1, ts: 9},
		{row: 3, col: 12, n: 2, ts: 1},
		{row: 3, col: 12, n: 3, ts: 1},
		{row: 4, col: 50, n: 1, ts: 49},
		{row: 4, col: 50, n: 4, ts: 25},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("row=%d,col=%d,n=%d", tt.row, tt.col, tt.n), func(t *testing.T) {
			s := &Screen{row: tt.row, col: tt.col}

			if ts := s.PrevTabStop(tt.n); ts != tt.ts {
				t.Errorf("PrevTabStop() should return %d, got %d", tt.ts, ts)
			}
		})
	}
}

func Test_Screen_NextTabStop(t *testing.T) {
	tests := []struct {
		row int
		col int
		n   int
		ts  int
	}{
		{row: 1, col: 1, n: 1, ts: 9},
		{row: 1, col: 1, n: 2, ts: 17},
		{row: 2, col: 4, n: 1, ts: 9},
		{row: 2, col: 4, n: 3, ts: 25},
		{row: 3, col: 12, n: 1, ts: 17},
		{row: 4, col: 50, n: 1, ts: 57},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("row=%d,col=%d,n=%d", tt.row, tt.col, tt.n), func(t *testing.T) {
			s := &Screen{row: tt.row, col: tt.col}

			if ts := s.NextTabStop(tt.n); ts != tt.ts {
				t.Errorf("NextTabStop() should return %d, got %d", tt.ts, ts)
			}
		})
	}
}

func Test_Screen_EraseLineAfter(t *testing.T) {
	lines := []string{"Hello World", "こんにちはABC世界"}

	tests := []struct {
		row      int
		col      int
		expected *Screen
	}{
		{
			row: 1,
			col: 1,
			expected: &Screen{
				row:   1,
				col:   1,
				lines: []string{"", "こんにちはABC世界"},
			},
		},
		{
			row: 1,
			col: 8,
			expected: &Screen{
				row:   1,
				col:   8,
				lines: []string{"Hello W", "こんにちはABC世界"},
			},
		},
		{
			row: 2,
			col: 1,
			expected: &Screen{
				row:   2,
				col:   1,
				lines: []string{"Hello World"},
			},
		},
		{
			row: 2,
			col: 6,
			expected: &Screen{
				row:   2,
				col:   6,
				lines: []string{"Hello World", "こん"},
			},
		},
		{
			row: 2,
			col: 17,
			expected: &Screen{
				row:   2,
				col:   17,
				lines: []string{"Hello World", "こんにちはABC世"},
			},
		},
		{
			row: 2,
			col: 19,
			expected: &Screen{
				row:   2,
				col:   19,
				lines: []string{"Hello World", "こんにちはABC世界"},
			},
		},
		{
			row: 4,
			col: 5,
			expected: &Screen{
				row:   4,
				col:   5,
				lines: []string{"Hello World", "こんにちはABC世界"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("row=%d,col=%d", tt.row, tt.col), func(t *testing.T) {
			s := &Screen{lines: append([]string{}, lines...), row: tt.row, col: tt.col}
			s.EraseLineAfter()

			opt := cmp.AllowUnexported(*tt.expected)
			if diff := cmp.Diff(tt.expected, s, opt); diff != "" {
				t.Errorf("%s", diff)
			}
		})
	}
}

func Test_Screen_EraseLineBefore(t *testing.T) {
	lines := []string{"Hello World", "こんにちはABC世界"}

	tests := []struct {
		row      int
		col      int
		expected *Screen
	}{
		{
			row: 1,
			col: 1,
			expected: &Screen{
				row:   1,
				col:   1,
				lines: []string{" ello World", "こんにちはABC世界"},
			},
		},
		{
			row: 1,
			col: 8,
			expected: &Screen{
				row:   1,
				col:   8,
				lines: []string{"        rld", "こんにちはABC世界"},
			},
		},
		{
			row: 2,
			col: 1,
			expected: &Screen{
				row:   2,
				col:   1,
				lines: []string{"Hello World", "  んにちはABC世界"},
			},
		},
		{
			row: 2,
			col: 6,
			expected: &Screen{
				row:   2,
				col:   6,
				lines: []string{"Hello World", "      ちはABC世界"},
			},
		},
		{
			row: 2,
			col: 17,
			expected: &Screen{
				row:   2,
				col:   17,
				lines: []string{"Hello World"},
			},
		},
		{
			row: 2,
			col: 19,
			expected: &Screen{
				row:   2,
				col:   19,
				lines: []string{"Hello World"},
			},
		},
		{
			row: 4,
			col: 5,
			expected: &Screen{
				row:   4,
				col:   5,
				lines: []string{"Hello World", "こんにちはABC世界"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("row=%d,col=%d", tt.row, tt.col), func(t *testing.T) {
			s := &Screen{lines: append([]string{}, lines...), row: tt.row, col: tt.col}
			s.EraseLineBefore()

			opt := cmp.AllowUnexported(*tt.expected)
			if diff := cmp.Diff(tt.expected, s, opt); diff != "" {
				t.Errorf("%s", diff)
			}
		})
	}
}

func Test_Screen_EraseLine(t *testing.T) {
	lines := []string{"Hello World", "こんにちはABC世界"}

	tests := []struct {
		row      int
		col      int
		expected *Screen
	}{
		{
			row: 1,
			col: 8,
			expected: &Screen{
				row:   1,
				col:   8,
				lines: []string{"", "こんにちはABC世界"},
			},
		},
		{
			row: 2,
			col: 8,
			expected: &Screen{
				row:   2,
				col:   8,
				lines: []string{"Hello World"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("row=%d,col=%d", tt.row, tt.col), func(t *testing.T) {
			s := &Screen{lines: append([]string{}, lines...), row: tt.row, col: tt.col}
			s.EraseLine()

			opt := cmp.AllowUnexported(*tt.expected)
			if diff := cmp.Diff(tt.expected, s, opt); diff != "" {
				t.Errorf("%s", diff)
			}
		})
	}
}

func Test_Screen_EraseScreenAfter(t *testing.T) {
	lines := []string{"Hello World", "こんにちはABC世界"}

	tests := []struct {
		row      int
		col      int
		expected *Screen
	}{
		{
			row: 1,
			col: 1,
			expected: &Screen{
				row:   1,
				col:   1,
				lines: []string{},
			},
		},
		{
			row: 1,
			col: 8,
			expected: &Screen{
				row:   1,
				col:   8,
				lines: []string{"Hello W"},
			},
		},
		{
			row: 2,
			col: 1,
			expected: &Screen{
				row:   2,
				col:   1,
				lines: []string{"Hello World"},
			},
		},
		{
			row: 2,
			col: 6,
			expected: &Screen{
				row:   2,
				col:   6,
				lines: []string{"Hello World", "こん"},
			},
		},
		{
			row: 2,
			col: 17,
			expected: &Screen{
				row:   2,
				col:   17,
				lines: []string{"Hello World", "こんにちはABC世"},
			},
		},
		{
			row: 2,
			col: 19,
			expected: &Screen{
				row:   2,
				col:   19,
				lines: []string{"Hello World", "こんにちはABC世界"},
			},
		},
		{
			row: 4,
			col: 5,
			expected: &Screen{
				row:   4,
				col:   5,
				lines: []string{"Hello World", "こんにちはABC世界"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("row=%d,col=%d", tt.row, tt.col), func(t *testing.T) {
			s := &Screen{lines: append([]string{}, lines...), row: tt.row, col: tt.col}
			s.EraseScreenAfter()

			opt := cmp.AllowUnexported(*tt.expected)
			if diff := cmp.Diff(tt.expected, s, opt); diff != "" {
				t.Errorf("%s", diff)
			}
		})
	}
}

func Test_Screen_EraseScreenBefore(t *testing.T) {
	lines := []string{"Hello World", "こんにちはABC世界"}

	tests := []struct {
		row      int
		col      int
		expected *Screen
	}{
		{
			row: 1,
			col: 1,
			expected: &Screen{
				row:   1,
				col:   1,
				lines: []string{" ello World", "こんにちはABC世界"},
			},
		},
		{
			row: 1,
			col: 8,
			expected: &Screen{
				row:   1,
				col:   8,
				lines: []string{"        rld", "こんにちはABC世界"},
			},
		},
		{
			row: 2,
			col: 1,
			expected: &Screen{
				row:   2,
				col:   1,
				lines: []string{"", "  んにちはABC世界"},
			},
		},
		{
			row: 2,
			col: 6,
			expected: &Screen{
				row:   2,
				col:   6,
				lines: []string{"", "      ちはABC世界"},
			},
		},
		{
			row: 2,
			col: 17,
			expected: &Screen{
				row:   2,
				col:   17,
				lines: []string{},
			},
		},
		{
			row: 2,
			col: 19,
			expected: &Screen{
				row:   2,
				col:   19,
				lines: []string{},
			},
		},
		{
			row: 4,
			col: 5,
			expected: &Screen{
				row:   4,
				col:   5,
				lines: []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("row=%d,col=%d", tt.row, tt.col), func(t *testing.T) {
			s := &Screen{lines: append([]string{}, lines...), row: tt.row, col: tt.col}
			s.EraseScreenBefore()

			opt := cmp.AllowUnexported(*tt.expected)
			if diff := cmp.Diff(tt.expected, s, opt); diff != "" {
				t.Errorf("%s", diff)
			}
		})
	}
}

func Test_Screen_EraseScreen(t *testing.T) {
	lines := []string{"Hello World", "こんにちはABC世界"}

	tests := []struct {
		row      int
		col      int
		expected *Screen
	}{
		{
			row: 1,
			col: 8,
			expected: &Screen{
				row:   1,
				col:   8,
				lines: []string{},
			},
		},
		{
			row: 2,
			col: 8,
			expected: &Screen{
				row:   2,
				col:   8,
				lines: []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("row=%d,col=%d", tt.row, tt.col), func(t *testing.T) {
			s := &Screen{lines: append([]string{}, lines...), row: tt.row, col: tt.col}
			s.EraseScreen()

			opt := cmp.AllowUnexported(*tt.expected)
			if diff := cmp.Diff(tt.expected, s, opt); diff != "" {
				t.Errorf("Screen differs from expected\n%s", diff)
			}
		})
	}
}

func Test_Screen_String(t *testing.T) {
	s := &Screen{lines: []string{"Hello World", "こんにちはABC世界"}, row: 2, col: 18}
	expected := "Hello World\nこんにちはABC世界"

	if actual := s.String(); actual != expected {
		t.Errorf("String() differs from expected\n%v", diff.LineDiff(expected, actual))
	}
}
