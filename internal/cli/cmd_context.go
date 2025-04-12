package cli

import (
	"iter"

	"github.com/eth-p/kubesel/pkg/kubeconfig"
	"github.com/eth-p/kubesel/pkg/kubeconfig/kcutils"
	"github.com/eth-p/kubesel/pkg/kubesel"
	"github.com/spf13/cobra"
)

var contextCommand = cobra.Command{
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
}

var ContextCommandOptions struct {
	KeepNamespace bool
}

func init() {
	RootCommand.AddCommand(&contextCommand)
	contextCommand.PersistentFlags().BoolVarP(
		&ContextCommandOptions.KeepNamespace,
		"keep-namespace", "n",
		false,
		"keep the current namespace",
	)

	createManagedPropertyCommands(&contextCommand, managedProperty[ContextListItem]{
		PropertyNameSingular: "context",
		PropertyNamePlural:   "contexts",
		ListGenerator:        ContextListItemIter,
		GetItemNames:         contextNames,
		Switch:               contextSwitchImpl,
	})
}

func contextSwitchImpl(ksel *kubesel.Kubesel, managedKc *kubesel.ManagedKubeconfig, target string) error {
	kcContext := kcutils.FindContext(target, ksel.GetMergedKubeconfig())

	managedKc.SetClusterName(*kcContext.Cluster)
	managedKc.SetAuthInfoName(*kcContext.User)

	if !ContextCommandOptions.KeepNamespace && kcContext.Namespace != nil {
		managedKc.SetNamespace(*kcContext.Namespace)
	}

	return managedKc.Save()
}

func contextNames() ([]string, error) {
	kubesel, err := Kubesel()
	if err != nil {
		return nil, err
	}

	return kubesel.GetAuthInfoNames(), nil
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
