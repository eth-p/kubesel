package cli

import (
	"iter"

	"github.com/eth-p/kubesel/pkg/kubeconfig"
	"github.com/eth-p/kubesel/pkg/kubesel"
	"github.com/spf13/cobra"
)

var clusterCommand = cobra.Command{
	Aliases: []string{
		"clusters",
		"cl",
	},

	Use:     "cluster [name]",
	GroupID: "Kubeconfig",

	Short: "Change to a different cluster",
	Long: `
		Change to a different Kubernetes cluster in the current shell.

		When selecting a cluster, you can use its full name as it
		appears in 'kubesel list clusters' or a fuzzy match of
		its name. If no cluster is specified or if the specified
		name fuzzily matches multiple clusters, a fzf picker will
		be opened.

		Note: You may need to change the user as well.
	`,
	Example: `
		kubesel cluster my.cluster.example  # full name
		kubesel cluster myclstr             # fuzzy match
		kubesel cluster                     # fzf picker
	`,
}

var ClusterCommandOptions struct {
}

func init() {
	RootCommand.AddCommand(&clusterCommand)
	createManagedPropertyCommands(&clusterCommand, managedProperty[clusterInfo]{
		PropertyNameSingular: "cluster",
		PropertyNamePlural:   "clusters",
		GetItemInfos:         clusterInfoIter,
		GetItemNames:         clusterNames,
		Switch:               clusterSwitchImpl,
	})
}

func clusterSwitchImpl(ksel *kubesel.Kubesel, managedKc *kubesel.ManagedKubeconfig, target string) error {
	managedKc.SetClusterName(target)
	return managedKc.Save()
}

func clusterNames() ([]string, error) {
	kubesel, err := Kubesel()
	if err != nil {
		return nil, err
	}

	return kubesel.GetClusterNames(), nil
}

type clusterInfo struct {
	Name     *string `yaml:"name" printer:"Name,order=1"`
	Server   *string `yaml:"server" printer:"Server,order=2"`
	ProxyURL *string `yaml:"proxy-url" printer:"Proxy URL,order=3,wide"`
}

func clusterInfoIter() (iter.Seq[clusterInfo], error) {
	kubesel, err := Kubesel()
	if err != nil {
		return nil, err
	}

	return func(yield func(clusterInfo) bool) {
		for _, kcNamedCluster := range kubesel.GetMergedKubeconfig().Clusters {
			kcCluster := kcNamedCluster.Cluster
			if kcCluster == nil {
				kcCluster = &kubeconfig.Cluster{}
			}

			item := clusterInfo{
				Name:     kcNamedCluster.Name,
				Server:   kcCluster.Server,
				ProxyURL: kcCluster.ProxyURL,
			}

			if !yield(item) {
				return
			}
		}
	}, nil
}
