package cobraerr

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestParseInvalidFlagError(t *testing.T) {
	cmd := cobra.Command{
		Use:           "test",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	var shorthandFlag bool
	var boolFlag bool
	var intFlag int
	cmd.Flags().BoolVarP(&shorthandFlag, "long", "l", false, "usage")
	cmd.Flags().BoolVar(&boolFlag, "bool", false, "usage")
	cmd.Flags().IntVar(&intFlag, "int", 0, "usage")

	testcases := map[string]struct {
		Args          []string
		ExpectedError InvalidFlagError
	}{
		"Invalid flag with shorthand": {
			// invalid argument "abc" for "-l, --long" flag: strconv.ParseBool: parsing "abc": invalid syntax
			Args: []string{"--long=abc"},
			ExpectedError: InvalidFlagError{
				Flag:          "long",
				FlagShorthand: "l",
				Value:         "abc",
				Cause:         `"abc" is not a boolean`,
			},
		},
		"Invalid boolean flag": {
			// invalid argument "abc" for "--bool" flag: strconv.ParseBool: parsing "abc": invalid syntax
			Args: []string{"--bool=abc"},
			ExpectedError: InvalidFlagError{
				Flag:  "bool",
				Value: "abc",
				Cause: `"abc" is not a boolean`,
			},
		},
	}

	t.Parallel()
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := runCobraCommand(t, &cmd, tc.Args)
			actualErr, ok := ParseInvalidFlagError(err)
			require.True(t, ok, "Failed to parse the Cobra error string:\n%s", err)

			// Compare.
			diff := cmp.Diff(&tc.ExpectedError, actualErr, cmpopts.EquateEmpty())
			require.Empty(t, diff, "--- Expected\n+++ Actual")
		})
	}
}
