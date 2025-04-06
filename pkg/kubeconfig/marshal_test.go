package kubeconfig

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/lithammer/dedent"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	. "github.com/eth-p/kubesel/internal/testutil"
)

var fullDocumentExpected = &Config{
	ApiVersion:     PtrFrom("the-apiVersion"),
	Kind:           PtrFrom("the-kind"),
	CurrentContext: PtrFrom("the-currentContext"),
	Clusters: []NamedCluster{
		{
			Name: PtrFrom("cluster-1"),
			Cluster: &Cluster{
				Server:                   PtrFrom("localhost"),
				TLSServerName:            PtrFrom("localhost.internal"),
				InsecureSkipTLSVerify:    PtrFrom(false),
				CertificateAuthorityData: PtrFrom("abc=="),
				CertificateAuthorityFile: PtrFrom("foo.pem"),
				ProxyURL:                 PtrFrom("http://localhost"),
				DisableCompression:       PtrFrom(true),
			},
		},
		{
			Name: PtrFrom("cluster-2"),
			Cluster: &Cluster{
				Server:                   PtrFrom("cluster.local"),
				CertificateAuthorityFile: PtrFrom("bar.pem"),
				Extensions: []NamedExtension{
					{
						Name: PtrFrom("my-ext"),
						Extension: &Extension{
							ApiVersion: PtrFrom("foo/v1"),
							Kind:       PtrFrom("SomeExtension"),
							Remaining: map[string]any{
								"a-field": "is here",
							},
						},
					},
				},
			},
		},
	},
	AuthInfos: []NamedAuthInfo{
		{
			Name: PtrFrom("my-user"),
			User: &AuthInfo{
				Username: PtrFrom("my-username"),
				Password: PtrFrom("my-password"),
			},
		},
	},
	Contexts: []NamedContext{
		{
			Name: PtrFrom("cluster-1-ctx"),
			Context: &Context{
				Cluster:   PtrFrom("cluster-1"),
				Namespace: PtrFrom("default"),
				User:      PtrFrom("my-user"),
			},
		},
		{
			Name: PtrFrom("cluster-2-ctx"),
			Context: &Context{
				Cluster:   PtrFrom("cluster-2"),
				Namespace: PtrFrom("default"),
				User:      PtrFrom("my-user"),
			},
		},
	},
	Preferences: &Preferences{
		Colors: PtrFrom(true),
	},
}

var fullDocumentYAML = []byte(dedent.Dedent(`
	apiVersion: the-apiVersion
	kind: the-kind
	current-context: the-currentContext
	clusters:
	  - name: cluster-1
	    cluster:
	      server: localhost
	      tls-server-name: localhost.internal
	      insecure-skip-tls-verify: false
	      certificate-authority: foo.pem
	      certificate-authority-data: "abc=="
	      proxy-url: "http://localhost"
	      disable-compression: true
	  - name: cluster-2
	    cluster:
	      server: cluster.local
	      certificate-authority: bar.pem
	      extensions:
	        - name: my-ext
	          extension:
	            apiVersion: foo/v1
	            kind: SomeExtension
	            a-field: is here
	users:
	  - name: my-user
	    user:
	      username: my-username
	      password: my-password
	contexts:
	  - name: cluster-1-ctx
	    context:
	      cluster: cluster-1
	      namespace: default
	      user: my-user
	  - name: cluster-2-ctx
	    context:
	      cluster: cluster-2
	      namespace: default
	      user: my-user
	preferences:
	  colors: true
`))

var fullDocumentJSON = (func() []byte {
	var asMap map[string]any
	err := yaml.Unmarshal([]byte(fullDocumentYAML), &asMap)
	if err != nil {
		panic(err)
	}

	raw, err := json.Marshal(asMap)
	if err != nil {
		panic(err)
	}

	return raw
})()

func TestUnmarshalJSON(t *testing.T) {
	t.Parallel()

	var actual Config
	err := json.Unmarshal(fullDocumentJSON, &actual)
	require.NoError(t, err, "unmarshalling")

	diff := cmp.Diff(*fullDocumentExpected, actual, cmpopts.EquateEmpty())
	require.Empty(t, diff, "--- Expected\n+++ Actual")
}

func TestMarshalJSON(t *testing.T) {
	t.Parallel()

	expected := make(map[string]any)
	actual := make(map[string]any)

	// Create the expected JSON.
	err := json.Unmarshal(fullDocumentJSON, &expected)
	require.NoError(t, err, "unmarshalling expected as unstructured")

	// Get the actual JSON.
	actualEncoded, err := json.Marshal(fullDocumentExpected)
	require.NoError(t, err, "marshalling actual as yaml")

	err = json.Unmarshal(actualEncoded, &actual)
	require.NoError(t, err, "unmarshalling")

	// Compare.
	diff := cmp.Diff(expected, actual, cmpopts.EquateEmpty())
	require.Empty(t, diff, "--- Expected\n+++ Actual")
}

func TestUnmarshalYAML(t *testing.T) {
	t.Parallel()

	var actual Config
	err := yaml.Unmarshal(fullDocumentYAML, &actual)
	require.NoError(t, err, "unmarshalling")

	diff := cmp.Diff(*fullDocumentExpected, actual, cmpopts.EquateEmpty())
	require.Empty(t, diff, "--- Expected\n+++ Actual")
}

func TestMarshalYAML(t *testing.T) {
	t.Parallel()

	expected := make(map[string]any)
	actual := make(map[string]any)

	// Create the expected YAML.
	err := yaml.Unmarshal(fullDocumentYAML, &expected)
	require.NoError(t, err, "unmarshalling expected as unstructured")

	// Get the actual YAML.
	actualEncoded, err := yaml.Marshal(fullDocumentExpected)
	require.NoError(t, err, "marshalling actual as yaml")

	err = yaml.Unmarshal(actualEncoded, &actual)
	require.NoError(t, err, "unmarshalling")

	// Compare.
	diff := cmp.Diff(expected, actual, cmpopts.EquateEmpty())
	require.Empty(t, diff, "--- Expected\n+++ Actual")
}

func BenchmarkUnmarshalJSON(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var config Config
		_ = json.Unmarshal(fullDocumentJSON, &config)
	}
}

func BenchmarkUnmarshalYAML(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var config Config
		_ = yaml.Unmarshal(fullDocumentYAML, &config)
	}
}
