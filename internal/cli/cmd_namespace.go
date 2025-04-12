package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var namespaceCommand = cobra.Command{
	RunE: NamespaceCommandMain,

	Aliases: []string{
		"names",
		"nsp",
		"ns",
	},

	Use:     "namespace namespace",
	GroupID: "Kubeconfig",

	Short: "Change to a different namespace",
	Long: `
		Change to a different Kubernetes namespace in the current shell.

		When selecting a namespace, you must use its full name.
	`,
	Example: `
		kubesel namespace kube-system  # full name
	`,

	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: nil,
}

var NamespaceCommandOptions struct {
}

func init() {
	RootCommand.AddCommand(&namespaceCommand)
}

func NamespaceCommandMain(cmd *cobra.Command, args []string) error {
	ksel, err := Kubesel()
	if err != nil {
		return err
	}

	managedConfig, err := ksel.GetManagedKubeconfig()
	if err != nil {
		return err
	}

	// Apply the namespace.
	managedConfig.SetNamespace(args[0])

	err = managedConfig.Save()
	if err != nil {
		return fmt.Errorf("error updating kubeconfig: %w", err)
	}

	return nil
}
