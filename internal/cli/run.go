package cli

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/eth-p/kubesel/internal/cobraerr"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

const (
	ExitCodeOK    = 0
	ExitCodeError = 1
	ExitCodeHelp  = 10
)

// Run is the entrypoint for the kubesel command-line interface.
func Run(args []string) (int, error) {
	RootCommand.SetArgs(args)
	cmd, err := RootCommand.ExecuteC()

	if err != nil {
		err = cobraerr.Parse(err) // try to parse cobra's unstructured errors

		var sb strings.Builder
		prepareErrorMessage(&sb, cmd, err)
		io.WriteString(RootCommand.ErrOrStderr(), sb.String()) //nolint:errcheck
		return ExitCodeError, err
	}

	if hasPrintedHelp {
		return ExitCodeHelp, err
	}

	return ExitCodeOK, nil
}

// DetectTerminalColors changes the default values for some options depending
// on whether kubsel is writing its output to a terminal.
func DetectTerminal() {
	type getFd interface {
		Fd() uintptr
	}

	if getFd, ok := RootCommand.OutOrStdout().(getFd); ok {
		fd := getFd.Fd()
		GlobalOptions.OutputIsTTY = term.IsTerminal(int(fd))
	}

	if GlobalOptions.OutputIsTTY {
		GlobalOptions.Color = true
	}
}

func prepareErrorMessage(sb *strings.Builder, cmd *cobra.Command, err error) {
	prependCommandToErrorMessage(sb, cmd)

	{
		var flagErr *cobraerr.InvalidFlagError
		if errors.As(err, &flagErr) {
			fmt.Fprintf(sb, "%s\n", err.Error())
			return
		}
	}

	{
		var flagErr *cobraerr.UnknownFlagError
		if errors.As(err, &flagErr) {
			fmt.Fprintf(sb, "%s\n", err.Error())
			return
		}
	}

	{
		var unknownCmdError *cobraerr.UnknownCommandError
		if errors.As(err, &unknownCmdError) {
			fmt.Fprintf(sb, "%s\n", err.Error())

			suggestions := cmd.SuggestionsFor(unknownCmdError.Command)
			if len(suggestions) > 0 {
				fmt.Fprintf(sb, "\nDid you mean:\n")
				for _, suggestion := range suggestions {
					fmt.Fprintf(sb, "  %s\n", suggestion)
				}
			}

			return
		}
	}

	// Unknown
	// fmt.Fprintln(sb, "\n----")
	// fmt.Fprintf(sb, "Got error %T\n", err)
	fmt.Fprintln(sb, err)
}

func prependCommandToErrorMessage(sb *strings.Builder, cmd *cobra.Command) {
	fmt.Fprintf(sb, "%s: ", cmd.CommandPath())
}
