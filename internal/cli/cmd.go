package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/eth-p/kubesel/internal/cobraerr"
	"github.com/eth-p/kubesel/pkg/kubesel"
	"github.com/spf13/cobra"
)

const (
	CommandGroupInfo       = "Info"
	CommandGroupKubesel    = "Kubesel"
	CommandGroupKubeconfig = "Kubeconfig"
)

// Command is the root `kubesel` command.
var Command = cobra.Command{
	Use: filepath.Base(os.Args[0]),

	CompletionOptions: cobra.CompletionOptions{
		HiddenDefaultCmd: initScriptLoadsCompletions,
	},

	SilenceErrors: true,
	SilenceUsage:  true,
}

var GlobalOptions struct {
	Color bool
}

var Kubesel = sync.OnceValues(kubesel.NewKubesel)

func init() {
	// Command groups.
	Command.AddGroup(&cobra.Group{
		ID:    CommandGroupInfo,
		Title: "Informational Commands:",
	})

	Command.AddGroup(&cobra.Group{
		ID:    CommandGroupKubeconfig,
		Title: "Kubeconfig Commands:",
	})

	Command.AddGroup(&cobra.Group{
		ID:    CommandGroupKubesel,
		Title: "Kubesel Commands:",
	})

	// Misc.
	Command.SetHelpCommandGroupID(
		CommandGroupKubesel,
	)

	// Persistent flags.
	Command.PersistentFlags().BoolVar(
		&GlobalOptions.Color,
		"color",
		true, // TODO: auto
		"Print with colors",
	)
}

// Run is the entrypoint for the kubesel command-line interfa
func Run(args []string) (int, error) {
	Command.SetArgs(args)
	cmd, err := Command.ExecuteC()

	if err != nil {
		err = cobraerr.Parse(err) // try to parse cobra's unstructured errors

		var sb strings.Builder
		prepareErrorMessage(&sb, cmd, err)
		io.WriteString(Command.ErrOrStderr(), sb.String()) //nolint:errcheck
		return 1, err
	}

	_ = cmd
	return 0, nil
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
