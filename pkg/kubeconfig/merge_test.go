package kubeconfig

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/lithammer/dedent"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestMerge(t *testing.T) {
	testcases := map[string]struct {
		First    string
		Other    string
		Expected string
	}{
		"First map entry is used": {
			First: `
				apiVersion: first-apiVersion
				kind: first-kind
				current-context: first-currentContext
				preferences:
				  colors: true
			`,
			Other: `
				apiVersion: second-apiVersion
				kind: second-kind
				current-context: second-currentContext
				preferences:
				  colors: false
			`,
			Expected: `
				apiVersion: first-apiVersion
				kind: first-kind
				current-context: first-currentContext
				preferences:
				  colors: true
			`,
		},
		"First non-nil map entry is used": {
			First: `
				apiVersion: first-apiVersion
			`,
			Other: `
				apiVersion: second-apiVersion
				kind: second-kind
				current-context: second-currentContext
				preferences:
				  colors: false
			`,
			Expected: `
				apiVersion: first-apiVersion
				kind: second-kind
				current-context: second-currentContext
				preferences:
				  colors: false
			`,
		},
		"NamedClusters are merged": {
			First: `
				clusters:
				 - name: first-cluster
				   cluster:
				     server: first-cluster-server
			`,
			Other: `
				clusters:
				 - name: second-cluster
				   cluster:
				     server: second-cluster-server
				     proxy-url: foo
			`,
			Expected: `
				clusters:
				 - name: first-cluster
				   cluster:
				     server: first-cluster-server
				 - name: second-cluster
				   cluster:
				     server: second-cluster-server
				     proxy-url: foo
			`,
		},
		"NamedClusters first specified by name is used": {
			First: `
				clusters:
				 - name: first-cluster
				   cluster:
				     server: first-cluster-server
			`,
			Other: `
				clusters:
				 - name: first-cluster
				   cluster:
				     server: second-cluster-server
				     proxy-url: foo
			`,
			Expected: `
				clusters:
				 - name: first-cluster
				   cluster:
				     server: first-cluster-server
			`,
		},
		"NamedContext are merged": {
			First: `
				contexts:
				 - name: first-cluster
				   context:
				     cluster: first-cluster
			`,
			Other: `
				contexts:
				 - name: second-cluster
				   context:
				     cluster: second-cluster
				     user: foo
			`,
			Expected: `
				contexts:
				 - name: first-cluster
				   context:
				     cluster: first-cluster
				 - name: second-cluster
				   context:
				     cluster: second-cluster
				     user: foo
			`,
		},
		"NamedContext first specified by name is used": {
			First: `
				contexts:
				 - name: first-cluster
				   context:
				     cluster: first-cluster
			`,
			Other: `
				contexts:
				 - name: first-cluster
				   context:
				     cluster: second-cluster
				     user: foo
			`,
			Expected: `
				contexts:
				 - name: first-cluster
				   context:
				     cluster: first-cluster
			`,
		},
		"NamedAuthInfos are merged": {
			First: `
				users:
				 - name: first-user
				   user:
				     token: first-user-token
			`,
			Other: `
				users:
				 - name: second-user
				   user:
				     token: second-user-token
				     username: second-user-name
			`,
			Expected: `
				users:
				 - name: first-user
				   user:
				     token: first-user-token
				 - name: second-user
				   user:
				     token: second-user-token
				     username: second-user-name
			`,
		},
		"NamedAuthInfos first specified by name is used": {
			First: `
				users:
				 - name: first-user
				   user:
				     token: first-user-token
			`,
			Other: `
				users:
				 - name: first-user
				   user:
				     token: second-user-token
				     username: second-user-name
			`,
			Expected: `
				users:
				 - name: first-user
				   user:
				     token: first-user-token
			`,
		},
		"Extensions are merged": {
			First: `
				extensions:
				 - name: first-extension
				   extension:
				     apiVersion: first-extension-apiVersion
				     kind: first-extension-kind
			`,
			Other: `
				extensions:
				 - name: second-extension
				   extension:
				     apiVersion: first-extension-apiVersion
				     kind: first-extension-kind
				     spec:
				       foo: bar
			`,
			Expected: `
				extensions:
				 - name: first-extension
				   extension:
				     apiVersion: first-extension-apiVersion
				     kind: first-extension-kind
				 - name: second-extension
				   extension:
				     apiVersion: first-extension-apiVersion
				     kind: first-extension-kind
				     spec:
				       foo: bar
			`,
		},
		"Extensions first specified by name is used": {
			First: `
				extensions:
				 - name: first-extension
				   extension:
				     apiVersion: first-extension-apiVersion
				     kind: first-extension-kind
			`,
			Other: `
				extensions:
				 - name: first-extension
				   extension:
				     apiVersion: first-extension-apiVersion
				     kind: first-extension-kind
				     spec:
				       foo: bar
			`,
			Expected: `
				extensions:
				 - name: first-extension
				   extension:
				     apiVersion: first-extension-apiVersion
				     kind: first-extension-kind
			`,
		},
	}

	t.Parallel()
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Unmarshal the kubeconfigs.
			var first Config
			var other Config
			var expected Config

			err := yaml.Unmarshal([]byte(dedent.Dedent(tc.First)), &first)
			require.NoError(t, err, "unmarshalling first kubeconfig")

			err = yaml.Unmarshal([]byte(dedent.Dedent(tc.Other)), &other)
			require.NoError(t, err, "unmarshalling other kubeconfig")

			err = yaml.Unmarshal([]byte(dedent.Dedent(tc.Expected)), &expected)
			require.NoError(t, err, "unmarshalling expected result")

			// Merge the kubeconfigs.
			actual := MergeConfig(&first, &other)
			diff := cmp.Diff(&expected, actual, cmpopts.EquateEmpty())
			require.Empty(t, diff, "--- Expected\n+++ Actual")
		})
	}
}
