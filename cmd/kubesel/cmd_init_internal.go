package main

import (
	"fmt"
	"os"

	"github.com/eth-p/kubesel/pkg/kubesel"
	"github.com/spf13/cobra"
)

var internalInitCommand = cobra.Command{
	RunE: internalInitCommandMain,

	Use:    "__init",
	Hidden: true,

	Short: "Create a new kubesel session file",
	Example: `
		# fish
		set -x KUBECONFIG (kubesel __init --pid=$fish_pid):"$KUBECONFIG"
	`,

	Args: cobra.NoArgs,

	SilenceErrors: true,
	SilenceUsage:  true,
}

var internalInitCommandOptions struct {
	OwnerPID kubesel.SessionOwnerPID
}

func init() {
	Command.AddCommand(&internalInitCommand)

	internalInitCommand.Flags().Int32Var(
		&internalInitCommandOptions.OwnerPID,
		"pid",
		-1,
		"the PID of the session owner",
	)

	internalInitCommand.MarkFlagRequired("pid")
}

func internalInitCommandMain(cmd *cobra.Command, args []string) error {
	ksel, err := Kubesel()
	if err != nil {
		return err
	}

	session, err := ksel.CurrentSession()
	if session != nil {
		os.Exit(2)
	}

	// Create the session and print the resulting file to standard out.
	owner, err := kubesel.SessionOwnerForProcess(internalInitCommandOptions.OwnerPID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "kubesel error creating session: %v\n", err)
		os.Exit(2)
	}

	session, err = ksel.CreateSession(*owner)
	if err != nil {
		fmt.Fprintf(os.Stderr, "kubesel error creating session: %v\n", err)
		os.Exit(2)
	}

	err = session.Save()
	if err != nil {
		fmt.Fprintf(os.Stderr, "kubesel error creating session: %v\n", err)
		os.Exit(2)
	}

	fmt.Fprintf(os.Stdout, "%s\n", session.Path())
	return nil
}
