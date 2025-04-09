package cli

import (
	"fmt"
	"os"

	"github.com/eth-p/kubesel/pkg/kubesel"
	"github.com/spf13/cobra"
)

var internalInitCommand = cobra.Command{
	RunE: internalInitCommandMain,

	Use:    "__init",
	Hidden: true,

	Short: "Create a new kubesel session file",
	Example: `
		# fish
		set -x KUBECONFIG (kubesel __init --pid=$fish_pid):"$KUBECONFIG"
	`,

	Args: cobra.NoArgs,

	SilenceErrors: true,
	SilenceUsage:  true,
}

var internalInitCommandOptions struct {
	OwnerPID kubesel.PidType
}

func init() {
	Command.AddCommand(&internalInitCommand)

	internalInitCommand.Flags().Int32Var(
		&internalInitCommandOptions.OwnerPID,
		"pid",
		-1,
		"the PID of the owner",
	)

	internalInitCommand.MarkFlagRequired("pid") // nolint:errcheck
}

func internalInitCommandMain(cmd *cobra.Command, args []string) error {
	ksel, err := Kubesel()
	if err != nil {
		return err
	}

	managedKubeconfig, err := ksel.GetManagedKubeconfig()
	if managedKubeconfig != nil {
		os.Exit(2)
	}

	// Create the managed kubeconfig and print the resulting file to stdout.
	owner, err := kubesel.OwnerForProcess(internalInitCommandOptions.OwnerPID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "kubesel error creating managed kubeconfig: %v\n", err)
		os.Exit(2)
	}

	managedKubeconfig, err = ksel.CreateManagedKubeconfig(*owner)
	if err != nil {
		fmt.Fprintf(os.Stderr, "kubesel error creating managed kubeconfig: %v\n", err)
		os.Exit(2)
	}

	err = managedKubeconfig.Save()
	if err != nil {
		fmt.Fprintf(os.Stderr, "kubesel error creating managed kubeconfig: %v\n", err)
		os.Exit(2)
	}

	fmt.Fprintf(os.Stdout, "%s\n", managedKubeconfig.Path())
	return nil
}
