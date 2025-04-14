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
	"github.com/adrg/xdg"
	"github.com/eth-p/kubesel/internal/cobraerr"
	"github.com/spf13/cobra"
)

//go:embed shell-init/init.*
var initScripts embed.FS

// initScriptLoadsCompletions is set to `false` when using the
// `no-init-completions` go build tag.
var initScriptLoadsCompletions = true

// initCommand describes the subcommand for creating a new kubesel session.
var initCommand = cobra.Command{
	PreRun: tryNormalGC,
	RunE:   initCommandMain,

	Use:     "init shell",
	GroupID: "Kubesel",

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
	KubeconfigFiles []string
}

func init() {
	RootCommand.AddCommand(&initCommand)
	initCommand.Flags().StringArrayVar(&InitCommandOptions.KubeconfigFiles, "add-kubeconfigs", []string{}, "kubeconfig files to add")
}

func initCommandMain(cmd *cobra.Command, args []string) error {
	// Find kubeconfig files specified as glob patterns.
	kcFiles, err := resolveKubeconfigFileGlobs()
	if err != nil {
		return err
	}

	// Parse the init script as a Go template.
	initScript, err := getInitScript(os.Args[0], args[0], kcFiles)
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
func getInitScript(argv0 string, shell string, extraKubeconfigFiles []string) (string, error) {
	templateSource, err := initScripts.ReadFile("shell-init/init." + shell)
	if err != nil {
		return "", fmt.Errorf("unsupported shell: %s", shell)
	}

	// Parse the script as a Go template.
	tpl := template.New("init-script").
		Funcs(template.FuncMap{
			"shellquote": shellescape.Quote,
			"join":       strings.Join,
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
		"add_kubeconfigs":    extraKubeconfigFiles,
	})

	if err != nil {
		return "", fmt.Errorf("failed to evaluate %s init script template: %w", shell, err)
	}

	return sb.String(), nil
}

// resolveKubeconfigFileGlobs searches for the files referenced in the
// `--add-kubeconfigs` flag.
func resolveKubeconfigFileGlobs() ([]string, error) {
	if len(InitCommandOptions.KubeconfigFiles) == 0 {
		return nil, nil
	}

	var results []string
	for _, glob := range InitCommandOptions.KubeconfigFiles {
		glob = expandTilde(glob)
		matches, err := filepath.Glob(glob)
		if err != nil {
			return nil, &cobraerr.InvalidFlagError{
				Flag:  "add-kubeconfigs",
				Value: glob,
				Cause: err.Error(),
			}
		}

		for _, file := range matches {
			results = append(results, file)
		}
	}

	if len(results) == 0 {
		return nil, &cobraerr.InvalidFlagError{
			Flag:  "add-kubeconfigs",
			Value: strings.Join(InitCommandOptions.KubeconfigFiles, " "),
			Cause: "no files match the provided glob",
		}
	}

	return results, nil
}

// expandTilde expands file paths that start with `~/` to `$HOME/`.
// If the path does not start with a tilde, it will be returned as-is.
func expandTilde(path string) string {
	cutPath, hasPrefix := strings.CutPrefix(path, "~/")
	if hasPrefix {
		return filepath.Join(xdg.Home, cutPath)
	} else {
		return path
	}
}
