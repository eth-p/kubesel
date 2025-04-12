package cobraerr

import (
	"io"
	"testing"

	"github.com/spf13/cobra"
)

func runCobraCommand(t *testing.T, cmd *cobra.Command, args []string) error {
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	cmd.SetArgs(args)
	return cmd.Execute()
}
