package cli

import (
	"fmt"

	"github.com/eth-p/kubesel/internal/printer"
	"github.com/spf13/cobra"
)

// listCommand describes the subcommand for listing information contained
// within kubeconfig files.
//
// Note: The subcommands are generated dynamically as part of program
// initialization. See `listCommandImpl` for the entrypoint of those
// subcommands.
var listCommand = cobra.Command{
	RunE:    showHelpIfNoArgs,
	Aliases: []string{"ls"},

	Use:     "list",
	GroupID: "Info",

	Short: "List info from kubeconfig files",
	Long: `
		List info from your kubeconfig files.

		Available output formats are:
		  list   (print only the names as a list)
		  table  (print as a table)
		  cols   (print as columns)

		With the table and column formats, the printed columns and
		their order can be changed by by appending '=COL1,COL2' to
		to the output format (e.g. '--output table=name,cluster').
	`,

	Args: cobra.NoArgs,
}

var ListCommandOptions struct {
	OutputFormat OutputFormat
}

func init() {
	RootCommand.AddCommand(&listCommand)
	listCommand.PersistentFlags().VarP(
		&ListCommandOptions.OutputFormat,
		"output", "o",
		"output format",
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
