package cli

import "github.com/spf13/cobra"

func init() {
	RootCommand.AddCommand(&cobra.Command{
		Use: "info",

		RunE: Info,
	})
}

func Info(cmd *cobra.Command, args []string) error {
	panic("Hi")
}
