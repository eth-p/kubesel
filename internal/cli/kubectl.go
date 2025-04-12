package cli

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

var (
	kubectlPath = sync.OnceValues(func() (string, error) {
		return exec.LookPath("kubectl")
	})
)

// runKubectl runs the kubectl executable found on the PATH and returns
// its output.
func runKubectl(ctx context.Context, args []string) (string, error) {
	kubectl, err := kubectlPath()
	if err != nil {
		return "", fmt.Errorf("cannot find kubectl executable: %w", err)
	}

	var stdout strings.Builder
	var stderr strings.Builder

	// Run kubectl.
	proc := exec.CommandContext(ctx, kubectl, args...)
	proc.Stdin = nil
	proc.Stdout = &stdout
	proc.Stderr = &stderr
	err = proc.Run()

	// Return stdout/stderr depending on exit code.
	if err != nil {
		return stderr.String(), fmt.Errorf("kubectl: %w", err)
	}

	return stdout.String(), nil
}
