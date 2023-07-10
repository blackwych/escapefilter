package escapefilter

import (
	"bufio"
	"io"
)

// processRune applys the effects to the Screen if the rune is a control character.
// processRune puts the rune to the Screen otherwise.
func processRune(s *Screen, r rune) {
	switch r {
	case '\u0008': // BS
		s.MoveCursor(s.Row(), s.Col()-1)
	case '\u0009': // HT
		ts := s.NextTabStop(1)
		s.MoveCursor(s.Row(), ts)
	case '\u000A': // LF
		// work as CR+LF
		s.MoveCursor(s.Row()+1, 1)
	case '\u000B': // VT
		s.MoveCursor(s.Row()+1, s.Col())
	case '\u000D': // CR
		s.MoveCursor(s.Row(), 1)
	default:
		s.PutRune(r)
	}
}

// EscapeFilter stores virtual screen and process text files contains ANSI escape code.
type EscapeFilter struct {
	screen *Screen
}

// New returns a new EscapeFilter.
func New() *EscapeFilter {
	return &EscapeFilter{screen: NewScreen()}
}

// Load loads contents from the Reader.
func (f *EscapeFilter) Load(rd io.Reader) error {
	brd := bufio.NewReader(rd)

	for {
		r, s, err := brd.ReadRune()
		if s == 0 {
			if err == io.EOF {
				break
			}
			return err
		}

		if r == '\u001B' {
			es, err := readEscapeSequence(brd)
			if err != nil {
				if err == invalidEscapeSequence {
					// ignore ESC
					continue
				} else {
					return err
				}
			}

			processEscapeSequence(f.screen, brd, es)
		} else {
			processRune(f.screen, r)
		}
	}

	return nil
}

// String returns the current screen content.
func (f *EscapeFilter) String() string {
	return f.screen.String()
}
