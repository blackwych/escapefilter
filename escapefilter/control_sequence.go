package escapefilter

import (
	"bufio"
	"errors"
	"strconv"
	"strings"
)

// controlSequence represents control sequence, which is CSI <param>* <intermediate>* <final>.
type controlSequence struct {
	param        string
	intermediate string
	final        string
}

// String returns the string representation of the control sequence.
func (c *controlSequence) String() string {
	var sb strings.Builder

	sb.WriteString("\u001B[")
	sb.WriteString(c.param)
	sb.WriteString(c.intermediate)
	sb.WriteString(c.final)

	return sb.String()
}

// invalidControlSequence represents the error of parsing control sequence.
var invalidControlSequence = errors.New("invalid control sequence")

// readControlSequence reads a control sequence from the Reader.
func readControlSequence(rd *bufio.Reader) (*controlSequence, error) {
	cs := &controlSequence{}

	const (
		PARAMETER = iota
		INTERMEDIATE
		FINAL
		END
	)

	for state := PARAMETER; state != END; {
		r, s, err := rd.ReadRune()
		if s == 0 {
			return cs, err
		}

		switch state {
		case PARAMETER:
			if '\u0030' <= r && r <= '\u003F' {
				cs.param += string(r)
			} else {
				rd.UnreadRune()
				state = INTERMEDIATE
			}
		case INTERMEDIATE:
			if '\u0020' <= r && r <= '\u002F' {
				cs.intermediate += string(r)
			} else {
				rd.UnreadRune()
				state = FINAL
			}
		case FINAL:
			if '\u0040' <= r && r <= '\u007E' {
				cs.final += string(r)
				state = END
			} else {
				rd.UnreadRune()
				return cs, invalidControlSequence
			}
		}
	}

	return cs, nil
}

// parseInt converts numerical string into an int value.
// parseInt returns the default value if the string is empty or unparsable.
func parseInt(str string, def int) (int, error) {
	if str == "" {
		return def, nil
	}

	n, err := strconv.Atoi(str)
	if err != nil {
		return def, err
	}

	return n, nil
}

// processControlSequence applys the effects of the control sequence to screen.
func processControlSequence(s *Screen, cs *controlSequence) error {
	switch cs.final {
	case "A": // CUU
		n, err := parseInt(cs.param, 1)
		if err != nil {
			return invalidControlSequence
		}

		s.MoveCursor(s.Row()-n, s.Col())
	case "B": // CUD
		n, err := parseInt(cs.param, 1)
		if err != nil {
			return invalidControlSequence
		}

		s.MoveCursor(s.Row()+n, s.Col())
	case "C": // CUF
		n, err := parseInt(cs.param, 1)
		if err != nil {
			return invalidControlSequence
		}

		s.MoveCursor(s.Row(), s.Col()+n)
	case "D": // CUB
		n, err := parseInt(cs.param, 1)
		if err != nil {
			return invalidControlSequence
		}

		s.MoveCursor(s.Row(), s.Col()-n)
	case "E": // CNL
		n, err := parseInt(cs.param, 1)
		if err != nil {
			return invalidControlSequence
		}

		s.MoveCursor(s.Row()+n, 1)
	case "F": // CPL
		n, err := parseInt(cs.param, 1)
		if err != nil {
			return invalidControlSequence
		}

		s.MoveCursor(s.Row()-n, 1)
	case "G": // CHA
		n, err := parseInt(cs.param, 1)
		if err != nil {
			return invalidControlSequence
		}

		s.MoveCursor(s.Row(), n)
	case "H": // CUP
		params := strings.Split(cs.param, ";")
		if len(params) != 2 {
			return invalidControlSequence
		}

		ns := []int{1, 1}
		for i, param := range params {
			var err error
			ns[i], err = parseInt(param, ns[i])
			if err != nil {
				return invalidControlSequence
			}
		}

		s.MoveCursor(ns[0], ns[1])
	case "I": // CHT
		n, err := parseInt(cs.param, 1)
		if err != nil {
			return invalidControlSequence
		}

		s.MoveCursor(s.Row(), s.NextTabStop(n))
	case "J": // ED
		n, err := parseInt(cs.param, 0)
		if err != nil {
			return invalidControlSequence
		}

		switch n {
		case 0:
			s.EraseScreenAfter()
		case 1:
			s.EraseScreenBefore()
		case 2:
			s.EraseScreen()
		}
	case "K": // EL
		n, err := parseInt(cs.param, 0)
		if err != nil {
			return invalidControlSequence
		}

		switch n {
		case 0:
			s.EraseLineAfter()
		case 1:
			s.EraseLineBefore()
		case 2:
			s.EraseLine()
		}
	case "Z": // CBT
		n, err := parseInt(cs.param, 1)
		if err != nil {
			return invalidControlSequence
		}

		s.MoveCursor(s.Row(), s.PrevTabStop(n))
	default:
		// unsupported, just ignore
	}

	return nil
}
