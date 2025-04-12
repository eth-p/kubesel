package cli

import (
	"github.com/eth-p/kubesel/internal/cobraerr"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// showHelpIfNoArgs is a cobra command that will show the command's help message
// if no arguments are provided, or return an [UnknownCommandError] otherwise.
func showHelpIfNoArgs(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return pflag.ErrHelp
	}

	return &cobraerr.UnknownCommandError{
		Command:       args[0],
		ParentCommand: cmd.Name(),
	}
}
