package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var NamespaceCommand = cobra.Command{
	RunE: NamespaceCommandMain,

	Aliases: []string{
		"names",
		"nsp",
		"ns",
	},

	Use:     "namespace namespace",
	GroupID: "Kubeconfig",

	Short: "Change the current namespace",
	Long: `
	`,
	Example: `
	`,

	Annotations: map[string]string{
		TypeNameAnnotation:       "namespace",
		PluralTypeNameAnnotation: "namespaces",
	},

	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: nil,
}

var NamespaceCommandOptions struct {
}

func init() {
	Command.AddCommand(&NamespaceCommand)
	// CreateListerFor(&NamespaceCommand, NamespaceListItemIter)
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
