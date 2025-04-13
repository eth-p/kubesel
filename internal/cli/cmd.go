package cli

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/charmbracelet/x/ansi"
	"github.com/eth-p/kubesel/internal/cobraprint"
	"github.com/eth-p/kubesel/internal/kubectl"
	"github.com/eth-p/kubesel/pkg/kubesel"
	"github.com/spf13/cobra"
)

const (
	CommandGroupInfo       = "Info"
	CommandGroupKubesel    = "Kubesel"
	CommandGroupKubeconfig = "Kubeconfig"
)

const (
	colorFlagName = "color"
	listFlagName  = "list"
)

// RootCommand is the root `kubesel` command.
var RootCommand = cobra.Command{
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

var (
	// Kubesel is the global instance of [kubesel.Kubesel] used by all
	// subcommands.
	Kubesel = sync.OnceValues(kubesel.NewKubesel)

	// Kubectl is the global instance of the [kubectl.Kubectl] wrapper used by
	// all subcommands.
	Kubectl = sync.OnceValues(kubectl.NewKubectlFromPATH)

	// hasPrintedHelp is used to determine if [ExitCodeHelp] should be returned.
	// This is set when the help function is called via `--help` or
	// `kubesel help`.
	hasPrintedHelp = false
	helpPrinter    = sync.OnceValue(makeHelpPrinter)
	errorPrinter   = sync.OnceValue(makeErrorPrinter)
)

func init() {
	// Command groups.
	RootCommand.AddGroup(&cobra.Group{
		ID:    CommandGroupInfo,
		Title: "Informational Commands:",
	})

	RootCommand.AddGroup(&cobra.Group{
		ID:    CommandGroupKubeconfig,
		Title: "Kubeconfig Commands:",
	})

	RootCommand.AddGroup(&cobra.Group{
		ID:    CommandGroupKubesel,
		Title: "Kubesel Commands:",
	})

	// Help.
	RootCommand.SetHelpCommandGroupID(
		CommandGroupKubesel,
	)

	RootCommand.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		helpPrinter().PrintCommandHelp(cmd, args)
		hasPrintedHelp = true
	})

	RootCommand.SetUsageFunc(func(cmd *cobra.Command) error {
		return helpPrinter().PrintCommandUsage(cmd)
	})

	// Persistent flags.
	RootCommand.PersistentFlags().BoolVar(
		&GlobalOptions.Color,
		colorFlagName,
		false, // Default is set by DetectTerminal
		"Print with colors",
	)

	RootCommand.PersistentFlags().Lookup(colorFlagName).DefValue = "auto"
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
		RootCommand.OutOrStdout(),
		opts,
	)
}

func makeErrorPrinter() *cobraprint.ErrorPrinter {
	opts := cobraprint.ErrorPrinterOptions{
		Indent:      "  ",
		HelpPrinter: helpPrinter(),
	}

	if GlobalOptions.Color {
		opts.ErrorCommandColor = ansi.SGR(ansi.BoldAttr, ansi.BrightRedForegroundColorAttr)
		opts.ErrorTextColor = ansi.SGR(ansi.RedForegroundColorAttr)
		opts.TipColor = ansi.SGR(ansi.YellowForegroundColorAttr)
	}

	return cobraprint.NewErrorPrinter(
		opts,
	)
}
