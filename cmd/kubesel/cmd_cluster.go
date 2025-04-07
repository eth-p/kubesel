package main

import (
	"fmt"
	"iter"
	"slices"

	"github.com/eth-p/kubesel/internal/fuzzy"
	"github.com/eth-p/kubesel/pkg/kubeconfig"
	"github.com/spf13/cobra"
)

var ClusterCommand = cobra.Command{
	RunE: ClusterCommandMain,
	Aliases: []string{
		"clusters",
		"cl",
	},

	Use:     "cluster [name]",
	GroupID: "Kubeconfig",

	Short:   "Change the current cluster",
	Long:    "",
	Example: "",

	Annotations: map[string]string{
		TypeNameAnnotation:       "cluster",
		PluralTypeNameAnnotation: "clusters",
	},

	Args:              cobra.RangeArgs(0, 1),
	ValidArgsFunction: nil,
}

var ClusterCommandOptions struct {
}

func init() {
	Command.AddCommand(&ClusterCommand)
	CreateListerFor(&ClusterCommand, ClusterListItemIter)
}

func ClusterCommandMain(cmd *cobra.Command, args []string) error {
	ksel, err := Kubesel()
	if err != nil {
		return err
	}

	managedConfig, err := ksel.GetManagedKubeconfig()
	if err != nil {
		return err
	}

	query := ""
	if len(args) > 0 {
		query = args[0]
	}

	available := ksel.GetClusterNames()
	desired, err := fuzzy.MatchOneOrPick(available, query)
	if err != nil {
		return err
	}

	// Safeguard.
	if !slices.Contains(available, desired) {
		return fmt.Errorf("unknown cluster: %v", desired)
	}

	managedConfig.SetClusterName(desired)
	err = managedConfig.Save()
	if err != nil {
		return fmt.Errorf("error updating kubeconfig: %w", err)
	}

	return nil
}

type ClusterListItem struct {
	Name     *string `yaml:"name" printer:"Name,order=1"`
	Server   *string `yaml:"server" printer:"Server,order=2"`
	ProxyURL *string `yaml:"proxy-url" printer:"Proxy URL,order=3,wide"`
}

func ClusterListItemIter() (iter.Seq[ClusterListItem], error) {
	kubesel, err := Kubesel()
	if err != nil {
		return nil, err
	}

	return func(yield func(ClusterListItem) bool) {
		for _, cluster := range kubesel.GetMergedKubeconfig().Clusters {
			clusterInfo := cluster.Cluster
			if clusterInfo == nil {
				clusterInfo = &kubeconfig.Cluster{}
			}

			item := ClusterListItem{
				Name:     cluster.Name,
				Server:   clusterInfo.Server,
				ProxyURL: clusterInfo.ProxyURL,
			}

			if !yield(item) {
				return
			}
		}
	}, nil
}
