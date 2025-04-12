package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/charmbracelet/x/ansi"
	"github.com/eth-p/kubesel/internal/cobraerr"
	"github.com/eth-p/kubesel/internal/cobraprint"
	"github.com/eth-p/kubesel/pkg/kubesel"
	"github.com/spf13/cobra"
	"golang.org/x/term"
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
	Color       bool // --color
	OutputIsTTY bool // not a flag
}

const (
	colorFlagName = "color"
	ExitCodeOK    = 0
	ExitCodeError = 1
	ExitCodeHelp  = 10
)

var (
	// Kubesel is the global instance of [kubesel.Kubesel] used by all subcommands.
	Kubesel = sync.OnceValues(kubesel.NewKubesel)

	// hasPrintedHelp is used to determine [ExitCodeHelp] should be returned.
	// This is set when the help function is called via `--help` or
	// `kubesel help`.
	hasPrintedHelp = false
	helpPrinter    = sync.OnceValue(makeHelpPrinter)
)

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

	// Help.
	Command.SetHelpCommandGroupID(
		CommandGroupKubesel,
	)

	Command.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		helpPrinter().PrintCommandHelp(cmd, args)
		hasPrintedHelp = true
	})

	Command.SetUsageFunc(func(cmd *cobra.Command) error {
		return helpPrinter().PrintCommandUsage(cmd)
	})

	// Persistent flags.
	Command.PersistentFlags().BoolVar(
		&GlobalOptions.Color,
		colorFlagName,
		false, // Default is set by DetectTerminal
		"Print with colors",
	)

	Command.PersistentFlags().Lookup(colorFlagName).DefValue = "auto"
}

// DetectTerminalColors changes the default values for some options depending
// on whether kubsel is writing its output to a terminal.
func DetectTerminal() {
	type getFd interface {
		Fd() uintptr
	}

	if getFd, ok := Command.OutOrStdout().(getFd); ok {
		fd := getFd.Fd()
		GlobalOptions.OutputIsTTY = term.IsTerminal(int(fd))
	}

	if GlobalOptions.OutputIsTTY {
		GlobalOptions.Color = true
	}
}

// Run is the entrypoint for the kubesel command-line interface.
func Run(args []string) (int, error) {
	Command.SetArgs(args)
	cmd, err := Command.ExecuteC()

	if err != nil {
		err = cobraerr.Parse(err) // try to parse cobra's unstructured errors

		var sb strings.Builder
		prepareErrorMessage(&sb, cmd, err)
		io.WriteString(Command.ErrOrStderr(), sb.String()) //nolint:errcheck
		return ExitCodeError, err
	}

	if hasPrintedHelp {
		return ExitCodeHelp, err
	}

	return ExitCodeOK, nil
}

func makeHelpPrinter() *cobraprint.HelpPrinter {
	opts := cobraprint.HelpPrinterOptions{
		Indent: "  ",
	}

	if GlobalOptions.Color {
		opts.HeadingColor = ansi.SGR(ansi.BoldAttr)
		opts.CommandNameColor = ansi.SGR(ansi.CyanForegroundColorAttr)
		opts.FlagNameColor = ansi.SGR(ansi.GreenForegroundColorAttr)
		opts.ArgTypeColor = ansi.SGR(ansi.UnderlineAttr)
	}

	return cobraprint.NewHelpPrinter(
		Command.OutOrStdout(),
		opts,
	)
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
