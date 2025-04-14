package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/eth-p/kubesel/pkg/kubeconfig/kcutils"
	"github.com/eth-p/kubesel/pkg/kubesel"
	"github.com/spf13/cobra"
)

var internalInitCommand = cobra.Command{
	RunE: internalInitCommandMain,

	Use:    "__init",
	Hidden: true,

	Short: "Create a new kubesel session file",
	Example: `
		# bash
		export KUBECONFIG="$(kubesel __init --pid=$$)"

		# zsh
		export KUBECONFIG="$(kubesel __init --pid=$$)"

		# fish
		set -gx KUBECONFIG (kubesel __init --pid=$fish_pid)
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
	managedKc, err := ksel.CreateManagedKubeconfig(*owner)

	// If the managed kubeconfig already exists, we'll re-use it.
	if errors.Is(err, kubesel.ErrAlreadyManaged) {
		printNewKubeconfigEnvVar(ksel, ksel.GetManagedKubeconfigPathForOwner(*owner))
		return nil
	}

	// If an error occurred, exit.
	if err != nil {
		fmt.Fprintf(os.Stderr, "kubesel error creating managed kubeconfig: %v\n", err)
		os.Exit(2)
	}

	// Use the same cluster, user, and namespace we had before.
	currentKc := ksel.GetMergedKubeconfig()
	if currentKc.CurrentContext != nil {
		currentContext := kcutils.FindContext(*currentKc.CurrentContext, currentKc)
		if currentContext != nil {
			if currentContext.Cluster != nil {
				managedKc.SetClusterName(*currentContext.Cluster)
			}
			if currentContext.User != nil {
				managedKc.SetAuthInfoName(*currentContext.User)
			}
			if currentContext.Namespace != nil {
				managedKc.SetNamespace(*currentContext.Namespace)
			}

			err = managedKc.Save()
			if err != nil {
				fmt.Fprintf(os.Stderr, "kubesel error creating managed kubeconfig: %v\n", err)
				os.Exit(2)
			}
		}
	}

	// Print the new KUBECONFIG environment variable.
	printNewKubeconfigEnvVar(ksel, managedKc.Path())
	return nil
}

func printNewKubeconfigEnvVar(ksel *kubesel.Kubesel, managedKcPath string) {
	var sb strings.Builder

	// Add the managed kubeconfig file at the start.
	sb.WriteString(managedKcPath)

	// Add the unmanaged kubeconfig files at the end.
	// Only print a file the first time it appears.
	seenFiles := make(map[string]bool, 0)
	for _, kcPath := range ksel.GetKubeconfigFilePaths() {
		if ksel.IsManagedKubeconfigPath(kcPath) {
			continue
		}

		if !seenFiles[kcPath] {
			seenFiles[kcPath] = true
			sb.WriteRune(filepath.ListSeparator)
			sb.WriteString(kcPath)
		}
	}

	// Print it.
	fmt.Fprintf(os.Stdout, "%s\n", sb.String())
}
