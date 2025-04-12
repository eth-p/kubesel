package cli

import (
	"fmt"
	"iter"
	"reflect"
	"slices"
	"strings"

	"github.com/eth-p/kubesel/internal/fuzzy"
	"github.com/eth-p/kubesel/pkg/kubesel"
	"github.com/spf13/cobra"
)

// managedProperty describes a kubeconfig property that is managed by kubesel.
type managedProperty[I any] struct {
	PropertyNameSingular string
	PropertyNamePlural   string
	ListGenerator        listGenerator[I]
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
		ListGenerator:        p.ListGenerator.upcast(),
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
// And `--list` is added as a flag to the provided command.
func createManagedPropertyCommands[I any](cmd *cobra.Command, prop managedProperty[I]) {
	prop.InfoStructType = reflect.TypeFor[I]()
	if prop.Aliases == nil {
		prop.Aliases = cmd.Aliases
	}

	upProp := prop.upcast() // I -> any
	createManagedPropertySwitchCommand(cmd, upProp)
	createManagedPropertyListSubcommand(cmd, upProp)
}

func createManagedPropertyListSubcommand(cmd *cobra.Command, prop *managedProperty[any]) {
	pluralName := prop.PropertyNamePlural
	cmdName, cmdAliases := createCommandNameAndAliases(pluralName, prop)
	subcmd := &cobra.Command{
		Use:     cmdName,
		Aliases: cmdAliases,

		Short: "List available " + pluralName,
		Long:  strings.ReplaceAll(listCommand.Long, "info", pluralName),

		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return generatedListCommandMain(prop, cmd, args)
		},
	}

	listCommand.AddCommand(subcmd)

	// Add the `--list` flag to the original command and set it as hidden.
	var printList bool
	cmd.Flags().BoolVar(&printList, listFlagName, false, "")
	flag := cmd.Flag(listFlagName)
	flag.Hidden = true
	flag.Usage = "Print the available " + prop.PropertyNamePlural

	// Wrap the target command's `RunE` function.
	realRunE := cmd.RunE
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if printList {
			UseListOutput("", &ListCommandOptions.OutputFormat)
			return subcmd.RunE(subcmd, []string{})
		}

		return realRunE(cmd, args)
	}
}

// createManagedPropertySwitchCommand updates the provided cobra command's
// run function to fuzzy-match/fuzzy-pick an item of the managed property type.
//
// Once the item is picked, the Switch function is called to change the
// property inside the managed kubeconfig.
func createManagedPropertySwitchCommand(cmd *cobra.Command, prop *managedProperty[any]) {
	if prop.Switch == nil {
		return
	}

	cmd.Args = cobra.RangeArgs(0, 1)
	cmd.ValidArgsFunction = nil
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ksel, err := Kubesel()
		if err != nil {
			return err
		}

		managedKc, err := ksel.GetManagedKubeconfig()
		if err != nil {
			return err
		}

		// Get the available item names.
		available, err := prop.GetItemNames()
		if err != nil {
			return err
		}

		// Fuzzy match/pick based on the query (or lack thereof)
		query := ""
		if len(args) > 0 {
			query = args[0]
		}

		desired, err := fuzzy.MatchOneOrPick(available, query)
		if err != nil {
			return err
		}

		// Safeguard.
		if !slices.Contains(available, desired) {
			return fmt.Errorf("unknown %s: %v", prop.PropertyNamePlural, desired)
		}

		// Switch.
		return prop.Switch(ksel, managedKc, desired)
	}
}

func createCommandNameAndAliases(name string, prop *managedProperty[any]) (string, []string) {
	aliases := make([]string, 0, len(prop.Aliases)+2)

	if prop.PropertyNameSingular != name {
		aliases = append(aliases, prop.PropertyNameSingular)
	}

	if prop.PropertyNamePlural != name {
		aliases = append(aliases, prop.PropertyNamePlural)
	}

	for _, alias := range prop.Aliases {
		if alias != name {
			aliases = append(aliases, alias)
		}
	}

	return name, aliases
}

// listGenerator is a function that creates an iterator over some type.
// This is used to get the items displayed by the `kubesel list` subcommand.
type listGenerator[I any] func() (iter.Seq[I], error)

func (g *listGenerator[I]) upcast() listGenerator[any] {
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
