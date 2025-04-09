package cli

import (
	"fmt"
	"strings"
)

type InvalidFlagError struct {
	Value   string
	Allowed []string
}

func (e *InvalidFlagError) Error() string {
	return fmt.Sprintf(
		"invalid option: %v\nallowed options are: %v",
		e.Value,
		strings.Join(e.Allowed, ", "),
	)
}

// ShellFlag is a flag that accepts known shell types.
type ShellFlag string

const (
	ShellTypeFish ShellFlag = "fish"
)

// String implements [pflag.Value].
func (f *ShellFlag) String() string {
	value := string(*f)
	if value == "" {
		return "none"
	}

	return value
}

// Set implements [pflag.Value].
func (f *ShellFlag) Set(v string) error {
	switch v {
	case string(ShellTypeFish):
		*f = ShellTypeFish
		return nil

	default:
		return &InvalidFlagError{
			Value: v,
			Allowed: []string{
				string(ShellTypeFish),
			},
		}
	}
}

// Type implements [pflag.Value].
func (f *ShellFlag) Type() string {
	return "shell"
}
