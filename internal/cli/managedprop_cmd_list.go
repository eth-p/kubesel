package cli

import (
	"fmt"
	"strings"

	"github.com/eth-p/kubesel/internal/printer"
	"github.com/spf13/cobra"
)

// createManagedPropertyListSubcommand creates a [cobra.Command] and attaches
// it to the top-level [listCommand] (`kubesel list`) command.
//
// This uses the [managedProperty.GetItemInfos] function to create an iterator
// of structs which are given to a struct printer. The struct printer used
// depends on the value of the `--output` flag.
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
			return managedPropertyListSubcommandMain(prop, cmd, args)
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

func managedPropertyListSubcommandMain(prop *managedProperty[any], cmd *cobra.Command, args []string) error {
	itemTyp, err := printer.ItemTypeOf(prop.InfoStructType)
	if err != nil {
		return err
	}

	// Create the item printer.
	ListCommandOptions.OutputFormat.DefaultIfUnset()
	printer, err := ListCommandOptions.OutputFormat.newPrinter(
		*itemTyp,
		cmd.OutOrStdout(),
	)

	if err != nil {
		return err
	}

	// Start iterating the items.
	iter, err := prop.GetItemInfos()
	if err != nil {
		return fmt.Errorf("cannot list %s: %w", prop.PropertyNamePlural, err)
	}

	// Print the items.
	for item := range iter {
		printer.Add(item)
	}

	printer.Close()
	return nil
}
