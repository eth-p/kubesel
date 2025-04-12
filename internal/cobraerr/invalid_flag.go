package cobraerr

import (
	"fmt"
	"strings"
)

// InvalidFlagError is returned when the user specifies an invalid command-line
// flag.
type InvalidFlagError struct {
	Flag          string
	FlagShorthand string
	Value         string
	Cause         string
}

func (e *InvalidFlagError) Error() string {
	return fmt.Sprintf("invalid argument %q for %q flag: %v", e.Value, e.Flag, e.Cause)
}

func ParseInvalidFlagError(err error) (*InvalidFlagError, bool) {
	str := err.Error()

	// Parse the `invalid argument `
	str, ok := strings.CutPrefix(str, "invalid argument ")
	if !ok {
		return nil, false
	}

	// Parse the value.
	value, str, ok := parseQuotedString(str)
	if !ok {
		return nil, false
	}

	// Parse the ` for `
	str, ok = strings.CutPrefix(str, " for ")
	if !ok {
		return nil, false
	}

	// Parse the flag name.
	flagName, str, ok := parseQuotedString(str)
	if !ok {
		return nil, false
	}

	flagShorthand, flag, ok := splitFlagName(flagName)
	if !ok {
		return nil, false
	}

	// Parse the ` flag: `
	str, ok = strings.CutPrefix(str, " flag: ")
	if !ok {
		return nil, false
	}

	// Parse the error cause.
	cause := parseInvalidFlagCause(value, str)

	return &InvalidFlagError{
		Flag:          flag,
		FlagShorthand: flagShorthand,
		Value:         value,
		Cause:         cause,
	}, true
}

// splitFlagName splits a `-l, --long` flag into its short and long names.
func splitFlagName(str string) (string, string, bool) {
	short, long, ok := strings.Cut(str, ", ")
	if !ok {
		short = ""
		long = str
	}

	if short != "" {
		short, ok = strings.CutPrefix(short, "-")
		if !ok {
			return "", "", false
		}
	}

	long, ok = strings.CutPrefix(long, "--")
	if !ok {
		return "", "", false
	}

	return short, long, true
}

// parseInvalidFlagCause converts common, verbose invalid flag errors to
// something more user-friendly.
func parseInvalidFlagCause(value string, str string) string {
	if strings.HasPrefix(str, "strconv.ParseBool:") {
		return fmt.Sprintf("%q is not a boolean", value)
	}

	return str
}
