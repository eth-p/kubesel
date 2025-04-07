package fuzzy

import (
	"errors"
	"fmt"

	fzf "github.com/junegunn/fzf/src"
)

var ErrUserCancelled = errors.New("user cancelled")

type PickOptions struct {
	Query string
}

// Pick opens a fzf TUI for the user pick a single item out of a list.
func Pick(items []string, opts *PickOptions) (string, error) {
	fzfOpts, err := fzf.ParseOptions(true, []string{})
	if err != nil {
		return "", err
	}

	// Create input channel to send data to fzf.
	var result string
	inputCh := make(chan string, 8)

	// Set fzf options.
	fzfOpts.ForceTtyIn = true
	fzfOpts.ClearOnExit = true

	if opts != nil {
		fzfOpts.Query = opts.Query
	}

	// Set fzf input/output.
	fzfOpts.Input = inputCh
	fzfOpts.Printer = func(s string) {
		result = s
	}

	// Start fzf and wait for it to finish.
	go func() {
		for _, item := range items {
			inputCh <- item
		}

		close(inputCh)
	}()

	returnCode, err := fzf.Run(fzfOpts)
	if returnCode == fzf.ExitInterrupt {
		return "", ErrUserCancelled
	}

	if err != nil {
		return "", fmt.Errorf("cannot open fzf: %w", err)
	}

	return result, nil
}
