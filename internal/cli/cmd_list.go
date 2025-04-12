package cli

import (
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
