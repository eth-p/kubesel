package cli

import (
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
		errorPrinter().PrintCommandError(
			RootCommand.ErrOrStderr(),
			cmd,
			err,
		)

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
