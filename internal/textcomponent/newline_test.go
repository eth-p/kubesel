package printer

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"
)

func TestNewlineComponentPrint(t *testing.T) {
	renderer := NewRenderer()
	Newline.Render(renderer)
	require.Equal(t, "\n", renderer.String())
}

func TestTrimComponentSimplify(t *testing.T) {
	testcases := map[string]struct {
		Input    Trim
		Expected []BasicComponent
	}{
		"No leading, no trailing": {
			Expected: []BasicComponent{
				Newline,
				&Text{Text: "foo"},
				Newline,
			},
			Input: Trim{
				Leading:  false,
				Trailing: false,
				Child: &Sequence{
					Children: []Component{
						Newline,
						&Text{Text: "foo"},
						Newline,
					},
				},
			},
		},
		"Trim leading": {
			Expected: []BasicComponent{
				&Text{Text: "foo"},
				Newline,
			},
			Input: Trim{
				Leading:  true,
				Trailing: false,
				Child: &Sequence{
					Children: []Component{
						Newline,
						&Text{Text: "foo"},
						Newline,
					},
				},
			},
		},
		"Trim trailing": {
			Expected: []BasicComponent{
				Newline,
				&Text{Text: "foo"},
			},
			Input: Trim{
				Leading:  false,
				Trailing: true,
				Child: &Sequence{
					Children: []Component{
						Newline,
						&Text{Text: "foo"},
						Newline,
					},
				},
			},
		},
		"Trim both": {
			Expected: []BasicComponent{
				&Text{Text: "foo"},
			},
			Input: Trim{
				Leading:  true,
				Trailing: true,
				Child: &Sequence{
					Children: []Component{
						Newline,
						&Text{Text: "foo"},
						Newline,
					},
				},
			},
		},
		"Trim leading when only newlines": {
			Expected: []BasicComponent{},
			Input: Trim{
				Leading: true,
				Child: &Sequence{
					Children: []Component{
						Newline,
						Newline,
					},
				},
			},
		},
		"Trim trailing when only newlines": {
			Expected: []BasicComponent{},
			Input: Trim{
				Trailing: true,
				Child: &Sequence{
					Children: []Component{
						Newline,
						Newline,
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
