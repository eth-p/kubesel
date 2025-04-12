package cli

import (
	"iter"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
)

// managedProperty describes a kubeconfig property that is managed by kubesel.
type managedProperty[I any] struct {
	PropertyNameSingular string
	PropertyNamePlural   string
	ListGenerator        listGenerator[I]

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
