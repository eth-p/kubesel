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
	cmd.ValidArgsFunction = createManagedPropertyCompletionFunc(prop)

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

func createManagedPropertyCompletionFunc(prop *managedProperty[any]) cobra.CompletionFunc {
	if prop.GetItemNames == nil {
		return nil
	}

	return func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
		items, err := getCompletionItemsFromNames(prop)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		// Sort the items.
		slices.SortFunc(items, func(a, b completionItem) int {
			return strings.Compare(a.name, b.name)
		})

		// Filter out the returned options by fuzzy matching.
		if toComplete != "" {
			items = filterCompletionItemsByFuzzy(items, toComplete)
		}

		// Return the matching names.
		cobraComps := make([]cobra.Completion, len(items))
		for i, item := range items {
			cobraComps[i] = item.asCobraCompletion()
		}

		return cobraComps, cobra.ShellCompDirectiveNoFileComp
	}
}

func filterCompletionItemsByFuzzy(items []completionItem, query string) []completionItem {
	matches := fuzzy.MatchesFunc(items, query, func(ci completionItem) string {
		return ci.name
	})

	// Sort the matches by score.
	slices.SortFunc(matches, func(a, b fuzzy.MatchResult[completionItem]) int {
		if a.Score != b.Score {
			return b.Score - a.Score
		}

		return strings.Compare(a.Item.name, b.Item.name)
	})

	// Return the original completion items.
	result := make([]completionItem, len(matches))
	for i, match := range matches {
		result[i] = match.Item
	}

	return result
}

type completionItem struct {
	name string
	// TODO: descriptions
}

func (c completionItem) asCobraCompletion() string {
	return c.name
}

// getCompletionItemsFromNames returns completion items by using the
// [managedProperty.GetItemNames] function to fetch the list of valid names.
func getCompletionItemsFromNames(prop *managedProperty[any]) ([]completionItem, error) {
	names, err := prop.GetItemNames()
	if err != nil {
		return nil, err
	}

	completions := make([]completionItem, len(names))
	for i, name := range names {
		completions[i].name = name
	}

	return completions, nil
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
