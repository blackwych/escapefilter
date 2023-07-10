package escapefilter

import (
	"strconv"
)

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
