package cli

import (
	"context"
	"iter"
	"strings"
	"time"

	"github.com/eth-p/kubesel/pkg/kubesel"
	"github.com/spf13/cobra"
)

var namespaceCommand = cobra.Command{
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

	PreRun: tryQuickGC,
}

var NamespaceCommandOptions struct {
}

func init() {
	RootCommand.AddCommand(&namespaceCommand)

	createManagedPropertyCommands(&namespaceCommand, managedProperty[namespaceInfo]{
		PropertyNameSingular: "namespace",
		PropertyNamePlural:   "namespaces",
		GetItemInfos:         namespaceInfoIter,
		GetItemNames:         namespaceNames,
		Switch:               namespaceSwitchImpl,
	})
}

func namespaceSwitchImpl(ksel *kubesel.Kubesel, managedKc *kubesel.ManagedKubeconfig, target string) error {
	managedKc.SetNamespace(target)
	return managedKc.Save()
}

func namespaceNames() ([]string, error) {
	kctl, err := Kubectl()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Get the namespaces using kubectl.
	output, err := kctl.Exec(ctx, []string{"get", "namespace", "--output=name", "--no-headers", "--server-print"})
	if err != nil {
		return nil, err
	}

	// Clean up the returned list.
	namespaces := strings.Split(output, "\n")
	for i, ns := range namespaces {
		namespaces[i] = strings.TrimPrefix(strings.Trim(ns, " \t\r\n"), "namespace/")
	}

	return namespaces, nil
}

type namespaceInfo struct {
	Name *string `yaml:"name" printer:"Name,order=1"`
}

func namespaceInfoIter() (iter.Seq[namespaceInfo], error) {
	namespaces, err := namespaceNames()
	if err != nil {
		return nil, err
	}

	return func(yield func(namespaceInfo) bool) {
		for _, namespace := range namespaces {
			item := namespaceInfo{
				Name: &namespace,
			}

			if !yield(item) {
				return
			}
		}
	}, nil
}
