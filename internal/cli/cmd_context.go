package cli

import (
	"fmt"
	"iter"
	"slices"

	"github.com/eth-p/kubesel/internal/fuzzy"
	"github.com/eth-p/kubesel/pkg/kubeconfig"
	"github.com/eth-p/kubesel/pkg/kubeconfig/kcutils"
	"github.com/eth-p/kubesel/pkg/kubesel"
	"github.com/spf13/cobra"
)

var ContextCommand = cobra.Command{
	RunE: ContextCommandMain,

	Aliases: []string{
		"contexts",
		"con",
		"ctx",
	},

	Use:     "context [name]",
	GroupID: "Kubeconfig",

	Short: "Change to a different cluster, user, and namespace",
	Long: `
		Change the current cluster, user, and namespace to the ones
		specified in a context.

		When selecting a context, you can use its full name as it
		appears in 'kubesel list contexts' or a fuzzy match of
		its name. If no context is specified or if the specified
		name fuzzily matches multiple contexts, a fzf picker will
		be opened.

		If the context does not specify a namespace (or the -n
		flag is set), the current namespace will be kept.
	`,
	Example: `
		kubesel cluster my.cluster.example  # full name
		kubesel cluster myclstr             # fuzzy match
		kubesel cluster                     # fzf picker
	`,

	Annotations: map[string]string{
		TypeNameAnnotation:       "context",
		PluralTypeNameAnnotation: "contexts",
	},

	Args:              cobra.RangeArgs(0, 1),
	ValidArgsFunction: nil,
}

var ContextCommandOptions struct {
	KeepNamespace bool
}

func init() {
	Command.AddCommand(&ContextCommand)
	ContextCommand.PersistentFlags().BoolVarP(
		&ContextCommandOptions.KeepNamespace,
		"keep-namespace", "n",
		false,
		"keep the current namespace",
	)

	CreateListerFor(&ContextCommand, ContextListItemIter)
}

func ContextCommandMain(cmd *cobra.Command, args []string) error {
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

	available := ksel.GetContextNames()
	desired, err := fuzzy.MatchOneOrPick(available, query)
	if err != nil {
		return err
	}

	// Safeguard.
	if !slices.Contains(available, desired) {
		return fmt.Errorf("unknown context: %v", desired)
	}

	// Apply the context.
	kcContext := kcutils.FindContext(desired, ksel.GetMergedKubeconfig())
	if !ContextCommandOptions.KeepNamespace && kcContext.Namespace != nil {
		managedConfig.SetNamespace(*kcContext.Namespace)
	}
	managedConfig.SetClusterName(*kcContext.Cluster)
	managedConfig.SetAuthInfoName(*kcContext.User)

	err = managedConfig.Save()
	if err != nil {
		return fmt.Errorf("error updating kubeconfig: %w", err)
	}

	return nil
}

type ContextListItem struct {
	Name      *string `yaml:"name" printer:"Name,order=0"`
	Cluster   *string `yaml:"cluster" printer:"Cluster,order=1"`
	User      *string `yaml:"user" printer:"User,order=2"`
	Namespace *string `yaml:"namespace" printer:"Namespace,order=3"`
}

func ContextListItemIter() (iter.Seq[ContextListItem], error) {
	ksel, err := Kubesel()
	if err != nil {
		return nil, err
	}

	return func(yield func(ContextListItem) bool) {
		for _, kcNamedContext := range ksel.GetMergedKubeconfig().Contexts {
			if kubesel.IsManagedContext(&kcNamedContext) {
				continue
			}

			kcContext := kcNamedContext.Context
			if kcContext == nil {
				kcContext = &kubeconfig.Context{}
			}

			item := ContextListItem{
				Name:      kcNamedContext.Name,
				Cluster:   kcContext.Cluster,
				User:      kcContext.User,
				Namespace: kcContext.Namespace,
			}

			if !yield(item) {
				return
			}
		}
	}, nil
}
