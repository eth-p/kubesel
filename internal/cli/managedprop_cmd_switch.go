package cli

import (
	"fmt"
	"slices"

	"github.com/eth-p/kubesel/internal/fuzzy"
	"github.com/spf13/cobra"
)

// createManagedPropertySwitchCommand updates the provided cobra command's
// run function to fuzzy-match/fuzzy-pick an item of the managed property type.
//
// Once the item is picked, the [managedProperty.Switch] function is called to
// update the managed kubeconfig with the selected item.
func createManagedPropertySwitchCommand(cmd *cobra.Command, prop *managedProperty[any]) {
	if prop.Switch == nil {
		return
	}

	// Flags.
	var mustExactMatch bool
	cmd.Flags().BoolVarP(&mustExactMatch,
		"exact", "e",
		false,
		prop.PropertyNameSingular+" must be exact match",
	)

	// Command.
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
		var desired string
		query := ""
		if len(args) > 0 {
			query = args[0]
		}

		if mustExactMatch {
			desired = query
		} else {
			desired, err = fuzzy.MatchOneOrPick(available, query)
			if err != nil {
				return err
			}
		}

		if desired == "" {
			return fmt.Errorf("no %s specified", prop.PropertyNameSingular)
		}

		// Safeguard.
		if !slices.Contains(available, desired) {
			return fmt.Errorf("unknown %s: %v", prop.PropertyNamePlural, desired)
		}

		// Switch.
		return prop.Switch(ksel, managedKc, desired)
	}
}
