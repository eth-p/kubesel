package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
		set -x KUBECONFIG (kubesel __init --pid=$fish_pid)
	`,

	Args: cobra.NoArgs,

	SilenceErrors: true,
	SilenceUsage:  true,
}

var internalInitCommandOptions struct {
	OwnerPID kubesel.PidType
}

func init() {
	RootCommand.AddCommand(&internalInitCommand)

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

	// Get the owner for the specified PID.
	owner, err := kubesel.OwnerForProcess(internalInitCommandOptions.OwnerPID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "kubesel error creating managed kubeconfig: %v\n", err)
		os.Exit(2)
	}

	// Create the managed kubeconfig.
	var managedKcPath string
	managedKc, err := ksel.CreateManagedKubeconfig(*owner)

	if err == nil {
		// If the managed kubeconfig is new, we need to save it.
		managedKcPath = managedKc.Path()
		err = managedKc.Save()
		if err != nil {
			fmt.Fprintf(os.Stderr, "kubesel error creating managed kubeconfig: %v\n", err)
			os.Exit(2)
		}
	} else if errors.Is(err, kubesel.ErrAlreadyManaged) {
		// If the managed kubeconfig already exists, we'll re-use it.
		managedKcPath = ksel.GetManagedKubeconfigPathForOwner(*owner)
	} else {
		// If there is any other error, fail.
		fmt.Fprintf(os.Stderr, "kubesel error creating managed kubeconfig: %v\n", err)
		os.Exit(2)
	}

	// Print the updated KUBECONFIG.
	var sb strings.Builder
	sb.WriteString(managedKcPath)

	for _, kcPath := range ksel.GetKubeconfigFilePaths() {
		if ksel.IsManagedKubeconfigPath(kcPath) {
			continue
		}

		sb.WriteRune(filepath.ListSeparator)
		sb.WriteString(kcPath)
	}

	fmt.Fprintf(os.Stdout, "%s\n", sb.String())
	return nil
}
