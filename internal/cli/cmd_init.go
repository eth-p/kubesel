package cli

import (
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"al.essio.dev/pkg/shellescape"
	"github.com/spf13/cobra"
)

//go:embed shell-init/init.*
var initScripts embed.FS

// initScriptLoadsCompletions is set to `false` when using the
// `no-init-completions` go build tag.
var initScriptLoadsCompletions = true

// InitCommand describes the subcommand for creating a new kubesel session.
var InitCommand = cobra.Command{
	RunE: InitCommandMain,

	Use: "init shell",

	Short: "Initialize kubesel in the current shell",
	Long: `
		Generate a shell script that when sourced, will initialize
		kubesel in the current shell.
	`,
	Example: `
		# bash
		source <(kubesel init bash)

		# zsh
		source <(kubesel init zsh)

		# fish
		kubesel init fish | source
	`,

	Args: cobra.ExactArgs(1),
	ValidArgs: []string{
		"bash",
		"fish",
		"zsh",
	},
}

var InitCommandOptions struct {
}

func init() {
	Command.AddCommand(&InitCommand)
}

func InitCommandMain(cmd *cobra.Command, args []string) error {
	// Parse the init script as a Go template.
	initScript, err := getInitScript(os.Args[0], args[0])
	if err != nil {
		return err
	}

	// Print it to standard out.
	_, err = io.WriteString(cmd.OutOrStdout(), initScript)
	if err != nil {
		return fmt.Errorf("failed to print shell init script: %w", err)
	}

	return nil
}

// getInitScript generates the init script for the specified shell.
//
// Supported shells are:
//   - bash
//   - fish
//   - zsh
func getInitScript(argv0 string, shell string) (string, error) {
	templateSource, err := initScripts.ReadFile("shell-init/init." + shell)
	if err != nil {
		return "", fmt.Errorf("unsupported shell: %s", shell)
	}

	// Parse the script as a Go template.
	tpl := template.New("init-script").
		Funcs(template.FuncMap{
			"shellquote": shellescape.Quote,
		})

	tpl, err = tpl.Parse(string(templateSource))
	if err != nil {
		return "", fmt.Errorf("failed to parse %s init script as template: %w", shell, err)
	}

	// Evaluate the template.
	var sb strings.Builder
	err = tpl.Execute(&sb, map[string]any{
		"kubesel_executable": argv0,
		"kubesel_name":       filepath.Base(argv0),
		"load_completions":   initScriptLoadsCompletions,
	})

	if err != nil {
		return "", fmt.Errorf("failed to evaluate %s init script template: %w", shell, err)
	}

	return sb.String(), nil
}
