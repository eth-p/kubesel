package kcutils

import (
	"testing"

	. "github.com/eth-p/kubesel/internal/testutil"
	"github.com/eth-p/kubesel/pkg/kubeconfig"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"
)

func TestExtensionsFrom(t *testing.T) {
	t.Parallel()

	expected := []kubeconfig.NamedExtension{
		{
			Name: PtrFrom("foo"),
			Extension: &kubeconfig.Extension{
				ApiVersion: PtrFrom("com.example/v1"),
				Kind:       PtrFrom("foo"),
			},
		},
	}

	t.Run("kubeconfig.Config", func(t *testing.T) {
		t.Parallel()
		obj := &kubeconfig.Config{
			Extensions: expected,
		}

		actual := ExtensionsFrom(obj)
		diff := cmp.Diff(expected, actual, cmpopts.EquateEmpty())
		require.Empty(t, diff, "--- Expected\n+++ Actual")
	})

	t.Run("kubeconfig.Cluster", func(t *testing.T) {
		t.Parallel()
		obj := &kubeconfig.Cluster{
			Extensions: expected,
		}

		actual := ExtensionsFrom(obj)
		diff := cmp.Diff(expected, actual, cmpopts.EquateEmpty())
		require.Empty(t, diff, "--- Expected\n+++ Actual")
	})

	t.Run("kubeconfig.Context", func(t *testing.T) {
		t.Parallel()
		obj := &kubeconfig.Context{
			Extensions: expected,
		}

		actual := ExtensionsFrom(obj)
		diff := cmp.Diff(expected, actual, cmpopts.EquateEmpty())
		require.Empty(t, diff, "--- Expected\n+++ Actual")
	})

	t.Run("kubeconfig.AuthInfo", func(t *testing.T) {
		t.Parallel()
		obj := &kubeconfig.AuthInfo{
			Extensions: expected,
		}

		actual := ExtensionsFrom(obj)
		diff := cmp.Diff(expected, actual, cmpopts.EquateEmpty())
		require.Empty(t, diff, "--- Expected\n+++ Actual")
	})

	t.Run("kubeconfig.Preferences", func(t *testing.T) {
		t.Parallel()
		obj := &kubeconfig.Preferences{
			Extensions: expected,
		}

		actual := ExtensionsFrom(obj)
		diff := cmp.Diff(expected, actual, cmpopts.EquateEmpty())
		require.Empty(t, diff, "--- Expected\n+++ Actual")
	})
}

func TestExtensionEncodeDecode(t *testing.T) {
	t.Parallel()

	type SimpleStruct struct {
		Name string
	}

	type SimpleStructWithTag struct {
		Name string `json:"the-name"`
	}

	type EmbeddingStruct struct {
		SimpleStruct
	}

	type InlineStructField struct {
		Inner SimpleStruct `json:",inline"`
	}

	t.Run("Decode works", func(t *testing.T) {
		t.Parallel()

		expected := SimpleStruct{
			Name: "foo",
		}

		var actual SimpleStruct
		err := DecodeExtension(&kubeconfig.Extension{
			Remaining: map[string]any{
				"name": "foo",
			},
		}, &actual)

		require.NoError(t, err)
		diff := cmp.Diff(expected, actual, cmpopts.EquateEmpty())
		require.Empty(t, diff, "--- Expected\n+++ Actual")
	})

	t.Run("Decode is case insensitive", func(t *testing.T) {
		t.Parallel()

		expected := SimpleStruct{
			Name: "foo",
		}

		var actual SimpleStruct
		err := DecodeExtension(&kubeconfig.Extension{
			Remaining: map[string]any{
				"NAME": "foo",
			},
		}, &actual)

		require.NoError(t, err)
		diff := cmp.Diff(expected, actual, cmpopts.EquateEmpty())
		require.Empty(t, diff, "--- Expected\n+++ Actual")
	})

	t.Run("Decode uses json struct tag", func(t *testing.T) {
		t.Parallel()

		expected := SimpleStructWithTag{
			Name: "foo",
		}

		var actual SimpleStructWithTag
		err := DecodeExtension(&kubeconfig.Extension{
			Remaining: map[string]any{
				"the-name": "foo",
			},
		}, &actual)

		require.NoError(t, err)
		diff := cmp.Diff(expected, actual, cmpopts.EquateEmpty())
		require.Empty(t, diff, "--- Expected\n+++ Actual")
	})

	t.Run("Decode inlines embedded structs", func(t *testing.T) {
		t.Parallel()

		expected := EmbeddingStruct{
			SimpleStruct: SimpleStruct{
				Name: "foo",
			},
		}

		var actual EmbeddingStruct
		err := DecodeExtension(&kubeconfig.Extension{
			Remaining: map[string]any{
				"name": "foo",
			},
		}, &actual)

		require.NoError(t, err)
		diff := cmp.Diff(expected, actual, cmpopts.EquateEmpty())
		require.Empty(t, diff, "--- Expected\n+++ Actual")
	})

	t.Run("Decode inlines structs tagged with inline", func(t *testing.T) {
		t.Parallel()

		expected := InlineStructField{
			Inner: SimpleStruct{
				Name: "foo",
			},
		}

		var actual InlineStructField
		err := DecodeExtension(&kubeconfig.Extension{
			Remaining: map[string]any{
				"name": "foo",
			},
		}, &actual)

		require.NoError(t, err)
		diff := cmp.Diff(expected, actual, cmpopts.EquateEmpty())
		require.Empty(t, diff, "--- Expected\n+++ Actual")
	})

	t.Run("Encode works", func(t *testing.T) {
		t.Parallel()

		expected := kubeconfig.Extension{
			Remaining: map[string]any{
				"Name": "foo",
			},
		}

		var actual kubeconfig.Extension
		err := EncodeExtension(&SimpleStruct{
			Name: "foo",
		}, &actual)

		require.NoError(t, err)
		diff := cmp.Diff(expected, actual, cmpopts.EquateEmpty())
		require.Empty(t, diff, "--- Expected\n+++ Actual")
	})

	t.Run("Encode uses specific case", func(t *testing.T) {
		t.Parallel()

		expected := kubeconfig.Extension{
			Remaining: map[string]any{
				"Name": "foo",
			},
		}

		var actual kubeconfig.Extension
		err := EncodeExtension(&SimpleStruct{
			Name: "foo",
		}, &actual)

		require.NoError(t, err)
		diff := cmp.Diff(expected, actual, cmpopts.EquateEmpty())
		require.Empty(t, diff, "--- Expected\n+++ Actual")
	})

	t.Run("Encode uses json struct tag", func(t *testing.T) {
		t.Parallel()

		expected := kubeconfig.Extension{
			Remaining: map[string]any{
				"the-name": "foo",
			},
		}

		var actual kubeconfig.Extension
		err := EncodeExtension(&SimpleStructWithTag{
			Name: "foo",
		}, &actual)

		require.NoError(t, err)
		diff := cmp.Diff(expected, actual, cmpopts.EquateEmpty())
		require.Empty(t, diff, "--- Expected\n+++ Actual")
	})

	t.Run("Encode inlines embedded structs", func(t *testing.T) {
		t.Parallel()

		expected := kubeconfig.Extension{
			Remaining: map[string]any{
				"Name": "foo",
			},
		}

		var actual kubeconfig.Extension
		err := EncodeExtension(&EmbeddingStruct{
			SimpleStruct: SimpleStruct{
				Name: "foo",
			},
		}, &actual)

		require.NoError(t, err)
		diff := cmp.Diff(expected, actual, cmpopts.EquateEmpty())
		require.Empty(t, diff, "--- Expected\n+++ Actual")
	})

	t.Run("Encode inlines structs tagged with inline", func(t *testing.T) {
		t.Parallel()

		expected := kubeconfig.Extension{
			Remaining: map[string]any{
				"Name": "foo",
			},
		}

		var actual kubeconfig.Extension
		err := EncodeExtension(&InlineStructField{
			Inner: SimpleStruct{
				Name: "foo",
			},
		}, &actual)

		require.NoError(t, err)
		diff := cmp.Diff(expected, actual, cmpopts.EquateEmpty())
		require.Empty(t, diff, "--- Expected\n+++ Actual")
	})
}
