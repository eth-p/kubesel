package cli

import (
	"iter"
	"reflect"

	"github.com/eth-p/kubesel/pkg/kubesel"
	"github.com/spf13/cobra"
)

// managedProperty describes a kubeconfig property that is managed by kubesel.
type managedProperty[I any] struct {
	PropertyNameSingular string
	PropertyNamePlural   string
	GetItemInfos         itemInfoGenerator[I]
	GetItemNames         func() ([]string, error)

	// Switch changes the active item of this managed property.
	// (e.g. switch to a different cluster or context)
	Switch func(ksel *kubesel.Kubesel, managedKc *kubesel.ManagedKubeconfig, target string) error

	// Set by createManagedPropertyCommands:

	Aliases        []string
	InfoStructType reflect.Type
}

func (p *managedProperty[I]) upcast() *managedProperty[any] {
	return &managedProperty[any]{
		PropertyNameSingular: p.PropertyNameSingular,
		PropertyNamePlural:   p.PropertyNamePlural,
		Aliases:              p.Aliases,
		InfoStructType:       p.InfoStructType,
		GetItemInfos:         p.GetItemInfos.upcast(),
		GetItemNames:         p.GetItemNames,
		Switch:               p.Switch,
	}
}

// createManagedPropertyCommands creates subcommands for common actions
// relating to kubeconfig properties managed by kubesel.
//
// The following subcommands are generated:
//   - `kubesel list <prop>`
//
// The following flags are added to the provided command:
//   - `--list`
//   - `--exact`
func createManagedPropertyCommands[I any](cmd *cobra.Command, prop managedProperty[I]) {
	prop.InfoStructType = reflect.TypeFor[I]()
	if prop.Aliases == nil {
		prop.Aliases = cmd.Aliases
	}

	upProp := prop.upcast() // I -> any
	createManagedPropertySwitchCommand(cmd, upProp)
	createManagedPropertyListSubcommand(cmd, upProp)
}

// createCommandNameAndAliases creates a name and aliases for a [cobra.Command].
//
// The `name` parameter is the command's primary name.
// The aliases exclude the primary name, and are generated from the
// [managedProperty.Aliases], the [managedProperty.PropertyNameSingular],
// and [managedProperty.PropertyNamePlural].
//
// For example, a managed property for contexts would have the singular name
// "context", the plural name "contexts", and the aliases "ctx/contexts". If
// this function was provided the name "ctx", the returned command name and
// aliases would be: "ctx", "context/contexts".
func createCommandNameAndAliases(name string, prop *managedProperty[any]) (string, []string) {
	aliases := make([]string, 0, len(prop.Aliases)+2)

	if prop.PropertyNameSingular != name {
		aliases = append(aliases, prop.PropertyNameSingular)
	}

	if prop.PropertyNamePlural != name {
		aliases = append(aliases, prop.PropertyNamePlural)
	}

	for _, alias := range prop.Aliases {
		if alias != name && alias != prop.PropertyNameSingular && alias != prop.PropertyNamePlural {
			aliases = append(aliases, alias)
		}
	}

	return name, aliases
}

// itemInfoGenerator is a function that creates an iterator over some type.
// This is used to get the items displayed by the `kubesel list` subcommand.
type itemInfoGenerator[I any] func() (iter.Seq[I], error)

func (g *itemInfoGenerator[I]) upcast() itemInfoGenerator[any] {
	return func() (iter.Seq[any], error) {
		realGenerator, err := (*g)()
		if err != nil {
			return nil, err
		}

		return func(yield func(any) bool) {
			for value := range realGenerator {
				if !yield(value) {
					break
				}
			}
		}, nil
	}
}
