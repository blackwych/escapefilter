package escapefilter

import (
	"bufio"
	"errors"
	"io"
	"unicode/utf8"
)

// invalidEscapeSequence represents the error of parsing escape sequence, which is ESC <char>.
var invalidEscapeSequence = errors.New("invalid escape sequence")

// readEscapeSequence reads an escape sequence from the Reader.
func readEscapeSequence(rd *bufio.Reader) (rune, error) {
	r, size, err := rd.ReadRune()
	if size == 0 {
		if err == io.EOF {
			return utf8.RuneError, invalidEscapeSequence
		} else {
			return utf8.RuneError, err
		}
	}

	if !('\u0040' <= r && r <= '\u005F') {
		rd.UnreadRune()
		return utf8.RuneError, invalidEscapeSequence
	}

	return r, nil
}

// processEscapeSequence applys the effects of the escape sequence to the screen.
func processEscapeSequence(s *Screen, rd *bufio.Reader, r rune) error {
	switch r {
	case '[': // CSI
		cs, err := readControlSequence(rd)
		if err != nil {
			if err == invalidControlSequence {
				return nil // just ignore
			} else {
				return err
			}
		}

		if err := processControlSequence(s, cs); err != nil {
			if err == invalidControlSequence {
				return nil // just ignore
			} else {
				return err
			}
		}
	case ']': // OSC
		osc, err := readOperatingSystemCommand(rd)
		if err != nil {
			if err == invalidOperatingSystemCommand {
				return nil // just ignore
			} else {
				return err
			}
		}

		if err := processOperatingSystemCommand(s, osc); err != nil {
			if err == invalidOperatingSystemCommand {
				return nil // just ignore
			} else {
				return err
			}
		}
	default:
		// unsupported, just ignore
	}

	return nil
}
