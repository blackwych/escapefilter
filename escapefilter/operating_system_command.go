package escapefilter

import (
	"bufio"
	"errors"
	"strings"
)

// operatingSystemCommand represents operating system command, which is OSC <command> ; <param>* <final>.
type operatingSystemCommand struct {
	command      string
	param        string
	final        string
}

// String returns the string representation of the control sequence.
func (c *operatingSystemCommand) String() string {
	var sb strings.Builder

	sb.WriteString("\u001B]")
	sb.WriteString(c.command)
	sb.WriteString(";")
	sb.WriteString(c.param)
	sb.WriteString(c.final)

	return sb.String()
}

// invalidOperatingSystemCommand represents the error of parsing operating system command.
var invalidOperatingSystemCommand = errors.New("invalid operating system command")

// readOperatingSystemCommand reads a control sequence from the Reader.
func readOperatingSystemCommand(rd *bufio.Reader) (*operatingSystemCommand, error) {
	osc := &operatingSystemCommand{}

	const (
		COMMAND = iota
		PARAMETER
		ST
		END
	)

	for state := COMMAND; state != END; {
		r, s, err := rd.ReadRune()
		if s == 0 {
			return osc, err
		}

		switch state {
		case COMMAND:
			switch {
			case '\u0030' <= r && r <= '\u0039':
				osc.command += string(r)
			case r == ';':
				state = PARAMETER
			default:
				rd.UnreadRune()
				return osc, invalidOperatingSystemCommand
			}
		case PARAMETER:
			switch r {
			case '\u0007': // BEL
				osc.final += string(r)
				state = END
			case '\u001B': // ESC
				state = ST
			default:
				osc.param += string(r)
			}
		case ST:
			if r == '\\' {
				osc.final += "\u001B" + string(r)
				state = END
			} else {
				osc.param += "\u001B" + string(r)
				state = PARAMETER
			}
		}
	}

	return osc, nil
}

// processOperatingSystemCommand applys the effects of the control sequence to screen.
func processOperatingSystemCommand(s *Screen, osc *operatingSystemCommand) error {
	// ignore all
	return nil
}
