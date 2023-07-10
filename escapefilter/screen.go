package escapefilter

import (
	"fmt"
	"github.com/mattn/go-runewidth"
	"strings"
)

// Screen stores character content and a cursor position.
type Screen struct {
	lines []string
	row   int
	col   int
}

// NewScreen returns a new empty Screen.
func NewScreen() *Screen {
	return &Screen{
		row: 1,
		col: 1,
	}
}

// findColumnPosition finds the rune index and column offset that is on the column col.
func findColumnPosition(line string, col int) (index int, offset int) {
	if col <= 0 {
		panic(fmt.Sprintf("col must be >= 1, got %d", col))
	}

	runes := []rune(line)
	c := 1

	for i, r := range runes {
		w := runewidth.RuneWidth(r)

		if c <= col && col < c+w {
			index = i
			offset = col - c
			return
		}

		c += w
	}

	index = len(runes) + (col - c)
	offset = 0
	return
}

// putRune puts a rune to the line and returns the new line and the column position.
func putRune(line string, col int, r rune) (newLine string, newCol int) {
	if col <= 0 {
		panic(fmt.Sprintf("col must be >= 1, got %d", col))
	}

	w := runewidth.RuneWidth(r)
	if w <= 0 {
		return line, col
	}

	l := runewidth.StringWidth(line)

	var sb strings.Builder

	if col > 1 {
		before := runewidth.Truncate(line, col-1, "")
		sb.WriteString(runewidth.FillRight(before, col-1))
	}

	sb.WriteRune(r)

	if col+w <= l {
		after := runewidth.TruncateLeft(line, col+w-1, "")
		sb.WriteString(runewidth.FillLeft(after, l-(col+w)+1))
	}

	return sb.String(), col + w
}

// PutRune puts a rune to the screen.
func (s *Screen) PutRune(r rune) {
	for len(s.lines) < s.row {
		s.lines = append(s.lines, "")
	}

	s.lines[s.row-1], s.col = putRune(s.lines[s.row-1], s.col, r)
}

// Row returns the current row position (1-based).
func (s *Screen) Row() int {
	return s.row
}

// Col returns the current col position (1-based).
func (s *Screen) Col() int {
	return s.col
}

// PrevTabStop returns the n-th last tab stop from the current position.
func (s *Screen) PrevTabStop(n int) int {
	if n <= 0 {
		panic(fmt.Sprintf("n must be >= 1, got %d", n))
	}

	const TABSTOP = 8
	if s.col <= TABSTOP {
		return 1
	}

	if col := ((s.col-2)/TABSTOP-n+1)*TABSTOP + 1; col >= 1 {
		return col
	} else {
		return 1
	}
}

// NextTabStop returns the n-th next tab stop from the current position.
func (s *Screen) NextTabStop(n int) int {
	if n <= 0 {
		panic(fmt.Sprintf("n must be >= 1, got %d", n))
	}

	const TABSTOP = 8
	return ((s.col-1)/TABSTOP+n)*TABSTOP + 1
}

// MoveCursor moves the cursor position to (row, col).
func (s *Screen) MoveCursor(row int, col int) {
	if row <= 0 {
		row = 1
	}

	if col <= 0 {
		col = 1
	}

	s.row = row
	s.col = col
}

// removeExtraBlankLines removes blank lines at the bottom.
func removeExtraBlankLines(lines []string) []string {
	var r int
	for r = len(lines); r >= 1 && lines[r-1] == ""; r-- {
	}
	return append([]string{}, lines[:r]...)
}

// EraseLineAfter erases characters from the current position to the end of the line.
func (s *Screen) EraseLineAfter() {
	if s.row <= len(s.lines) {
		if s.col > 1 {
			s.lines[s.row-1] = runewidth.Truncate(s.lines[s.row-1], s.col-1, "")
		} else {
			s.lines[s.row-1] = ""
		}
	}

	s.lines = removeExtraBlankLines(s.lines)
}

// EraseLineBefore erases characters from the current position to the beginning of the line.
func (s *Screen) EraseLineBefore() {
	if s.row <= len(s.lines) {
		l := runewidth.StringWidth(s.lines[s.row-1])
		if s.col < l {
			rest := runewidth.TruncateLeft(s.lines[s.row-1], s.col, "")
			s.lines[s.row-1] = runewidth.FillLeft(rest, l)
		} else {
			s.lines[s.row-1] = ""
		}
	}

	s.lines = removeExtraBlankLines(s.lines)
}

// EraseLineBefore erases characters in the current row.
func (s *Screen) EraseLine() {
	if s.row <= len(s.lines) {
		s.lines[s.row-1] = ""
	}

	s.lines = removeExtraBlankLines(s.lines)
}

// EraseScreenAfter erases characters from the current position to the end of the screen.
func (s *Screen) EraseScreenAfter() {
	if s.row > len(s.lines) {
		return
	}

	s.lines = append([]string{}, s.lines[:s.row]...)
	s.EraseLineAfter()
}

// EraseScreenBefore erases characters from the current position to the beginning of the screen.
func (s *Screen) EraseScreenBefore() {
	if s.row > len(s.lines) {
		s.EraseScreen()
		return
	}

	for r := 1; r < s.row; r++ {
		s.lines[r-1] = ""
	}

	s.EraseLineBefore()
}

// EraseScreen erases characters in the entire screen.
func (s *Screen) EraseScreen() {
	s.lines = []string{}
}

// String returns string content of the screen.
// If the cursor is farther than the end of the content, additional lines and spaces will be added.
func (s *Screen) String() string {
	var sb strings.Builder

	n := len(s.lines)

	for r := 1; r <= n-1; r++ {
		sb.WriteString(s.lines[r-1])
		sb.WriteRune('\n')
	}

	if s.row == n {
		sb.WriteString(runewidth.FillRight(s.lines[n-1], s.col-1))
	} else {
		sb.WriteString(s.lines[n-1])
	}

	for r := n + 1; r <= s.row; r++ {
		sb.WriteRune('\n')
	}

	return sb.String()
}
