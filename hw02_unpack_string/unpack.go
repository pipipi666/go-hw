package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	// Place your code here.
	sb := strings.Builder{}
	lastChar := ""

	for _, char := range s {
		str := string(char)
		num, err := strconv.Atoi(str)

		if err == nil {
			if len(lastChar) > 0 {
				sb.WriteString(strings.Repeat(lastChar, num))
				lastChar = ""
			} else {
				return "", ErrInvalidString
			}
		} else {
			if len(lastChar) > 0 {
				sb.WriteString(lastChar)
			}

			lastChar = str
		}
	}

	if len(lastChar) > 0 {
		sb.WriteString(lastChar)
	}

	return sb.String(), nil
}
