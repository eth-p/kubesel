package cobraerr

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestParseUnknownFlagError(t *testing.T) {
	cmd := cobra.Command{
		Use:           "test",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	var boolFlag bool
	cmd.Flags().BoolVarP(&boolFlag, "bool", "b", false, "usage")

	testcases := map[string]struct {
		Args          []string
		ExpectedError UnknownFlagError
	}{
		"Unknown shorthand flag": {
			// unknown shorthand flag: 'x' in -x
			Args: []string{"-x"},
			ExpectedError: UnknownFlagError{
				IsShorthandFlag: true,
				Flag:            "x",
				FlagSet:         "x",
			},
		},
		"Unknown shorthand flag in set of multiple shorthand flags": {
			// unknown shorthand flag: 'x' in -axy
			Args: []string{"-axy"},
			ExpectedError: UnknownFlagError{
				IsShorthandFlag: true,
				Flag:            "a",
				FlagSet:         "axy",
			},
		},
		"Unknown shorthand flag that is apostrophe": {
			// unknown shorthand flag: '\'' in -'
			Args: []string{"-'"},
			ExpectedError: UnknownFlagError{
				IsShorthandFlag: true,
				Flag:            "'",
				FlagSet:         "'",
			},
		},
		"Unknown shorthand flag that is backslash": {
			// unknown shorthand flag: '\\' in -\
			Args: []string{"-\\"},
			ExpectedError: UnknownFlagError{
				IsShorthandFlag: true,
				Flag:            "\\",
				FlagSet:         "\\",
			},
		},
		"Unknown long flag": {
			// unknown flag: --x
			Args: []string{"--x"},
			ExpectedError: UnknownFlagError{
				IsShorthandFlag: false,
				Flag:            "x",
				FlagSet:         "",
			},
		},
		"Unknown long flag with apostrophe": {
			// unknown flag: --x'
			Args: []string{"--x'"},
			ExpectedError: UnknownFlagError{
				IsShorthandFlag: false,
				Flag:            "x'",
				FlagSet:         "",
			},
		},
		"Unknown long flag with quote": {
			// unknown flag: --""
			Args: []string{"--\""},
			ExpectedError: UnknownFlagError{
				IsShorthandFlag: false,
				Flag:            "\"",
				FlagSet:         "",
			},
		},
		"Unknown long flag with backslash": {
			// unknown flag: --\
			Args: []string{"--\\"},
			ExpectedError: UnknownFlagError{
				IsShorthandFlag: false,
				Flag:            "\\",
				FlagSet:         "",
			},
		},
	}

	t.Parallel()
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := runCobraCommand(t, &cmd, tc.Args)
			actualErr, ok := ParseUnknownFlagError(err)
			require.True(t, ok, "Failed to parse the Cobra error string:\n%s", err)

			// Compare.
			diff := cmp.Diff(&tc.ExpectedError, actualErr, cmpopts.EquateEmpty())
			require.Empty(t, diff, "--- Expected\n+++ Actual")
		})
	}
}
