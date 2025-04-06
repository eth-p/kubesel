package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ListCommand describes the subcommand for listing information contained
// within kubeconfig files.
//
// Note: The subcommands are generated dynamically as part of program
// initialization. See `listCommandImpl` for the entrypoint of those
// subcommands.
var StatusCommand = cobra.Command{
	RunE: StatusCommandMain,

	Use:     "status",
	GroupID: "Kubesel",

	Short: "Show kubesel status",
	Long: `
	`,

	Args: cobra.NoArgs,
}

func init() {
	Command.AddCommand(&StatusCommand)
}

func StatusCommandMain(cmd *cobra.Command, args []string) error {
	ksel, err := Kubesel()
	if err != nil {
		return err
	}

	session, err := ksel.CurrentSession()
	fmt.Printf("%#v\n%v\n", session, err)
	return nil
}
