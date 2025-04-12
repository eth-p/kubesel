package cobraerr

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestParseUnknownCommandError(t *testing.T) {
	cmd := cobra.Command{
		Use:                        "test",
		SilenceErrors:              true,
		SilenceUsage:               true,
		SuggestionsMinimumDistance: 2,
	}

	cmd.AddCommand(&cobra.Command{
		Use:           "subcommand",
		SilenceErrors: true,
		SilenceUsage:  true,
	})

	testcases := map[string]struct {
		Args          []string
		ExpectedError UnknownCommandError
	}{
		"Unknown subcommand": {
			// unknown command "unknown" for "test"
			Args: []string{"unknown"},
			ExpectedError: UnknownCommandError{
				Command:       "unknown",
				ParentCommand: cmd.Name(),
			},
		},
		"Unknown subcommand with quote in its name": {
			// unknown command "foo\"bar" for "test"
			Args: []string{"foo\"bar"},
			ExpectedError: UnknownCommandError{
				Command:       "foo\"bar",
				ParentCommand: cmd.Name(),
			},
		},
		"Unknown subcommand with backslash in its name": {
			// unknown command "foo\\bar" for "test"
			Args: []string{"foo\\bar"},
			ExpectedError: UnknownCommandError{
				Command:       "foo\\bar",
				ParentCommand: cmd.Name(),
			},
		},
	}

	t.Parallel()
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := runCobraCommand(t, &cmd, tc.Args)
			actualErr, ok := ParseUnknownCommandError(err)
			require.True(t, ok, "Failed to parse the Cobra error string:\n%s", err)

			// Compare.
			diff := cmp.Diff(&tc.ExpectedError, actualErr, cmpopts.EquateEmpty())
			require.Empty(t, diff, "--- Expected\n+++ Actual")
		})
	}
}
