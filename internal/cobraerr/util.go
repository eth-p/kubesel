package cobraerr

import (
	"strconv"
	"unicode/utf8"
)

func firstRune(str string) (rune, bool) {
	for _, c := range str {
		return c, true
	}

	return '\x00', false
}

// parseQuotedString parses a quoted string, returning the unquoted result,
// the remaining characters, and a boolean indicating if parsing was successful.
func parseQuotedString(str string) (string, string, bool) {
	originalStr := str

	// Parse the quote char.
	quoteChar, ok := firstRune(str)
	if !ok || (quoteChar != '\'' && quoteChar != '"') {
		return "", originalStr, false
	}

	str = str[utf8.RuneLen(quoteChar):]

	// Iterate through the string, parsing each rune until reaching the quote
	// character again.
	var chars []rune
	for {
		var ch rune
		var err error

		prevStr := str
		ch, _, str, err = strconv.UnquoteChar(str, byte(quoteChar))
		if err != nil {
			nextCh, ok := firstRune(prevStr)
			if ok && nextCh == quoteChar {
				str = prevStr[utf8.RuneLen(quoteChar):]
				break
			}

			return "", originalStr, false
		}

		chars = append(chars, ch)
	}

	return string(chars), str, true
}
