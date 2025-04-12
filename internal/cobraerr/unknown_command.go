package cobraerr

import (
	"fmt"
	"strings"
)

// UnknownCommandError is returned when the user specifies an unknown command
// or subcommand.
type UnknownCommandError struct {
	Command       string
	ParentCommand string
}

func (e *UnknownCommandError) Error() string {
	return fmt.Sprintf("unknown command %q for %q", e.Command, e.ParentCommand)
}

func ParseUnknownCommandError(err error) (*UnknownCommandError, bool) {
	// unknown command "foo" for "test"
	str := err.Error()

	// Parse the `unknown command `
	str, ok := strings.CutPrefix(str, "unknown command ")
	if !ok {
		return nil, false
	}

	// Parse the `"foo"`
	command, str, ok := parseQuotedString(str)
	if !ok {
		return nil, false
	}

	// Parse the ` for `
	str, ok = strings.CutPrefix(str, " for ")
	if !ok {
		return nil, false
	}

	// Parse the `"test"`
	parentCommand, _, ok := parseQuotedString(str)
	if !ok {
		return nil, false
	}

	return &UnknownCommandError{
		Command:       command,
		ParentCommand: parentCommand,
	}, true
}
