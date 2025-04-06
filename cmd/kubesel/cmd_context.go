package main

import (
	"fmt"
	"iter"
	"slices"

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

	Short: "Change the current context",
	Long: `
	`,
	Example: `
	`,

	Annotations: map[string]string{
		TypeNameAnnotation:       "context",
		PluralTypeNameAnnotation: "contexts",
	},

	Args:              cobra.RangeArgs(0, 1),
	ValidArgsFunction: nil,
}

var ContextCommandOptions struct {
}

func init() {
	Command.AddCommand(&ContextCommand)
	CreateListerFor(&ContextCommand, ContextListItemIter)
}

func ContextCommandMain(cmd *cobra.Command, args []string) error {
	ksel, err := Kubesel()
	if err != nil {
		return err
	}

	session, err := ksel.CurrentSession()
	if err != nil {
		return err
	}

	knownContexts := ksel.GetAuthInfoNames()
	desiredContext := ""

	// Select the context.
	if len(args) == 0 {
		// TODO: picker
	} else {
		desiredContext = args[0]
	}

	if !slices.Contains(knownContexts, desiredContext) {
		return fmt.Errorf("unknown context: %v", desiredContext)
	}

	kcContext := kcutils.FindContext(desiredContext, ksel.GetMergedKubeconfig())
	if kcContext.Namespace != nil {
		session.SetNamespace(*kcContext.Namespace)
	}
	session.SetClusterName(*kcContext.Cluster)
	session.SetAuthInfoName(*kcContext.User)

	err = session.Save()
	if err != nil {
		return fmt.Errorf("error saving session: %w", err)
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
