package cobraerr

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseQuotedString(t *testing.T) {
	testcases := map[string]struct {
		Input          string
		ExpectedOk     bool
		ExpectedParsed string
		ExpectedRest   string
	}{
		"Simple string": {
			Input:          "\"hello world\" and rest",
			ExpectedOk:     true,
			ExpectedParsed: "hello world",
			ExpectedRest:   " and rest",
		},
	}

	t.Parallel()
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actualParsed, actualRest, actualOk := parseQuotedString(tc.Input)
			require.Equal(t, tc.ExpectedOk, actualOk, "Was able to parse quoted string")
			require.Equal(t, tc.ExpectedParsed, actualParsed, "The parsed string")
			require.Equal(t, tc.ExpectedRest, actualRest, "The remaining text")
		})
	}
}
