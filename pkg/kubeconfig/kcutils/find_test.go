package kcutils

import (
	"testing"

	. "github.com/eth-p/kubesel/internal/testutil"
	"github.com/eth-p/kubesel/pkg/kubeconfig"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"
)

func TestFindContext(t *testing.T) {
	testcases := map[string]struct {
		Kubeconfig *kubeconfig.Config
		Expected   *kubeconfig.Context
		Name       string
	}{
		"Returns nil if no context with name": {
			Name: "foo",
			Kubeconfig: &kubeconfig.Config{
				Contexts: []kubeconfig.NamedContext{
					{
						Name: PtrFrom("localhost"),
						Context: &kubeconfig.Context{
							Cluster: PtrFrom("localhost"),
						},
					},
				},
			},
			Expected: nil,
		},

		"Returns first found with name": {
			Name: "localhost",
			Kubeconfig: &kubeconfig.Config{
				Contexts: []kubeconfig.NamedContext{
					{
						Name: PtrFrom("localhost"),
						Context: &kubeconfig.Context{
							Cluster: PtrFrom("localhost"),
						},
					},
					{
						Name: PtrFrom("localhost"),
						Context: &kubeconfig.Context{
							Cluster: PtrFrom("not-localhost"),
						},
					},
				},
			},
			Expected: &kubeconfig.Context{
				Cluster: PtrFrom("localhost"),
			},
		},
	}

	t.Parallel()
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actual := FindContext(tc.Name, tc.Kubeconfig)

			// Compare.
			diff := cmp.Diff(tc.Expected, actual, cmpopts.EquateEmpty())
			require.Empty(t, diff, "--- Expected\n+++ Actual")
		})
	}
}

func TestFindCluster(t *testing.T) {
	testcases := map[string]struct {
		Kubeconfig *kubeconfig.Config
		Expected   *kubeconfig.Cluster
		Name       string
	}{
		"Returns nil if no cluster with name": {
			Name: "foo",
			Kubeconfig: &kubeconfig.Config{
				Clusters: []kubeconfig.NamedCluster{
					{
						Name: PtrFrom("localhost"),
						Cluster: &kubeconfig.Cluster{
							Server: PtrFrom("localhost"),
						},
					},
				},
			},
			Expected: nil,
		},

		"Returns first found with name": {
			Name: "localhost",
			Kubeconfig: &kubeconfig.Config{
				Clusters: []kubeconfig.NamedCluster{
					{
						Name: PtrFrom("localhost"),
						Cluster: &kubeconfig.Cluster{
							Server: PtrFrom("localhost"),
						},
					},
					{
						Name: PtrFrom("localhost"),
						Cluster: &kubeconfig.Cluster{
							Server: PtrFrom("not-localhost"),
						},
					},
				},
			},
			Expected: &kubeconfig.Cluster{
				Server: PtrFrom("localhost"),
			},
		},
	}

	t.Parallel()
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actual := FindCluster(tc.Name, tc.Kubeconfig)

			// Compare.
			diff := cmp.Diff(tc.Expected, actual, cmpopts.EquateEmpty())
			require.Empty(t, diff, "--- Expected\n+++ Actual")
		})
	}
}

func TestFindAuthInfo(t *testing.T) {
	testcases := map[string]struct {
		Kubeconfig *kubeconfig.Config
		Expected   *kubeconfig.AuthInfo
		Name       string
	}{
		"Returns nil if no user with name": {
			Name: "foo",
			Kubeconfig: &kubeconfig.Config{
				Clusters: []kubeconfig.NamedCluster{
					{
						Name: PtrFrom("localhost"),
						Cluster: &kubeconfig.Cluster{
							Server: PtrFrom("localhost"),
						},
					},
				},
			},
			Expected: nil,
		},

		"Returns first found with name": {
			Name: "localhost",
			Kubeconfig: &kubeconfig.Config{
				AuthInfos: []kubeconfig.NamedAuthInfo{
					{
						Name: PtrFrom("localhost"),
						User: &kubeconfig.AuthInfo{
							Username: PtrFrom("localhost"),
						},
					},
					{
						Name: PtrFrom("localhost"),
						User: &kubeconfig.AuthInfo{
							Username: PtrFrom("not-localhost"),
						},
					},
				},
			},
			Expected: &kubeconfig.AuthInfo{
				Username: PtrFrom("localhost"),
			},
		},
	}

	t.Parallel()
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actual := FindAuthInfo(tc.Name, tc.Kubeconfig)

			// Compare.
			diff := cmp.Diff(tc.Expected, actual, cmpopts.EquateEmpty())
			require.Empty(t, diff, "--- Expected\n+++ Actual")
		})
	}
}

func TestFindExtension(t *testing.T) {
	testcases := map[string]struct {
		Extensions []kubeconfig.NamedExtension
		Expected   *kubeconfig.Extension
		Name       string
	}{
		"Returns nil if no extension with name": {
			Name: "foo",
			Extensions: []kubeconfig.NamedExtension{
				{
					Name: PtrFrom("test"),
					Extension: &kubeconfig.Extension{
						ApiVersion: PtrFrom("com.example/v1"),
						Kind:       PtrFrom("Foo"),
					},
				},
			},
			Expected: nil,
		},

		"Returns first found with name": {
			Name: "test",
			Extensions: []kubeconfig.NamedExtension{
				{
					Name: PtrFrom("test"),
					Extension: &kubeconfig.Extension{
						ApiVersion: PtrFrom("com.example/v1"),
						Kind:       PtrFrom("Foo"),
					},
				},
				{
					Name: PtrFrom("test"),
					Extension: &kubeconfig.Extension{
						ApiVersion: PtrFrom("com.example/v1"),
						Kind:       PtrFrom("Bar"),
					},
				},
			},
			Expected: &kubeconfig.Extension{
				ApiVersion: PtrFrom("com.example/v1"),
				Kind:       PtrFrom("Foo"),
			},
		},
	}

	t.Parallel()
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actual := FindExtension(tc.Name, tc.Extensions)

			// Compare.
			diff := cmp.Diff(tc.Expected, actual, cmpopts.EquateEmpty())
			require.Empty(t, diff, "--- Expected\n+++ Actual")
		})
	}
}

func TestFindExtensionsByKind(t *testing.T) {
	testcases := map[string]struct {
		Extensions []kubeconfig.NamedExtension
		Expected   []*kubeconfig.Extension
		ApiVersion string
		Kind       string
	}{
		"Returns nil if no matching extension": {
			ApiVersion: "com.other",
			Kind:       "Bar",
			Extensions: []kubeconfig.NamedExtension{
				{
					Name: PtrFrom("test"),
					Extension: &kubeconfig.Extension{
						ApiVersion: PtrFrom("com.example/v1"),
						Kind:       PtrFrom("Foo"),
					},
				},
			},
			Expected: nil,
		},

		"Returns all found with name": {
			ApiVersion: "com.example/v1",
			Kind:       "Foo",
			Extensions: []kubeconfig.NamedExtension{
				{
					Name: PtrFrom("test"),
					Extension: &kubeconfig.Extension{
						ApiVersion: PtrFrom("com.example/v1"),
						Kind:       PtrFrom("Foo"),
						Remaining: map[string]any{
							"position": 1,
						},
					},
				},
				{
					Name: PtrFrom("test"),
					Extension: &kubeconfig.Extension{
						ApiVersion: PtrFrom("com.example/v1"),
						Kind:       PtrFrom("Foo"),
						Remaining: map[string]any{
							"position": 2,
						},
					},
				},
			},
			Expected: []*kubeconfig.Extension{
				{
					ApiVersion: PtrFrom("com.example/v1"),
					Kind:       PtrFrom("Foo"),
					Remaining: map[string]any{
						"position": 1,
					},
				},
				{
					ApiVersion: PtrFrom("com.example/v1"),
					Kind:       PtrFrom("Foo"),
					Remaining: map[string]any{
						"position": 2,
					},
				},
			},
		},
	}

	t.Parallel()
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actual := FindExtensionsByKind(tc.ApiVersion, tc.Kind, tc.Extensions)

			// Compare.
			diff := cmp.Diff(tc.Expected, actual, cmpopts.EquateEmpty())
			require.Empty(t, diff, "--- Expected\n+++ Actual")
		})
	}
}
