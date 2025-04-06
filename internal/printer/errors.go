package printer

import (
	"fmt"
	"reflect"
	"strings"
)

// UnknownFieldError is returned by some printers when their options specify
// fields that do not exist in the type they are printing.
type UnknownFieldError struct {
	Available []string
	Field     string
}

func (e *UnknownFieldError) Error() string {
	return fmt.Sprintf(
		"unknown field: %s\navailable fields: %s",
		e.Field,
		strings.Join(e.Available, ", "),
	)
}

func newUnknownFieldError(available ItemStructFieldList, missing string) error {
	availableStrings := make([]string, len(available))
	for i, field := range available {
		availableStrings[i] = normalizeName(field.Name)
	}

	return &UnknownFieldError{
		Available: availableStrings,
		Field:     missing,
	}
}

// InvalidTypeError is returned when a specific type is expected but not
// provided.
type InvalidTypeError struct {
	Expected string
	Actual   reflect.Type
}

func (e *InvalidTypeError) Error() string {
	return fmt.Sprintf(
		"invalid type: expected %s, got %q",
		e.Expected,
		e.Actual.String(),
	)
}
