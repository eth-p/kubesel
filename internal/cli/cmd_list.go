package cli

import (
	"fmt"

	"github.com/eth-p/kubesel/internal/printer"
	"github.com/spf13/cobra"
)

// ListCommand describes the subcommand for listing information contained
// within kubeconfig files.
//
// Note: The subcommands are generated dynamically as part of program
// initialization. See `listCommandImpl` for the entrypoint of those
// subcommands.
var ListCommand = cobra.Command{
	RunE:    showHelpIfNoArgs,
	Aliases: []string{"ls"},

	Use:     "list",
	GroupID: "Info",

	Short: "Show from kubeconfig files",
	Long: `
		Display information contained inside your kubeconfig files.
	`,

	Args: cobra.NoArgs,
}

var ListCommandOptions struct {
	OutputFormat OutputFormat
}

func init() {
	Command.AddCommand(&ListCommand)
	ListCommand.PersistentFlags().VarP(
		&ListCommandOptions.OutputFormat,
		"output", "o",
		"The format to print listed info",
	)
}

func listCommandImpl(lister *lister, cmd *cobra.Command, args []string) error {
	itemTyp, err := printer.ItemTypeOf(lister.itemType)
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
	iter, err := lister.itemGenerator()
	if err != nil {
		return fmt.Errorf("cannot list %s: %w", lister.typeNamePlural, err)
	}

	// Print the items.
	for item := range iter {
		printer.Add(item)
	}

	printer.Close()
	return nil
}
