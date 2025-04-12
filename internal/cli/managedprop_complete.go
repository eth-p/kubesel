package cli

import (
	"slices"
	"strings"

	"github.com/eth-p/kubesel/internal/fuzzy"
	"github.com/spf13/cobra"
)

// createManagedPropertyCompletionFunc creates a [cobra.CompletionFunc] for
// the provided managed property.
//
// This use the [managedProperty.GetItemNames] function to get the full list
// of items, then filters it down with fuzzy matching.
//
// If [managedProperty.GetItemNames], this returns nil instead.
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
