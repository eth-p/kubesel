package cli

import (
	"fmt"
	"io"
	"math/rand/v2"
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
	debugFlagName = "debug"
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
	Debug       bool // --debug (hidden flag)
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

	// gcWait is a WaitGroup that waits for automatic garbage collection to
	// finish running.
	gcWait sync.WaitGroup
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
		hasPrintedHelp = true
		out := cmd.OutOrStdout()
		usage := helpPrinter().PrintCommandHelp(cmd, args)
		_, _ = io.WriteString(out, usage)
	})

	RootCommand.SetUsageFunc(func(cmd *cobra.Command) error {
		out := cmd.OutOrStdout()
		usage := helpPrinter().PrintCommandUsage(cmd)
		_, _ = io.WriteString(out, usage)
		return nil
	})

	// Persistent flags.
	RootCommand.PersistentFlags().BoolVar(
		&GlobalOptions.Color,
		colorFlagName,
		false, // Default is set by DetectTerminal
		"Print with colors",
	)

	RootCommand.PersistentFlags().BoolVar(
		&GlobalOptions.Debug,
		debugFlagName,
		false,
		"Print with colors",
	)

	RootCommand.PersistentFlags().Lookup(colorFlagName).DefValue = "auto"
	RootCommand.PersistentFlags().Lookup(debugFlagName).Hidden = true
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

	return cobraprint.NewHelpPrinter(opts)
}

func makeErrorPrinter() *cobraprint.ErrorPrinter {
	opts := cobraprint.ErrorPrinterOptions{
		Indent:           "  ",
		BlockquoteIndent: "▏ ",
		HelpPrinter:      helpPrinter(),
	}

	if GlobalOptions.Color {
		opts.ErrorCommandColor = ansi.SGR(ansi.BoldAttr, ansi.BrightRedForegroundColorAttr)
		opts.ErrorTextColor = ansi.SGR(ansi.RedForegroundColorAttr)
		opts.TipColor = ansi.SGR(ansi.YellowForegroundColorAttr)
	}

	return cobraprint.NewErrorPrinter(opts)
}

// tryQuickGC has a 1 in 2 chance to run a background garbage collection over
// 5 files. The files checked are nondeterministic, and _eventually_ all files
// will end up checked.
func tryQuickGC(cmd *cobra.Command, args []string) {
	tryGC(2, 5)
}

// tryNormalGC has a 1 in 5 chance to run a background garbage collection over
// 25 files. The files checked are nondeterministic, and _eventually_ all files
// will end up checked.
func tryNormalGC(cmd *cobra.Command, args []string) {
	tryGC(5, 25)
}

func tryGC(chance, maxFiles int) {
	debugf("Trying GC. chance=%d maxFiles=%d\n", chance, maxFiles)
	randResult := rand.IntN(chance)
	if randResult != 0 {
		return
	}

	gcWait.Add(1)
	defer gcWait.Done()

	ksel, err := Kubesel()
	if err != nil {
		debugf("Cannot run GC. err=%v\n", err)
		return
	}

	debugf("Running GC.\n")
	res, err := ksel.GarbageCollect(&kubesel.GarbageCollectOptions{
		MaxFilesToCheck: maxFiles,
	})
	debugf("Finished GC. res=%v err=%v\n", res, err)
}

func debugf(pattern string, args ...any) {
	if GlobalOptions.Debug {
		fmt.Fprintf(os.Stderr, pattern, args...)
	}
}
