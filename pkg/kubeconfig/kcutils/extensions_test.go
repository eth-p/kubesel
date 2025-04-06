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
