package printer

import (
	"testing"

	"github.com/charmbracelet/x/ansi"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"
)

func TestTextRender(t *testing.T) {
	testcases := map[string]struct {
		Input    Text
		Expected string
	}{
		"No color": {
			Expected: "foo",
			Input: Text{
				Text: "foo",
			},
		},
		"With color": {
			Expected: "\x1B[32mfoo\x1B[0m",
			Input: Text{
				Color: ansi.SGR(ansi.GreenForegroundColorAttr),
				Text:  "foo",
			},
		},
	}

	t.Parallel()
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			renderer := NewRenderer()
			renderer.Render(&tc.Input)
			require.Equal(t, tc.Expected, renderer.String())
		})
	}
}

func TestTextSimplify(t *testing.T) {
	testcases := map[string]struct {
		Input    Text
		Expected []BasicComponent
	}{
		"No text": {
			Expected: []BasicComponent{
				&Text{
					Text: "",
				},
			},
			Input: Text{
				Text: "",
			},
		},
		"No newline": {
			Expected: []BasicComponent{
				&Text{
					Text: "foo",
				},
			},
			Input: Text{
				Text: "foo",
			},
		},
		"One newline": {
			Expected: []BasicComponent{
				&Text{
					Text: "foo",
				},
				Newline,
				&Text{
					Text: "bar",
				},
			},
			Input: Text{
				Text: "foo\nbar",
			},
		},
		"Two newlines": {
			Expected: []BasicComponent{
				&Text{
					Text: "foo",
				},
				Newline,
				&Text{
					Text: "bar",
				},
				Newline,
				&Text{
					Text: "baz",
				},
			},
			Input: Text{
				Text: "foo\nbar\nbaz",
			},
		},
		"Trailing newline": {
			Expected: []BasicComponent{
				&Text{
					Text: "foo",
				},
				Newline,
			},
			Input: Text{
				Text: "foo\n",
			},
		},
	}

	t.Parallel()
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			renderer := NewRenderer()
			actual := tc.Input.Simplify(renderer)
			diff := cmp.Diff(tc.Expected, actual, cmpopts.EquateEmpty())
			require.Empty(t, diff, "--- Expected\n+++ Actual")
		})
	}
}
