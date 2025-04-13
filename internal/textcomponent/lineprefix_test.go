package printer

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"
)

func TestLinePrefixComponentPrint(t *testing.T) {
	testcases := map[string]struct {
		Input    LinePrefix
		Expected string
	}{
		"No child": {
			Expected: "",
			Input: LinePrefix{
				Prefix: &Text{
					Text: ">",
				},
			},
		},
		"With child": {
			Expected: ">foo",
			Input: LinePrefix{
				Prefix: &Text{
					Text: ">",
				},
				Child: &Text{
					Text: "foo",
				},
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

func TestLinePrefixComponentSimplify(t *testing.T) {
	testcases := map[string]struct {
		Input    LinePrefix
		Expected []BasicComponent
	}{
		"No child, no prefix": {
			Expected: []BasicComponent{},
			Input:    LinePrefix{},
		},
		"No child": {
			Expected: []BasicComponent{},
			Input: LinePrefix{
				Prefix: &Text{
					Text: ">",
				},
			},
		},
		"No prefix": {
			Expected: []BasicComponent{
				&Text{
					Text: "foo",
				},
			},
			Input: LinePrefix{
				Child: &Text{
					Text: "foo",
				},
			},
		},
		"No newline": {
			Expected: []BasicComponent{
				&Text{
					Text: ">",
				},
				&Text{
					Text: "foo",
				},
			},
			Input: LinePrefix{
				Prefix: &Text{
					Text: ">",
				},
				Child: &Text{
					Text: "foo",
				},
			},
		},
		"With newline": {
			Expected: []BasicComponent{
				&Text{
					Text: ">",
				},
				&Text{
					Text: "foo",
				},
				Newline,
				&Text{
					Text: ">",
				},
				&Text{
					Text: "bar",
				},
			},
			Input: LinePrefix{
				Prefix: &Text{
					Text: ">",
				},
				Child: &Text{
					Text: "foo\nbar",
				},
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
