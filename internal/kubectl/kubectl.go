package kubectl

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// Kubectl is a utility for running the locally-installed `kubectl` command.
type Kubectl struct {
	executable string
}

// NewKubectlFromPATH returns a [Kubectl] wrapper for the `kubectl` command
// found on the PATH.
func NewKubectlFromPATH() (*Kubectl, error) {
	executable, err := exec.LookPath("kubectl")
	if err != nil {
		return nil, &NotInstalledError{}
	}

	return NewKubectl(executable)
}

// NewKubectl returns a [Kubectl] wrapper for the specified `kubectl`
// executable.
func NewKubectl(executable string) (*Kubectl, error) {
	return &Kubectl{
		executable: executable,
	}, nil
}

// Exec runs the kubectl executable and returns its standard output.
// If kubectl fails to execute, a [KubectlError] will be returned.
func (k Kubectl) Exec(ctx context.Context, args []string) (string, error) {
	var stdout strings.Builder
	var stderr strings.Builder

	proc := exec.CommandContext(ctx, k.executable, args...)
	proc.Stdin = nil
	proc.Stdout = &stdout
	proc.Stderr = &stderr
	err := proc.Run()

	// Did kubectl exit unsuccessfully?
	var exitError *exec.ExitError
	if errors.As(err, &exitError) {
		if exitError.ExitCode() == -1 && errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return "", &ExecTimeoutError{
				cause: err,
			}
		}

		return stdout.String(), &KubectlError{
			ExitCode: exitError.ExitCode(),
			Details:  stdout.String(),
		}
	}

	// Did something else fail?
	if err != nil {
		return "", fmt.Errorf("unexpected error: %w", err)
	}

	return stdout.String(), nil
}
