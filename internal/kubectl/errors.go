package kubectl

import (
	"fmt"
	"strings"
)

// KubectlError is returned when `kubectl` returned a non-zero exit code.
type KubectlError struct {
	ExitCode int
	Details  string
}

// Error implements error.
func (e *KubectlError) Error() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "kubectl encountered an error (%v)", e.ExitCode)

	if e.Details != "" {
		sb.WriteString("\n\n")
		sb.WriteString(e.Details)
	}

	return sb.String()
}

// ExecTimeoutError is returned when `kubectl` could not complete before the
// context deadline passed.
type ExecTimeoutError struct {
	cause error
}

func (e *ExecTimeoutError) Error() string {
	return "kubectl took too long to run"
}

func (e *ExecTimeoutError) Unwrap() error {
	return e.cause
}

// NotInstalledError is returned when `kubectl` could not be found on the PATH.
type NotInstalledError struct {
}

func (e *NotInstalledError) Error() string {
	return "kubectl is not installed"
}
