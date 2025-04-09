package cli

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"os"

	"al.essio.dev/pkg/shellescape"
	"github.com/eth-p/kubesel/pkg/kubesel"
	"github.com/spf13/cobra"
)

//go:embed shell-init/*
var initScripts embed.FS

// InitCommand describes the subcommand for creating a new kubesel session.
var InitCommand = cobra.Command{
	RunE: InitCommandMain,

	Use: "init shell",

	Short: "Initialize kubesel in the current shell",
	Long: `
		Generate a shell script that when sourced, will initialize kubesel in
		the current shell.
	`,
	Example: `
		# bash
		source <(kubesel init bash)

		# fish
		kubesel init fish | source
	`,

	Args: cobra.ExactArgs(1),
	ValidArgs: []string{
		"bash",
		"fish",
	},
}

var InitCommandOptions struct {
	InheritExisting bool
}

func init() {
	Command.AddCommand(&InitCommand)

	InitCommand.Flags().BoolVar(
		&InitCommandOptions.InheritExisting,
		"inherit-existing",
		false,
		"quietly exit if already managing kubeconfig",
	)
}

func InitCommandMain(cmd *cobra.Command, args []string) error {
	ksel, err := Kubesel()
	if err != nil {
		return err
	}

	initScript, err := initScripts.ReadFile("shell-init/init." + args[0])
	if err != nil {
		return fmt.Errorf("unsupported shell: %s", args[0])
	}

	managedKubeconfig, err := ksel.GetManagedKubeconfig()
	if !errors.Is(err, kubesel.ErrUnmanaged) {
		if InitCommandOptions.InheritExisting {
			return nil
		}

		return fmt.Errorf("already managing kubeconfig: %s", managedKubeconfig.Path())
	}

	// Replace `@@KUBESEL@@` with quoted path to kubesel executable and print
	// it to standard out.
	templatedInitScript := bytes.ReplaceAll(
		initScript,
		[]byte("@@KUBESEL@@"),
		[]byte(shellescape.Quote(os.Args[0])),
	)

	cmd.OutOrStdout().Write(templatedInitScript)
	return nil
}
