package printer

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"
)

func TestSequenceRender(t *testing.T) {
	testcases := map[string]struct {
		Input    Sequence
		Expected string
	}{
		"No children": {
			Expected: "",
			Input:    Sequence{},
		},
		"With children": {
			Expected: "foobar",
			Input: Sequence{
				Children: []Component{
					&Text{
						Text: "foo",
					},
					&Text{
						Text: "bar",
					},
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

func TestSequenceSimplify(t *testing.T) {
	testcases := map[string]struct {
		Input    Sequence
		Expected []BasicComponent
	}{
		"No children": {
			Expected: []BasicComponent{},
			Input:    Sequence{},
		},
		"With child": {
			Expected: []BasicComponent{
				&Text{
					Text: "bar",
				},
			},
			Input: Sequence{
				Children: []Component{
					&Text{
						Text: "bar",
					},
				},
			},
		},
		"With children that can be simplified": {
			Expected: []BasicComponent{
				&Text{
					Text: "foo",
				},
				Newline,
				&Text{
					Text: "bar",
				},
			},
			Input: Sequence{
				Children: []Component{
					&Text{
						Text: "foo\nbar",
					},
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
