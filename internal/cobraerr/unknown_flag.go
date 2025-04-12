package cobraerr

import (
	"fmt"
	"strings"
)

// UnknownFlagError is returned when the user specifies an unknown command-line
// flag.
type UnknownFlagError struct {
	IsShorthandFlag bool
	Flag            string
	FlagSet         string
}

func (e *UnknownFlagError) Error() string {
	if e.IsShorthandFlag {
		return fmt.Sprintf("unknown shorthand flag: %q in %s", e.Flag, e.FlagSet)
	} else {
		return fmt.Sprintf("unknown flag: --%s", e.Flag)
	}
}

func ParseUnknownFlagError(err error) (*UnknownFlagError, bool) {
	errStr := err.Error()

	// unknown shorthand flag: 'x' in -x
	rest, found := strings.CutPrefix(errStr, "unknown shorthand flag: ")
	if found {
		unknownFlag, flagSet, ok := parseUnknownShorthandFlag(rest)
		if ok {
			return &UnknownFlagError{
				IsShorthandFlag: true,
				Flag:            unknownFlag,
				FlagSet:         flagSet,
			}, true
		}

		return nil, false
	}

	// unknown flag: --x
	rest, found = strings.CutPrefix(errStr, "unknown flag: --")
	if found {
		return &UnknownFlagError{
			IsShorthandFlag: false,
			Flag:            rest,
		}, true
	}

	return nil, false
}

// parseUnknownShorthandFlag parses the shorthand flag name and set of
// shorthand flags out of the following string:
//
//	'x' in -x
//	'\\' in -\
//	'\'' in -'
func parseUnknownShorthandFlag(str string) (string, string, bool) {
	// Parse the `'x'`
	flag, str, ok := parseQuotedString(str)
	if !ok {
		return "", "", false
	}

	// Parse the ` in -x`
	flagSet, ok := strings.CutPrefix(str, " in -")
	if !ok {
		return "", "", false
	}

	return flag, flagSet, ok
}
