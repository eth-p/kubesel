package kubeconfig

// Config is the root of the kubeconfig file.
//
// See: [documentation] or [kubernetes source]
//
// [documentation]: https://kubernetes.io/docs/reference/config-api/kubeconfig.v1/#Config
// [kubernetes source]: https://github.com/kubernetes/client-go/blob/2086688a727d00268de695a9c701e52458939701/tools/clientcmd/api/v1/types.go#L28
type Config struct {
	ApiVersion     *string          `yaml:"apiVersion,omitempty"      json:"apiVersion,omitempty"`
	Kind           *string          `yaml:"kind,omitempty"            json:"kind,omitempty"`
	CurrentContext *string          `yaml:"current-context,omitempty" json:"current-context,omitempty" validate:"required"`
	Preferences    *Preferences     `yaml:"preferences,omitempty"     json:"preferences,omitempty"     validate:"required"`
	Clusters       []NamedCluster   `yaml:"clusters,omitempty"        json:"clusters,omitempty"        validate:"required"`
	Contexts       []NamedContext   `yaml:"contexts,omitempty"        json:"contexts,omitempty"        validate:"required"`
	AuthInfos      []NamedAuthInfo  `yaml:"users,omitempty"           json:"users,omitempty"           validate:"required"`
	Extensions     []NamedExtension `yaml:"extensions,omitempty"      json:"extensions,omitempty"`

	// Remaining contains any remaining unknown fields.
	Remaining map[string]any `yaml:",inline,omitempty" json:",inline,omitempty"`
}

// NamedCluster provides a name for a [Cluster].
//
// See: [documentation] or [kubernetes source]
//
// [kubernetes source]: https://github.com/kubernetes/client-go/blob/2086688a727d00268de695a9c701e52458939701/tools/clientcmd/api/v1/types.go#L163
// [documentation]: https://kubernetes.io/docs/reference/config-api/kubeconfig.v1/#NamedCluster
type NamedCluster struct {
	Name    *string  `yaml:"name,omitempty"    json:"name,omitempty"    validate:"required"`
	Cluster *Cluster `yaml:"cluster,omitempty" json:"cluster,omitempty" validate:"required"`

	// Remaining contains any remaining unknown fields.
	Remaining map[string]any `yaml:",inline,omitempty" json:",inline,omitempty"`
}

func (n NamedCluster) key() *string {
	return n.Name
}

// Cluster describes how to communicate with a Kubernetes cluster.
//
// [documentation]: https://kubernetes.io/docs/reference/config-api/kubeconfig.v1/#Cluster
// [kubernetes source]: https://github.com/kubernetes/client-go/blob/2086688a727d00268de695a9c701e52458939701/tools/clientcmd/api/v1/types.go#L63
type Cluster struct {
	Server                   *string          `yaml:"server,omitempty"                     json:"server,omitempty"                     validate:"required"`
	TLSServerName            *string          `yaml:"tls-server-name,omitempty"            json:"tls-server-name,omitempty"`
	InsecureSkipTLSVerify    *bool            `yaml:"insecure-skip-tls-verify,omitempty"   json:"insecure-skip-tls-verify,omitempty"`
	CertificateAuthorityFile *string          `yaml:"certificate-authority,omitempty"      json:"certificate-authority,omitempty"`
	CertificateAuthorityData *string          `yaml:"certificate-authority-data,omitempty" json:"certificate-authority-data,omitempty"`
	ProxyURL                 *string          `yaml:"proxy-url,omitempty"                  json:"proxy-url,omitempty"`
	DisableCompression       *bool            `yaml:"disable-compression,omitempty"        json:"disable-compression,omitempty"`
	Extensions               []NamedExtension `yaml:"extensions,omitempty"                 json:"extensions,omitempty"`

	// Remaining contains any remaining unknown fields.
	Remaining map[string]any `yaml:",inline,omitempty" json:",inline,omitempty"`
}

// NamedContext provides a name for a [Context].
//
// See: [documentation] or [kubernetes source]
//
// [documentation]: https://kubernetes.io/docs/reference/config-api/kubeconfig.v1/#NamedContext
// [kubernetes source]: https://github.com/kubernetes/client-go/blob/2086688a727d00268de695a9c701e52458939701/tools/clientcmd/api/v1/types.go#L171
type NamedContext struct {
	Name    *string  `yaml:"name,omitempty"    json:"name,omitempty"    validate:"required"`
	Context *Context `yaml:"context,omitempty" json:"context,omitempty" validate:"required"`

	// Remaining contains any remaining unknown fields.
	Remaining map[string]any `yaml:",inline,omitempty" json:",inline,omitempty"`
}

func (n NamedContext) key() *string {
	return n.Name
}

// Context specifies a [NamedCluster] to communicate with, a [NamedAuthInfo] to
// authenticate with, and a target namespace to work with.
//
// See: [documentation] or [kubernetes source]
//
// [documentation]: https://kubernetes.io/docs/reference/config-api/kubeconfig.v1/#Context
// [kubernetes source]: https://github.com/kubernetes/client-go/blob/2086688a727d00268de695a9c701e52458939701/tools/clientcmd/api/v1/types.go#L149
type Context struct {
	Cluster    *string          `yaml:"cluster,omitempty"    json:"cluster,omitempty"    validate:"required"`
	User       *string          `yaml:"user,omitempty"       json:"user,omitempty"       validate:"required"`
	Namespace  *string          `yaml:"namespace,omitempty"  json:"namespace,omitempty"`
	Extensions []NamedExtension `yaml:"extensions,omitempty" json:"extensions,omitempty"`

	// Remaining contains any remaining unknown fields.
	Remaining map[string]any `yaml:",inline,omitempty" json:",inline,omitempty"`
}

// NamedAuthInfo provides a name for an [AuthInfo].
//
// See: [documentation] or [kubernetes source]
//
// [documentation]: https://kubernetes.io/docs/reference/config-api/kubeconfig.v1/#NamedAuthInfo
// [kubernetes source]: https://github.com/kubernetes/client-go/blob/2086688a727d00268de695a9c701e52458939701/tools/clientcmd/api/v1/types.go#L179
type NamedAuthInfo struct {
	Name *string   `yaml:"name,omitempty" json:"name,omitempty" validate:"required"`
	User *AuthInfo `yaml:"user,omitempty" json:"user,omitempty" validate:"required"`

	// Remaining contains any remaining unknown fields.
	Remaining map[string]any `yaml:",inline,omitempty" json:",inline,omitempty"`
}

func (n NamedAuthInfo) key() *string {
	return n.Name
}

// AuthInfo describes the authentication method for a user.
//
// See: [documentation] or [kubernetes source]
//
// [documentation]: https://kubernetes.io/docs/reference/config-api/kubeconfig.v1/#AuthInfo
// [kubernetes source]: https://github.com/kubernetes/client-go/blob/2086688a727d00268de695a9c701e52458939701/tools/clientcmd/api/v1/types.go#L100
type AuthInfo struct {
	ClientCertificateFile *string             `yaml:"client-certificate,omitempty"      json:"client-certificate,omitempty"`
	ClientCertificateData *string             `yaml:"client-certificate-data,omitempty" json:"client-certificate-data,omitempty"`
	ClientKeyFile         *string             `yaml:"client-key,omitempty"              json:"client-key,omitempty"`
	ClientKeyData         *string             `yaml:"client-key-data,omitempty"         json:"client-key-data,omitempty"`
	TokenFile             *string             `yaml:"tokenFile,omitempty"               json:"tokenFile,omitempty"`
	Token                 *string             `yaml:"token,omitempty"                   json:"token,omitempty"`
	As                    *string             `yaml:"as,omitempty"                      json:"as,omitempty"`
	AsUID                 *string             `yaml:"as-uid,omitempty"                  json:"as-uid,omitempty"`
	AsGroups              []string            `yaml:"as-groups,omitempty"               json:"as-groups,omitempty"`
	AsUserExtra           map[string][]string `yaml:"as-user-extra,omitempty"           json:"as-user-extra,omitempty"`
	Username              *string             `yaml:"username,omitempty"                json:"username,omitempty"`
	Password              *string             `yaml:"password,omitempty"                json:"password,omitempty"`
	AuthProvider          *AuthProviderConfig `yaml:"auth-provider,omitempty"           json:"auth-provider,omitempty"`
	Exec                  *ExecConfig         `yaml:"exec,omitempty"                    json:"exec,omitempty"`
	Extensions            []NamedExtension    `yaml:"extensions,omitempty"              json:"extensions,omitempty"`

	// Remaining contains any remaining unknown fields.
	Remaining map[string]any `yaml:",inline,omitempty" json:",inline,omitempty"`
}

// AuthProviderConfig specifies an authentication provider and its
// configuration.
//
// See: [documentation] or [kubernetes source]
//
// [documentation]: https://kubernetes.io/docs/reference/config-api/kubeconfig.v1/#AuthProviderConfig
// [kubernetes source]: https://github.com/kubernetes/client-go/blob/2086688a727d00268de695a9c701e52458939701/tools/clientcmd/api/v1/types.go#L195
type AuthProviderConfig struct {
	Name   *string           `yaml:"name,omitempty"   json:"name,omitempty"   validate:"required"`
	Config map[string]string `yaml:"config,omitempty" json:"config,omitempty" validate:"required"`

	// Remaining contains any remaining unknown fields.
	Remaining map[string]any `yaml:",inline,omitempty" json:",inline,omitempty"`
}

// ExecConfig describes an external command used for authentication.
//
// See: [documentation] or [kubernetes source]
//
// [documentation]: https://kubernetes.io/docs/reference/config-api/kubeconfig.v1/#ExecConfig
// [kubernetes source]: https://github.com/kubernetes/client-go/blob/2086688a727d00268de695a9c701e52458939701/tools/clientcmd/api/v1/types.go#L205
type ExecConfig struct {
	Command            *string      `yaml:"command,omitempty"             json:"command,omitempty"             validate:"required"`
	Args               []string     `yaml:"args,omitempty"                json:"args,omitempty"`
	Env                []ExecEnvVar `yaml:"env,omitempty"                 json:"env,omitempty"`
	ApiVersion         *string      `yaml:"apiVersion,omitempty"          json:"apiVersion,omitempty"          validate:"required"`
	InstallHint        *string      `yaml:"installHint,omitempty"         json:"installHint,omitempty"         validate:"required"`
	ProvideClusterInfo *bool        `yaml:"providerClusterInfo,omitempty" json:"providerClusterInfo,omitempty" validate:"required"`
	InteractiveMode    *string      `yaml:"interactiveMode,omitempty"     json:"interactiveMode,omitempty"`

	// Remaining contains any remaining unknown fields.
	Remaining map[string]any `yaml:",inline,omitempty" json:",inline,omitempty"`
}

// ExecEnvVar is an environment variable provided to the external command
// used for authentication.
//
// See: [documentation] or [kubernetes source]
//
// [documentation]: https://kubernetes.io/docs/reference/config-api/kubeconfig.v1/#ExecEnvVar
// [kubernetes source]: https://github.com/kubernetes/client-go/blob/2086688a727d00268de695a9c701e52458939701/tools/clientcmd/api/v1/types.go#L248
type ExecEnvVar struct {
	Name  *string `yaml:"name,omitempty"  json:"name,omitempty"  validate:"required"`
	Value *string `yaml:"value,omitempty" json:"value,omitempty" validate:"required"`

	// Remaining contains any remaining unknown fields.
	Remaining map[string]any `yaml:",inline,omitempty" json:",inline,omitempty"`
}

// NamedExtension provides a name for an [Extension].
//
// See: [documentation] or [kubernetes source]
//
// [documentation]: https://kubernetes.io/docs/reference/config-api/kubeconfig.v1/#NamedExtension
// [kubernetes source]: https://github.com/kubernetes/client-go/blob/2086688a727d00268de695a9c701e52458939701/tools/clientcmd/api/v1/types.go#L187
type NamedExtension struct {
	Name      *string    `yaml:"name,omitempty"      json:"name,omitempty"      validate:"required"`
	Extension *Extension `yaml:"extension,omitempty" json:"extension,omitempty" validate:"required"`

	// Remaining contains any remaining unknown fields.
	Remaining map[string]any `yaml:",inline,omitempty" json:",inline,omitempty"`
}

func (n NamedExtension) key() *string {
	return n.Name
}

// Preferences specifies kubecli client preferences.
//
// See: [documentation] or [kubernetes source]
//
// Deprecated: Replaced by kuberc in Kubernetes 1.33, KEP-3104 ([proposal], [tracker])
//
// [proposal]: https://github.com/kubernetes/enhancements/blob/master/keps/sig-cli/3104-introduce-kuberc/README.md
// [tracker]: https://github.com/kubernetes/enhancements/issues/3104
// [documentation]: https://kubernetes.io/docs/reference/config-api/kubeconfig.v1/#Preferences
// [kubernetes source]: https://github.com/kubernetes/client-go/blob/2086688a727d00268de695a9c701e52458939701/tools/clientcmd/api/v1/types.go#L54
type Preferences struct {
	Colors     *bool            `yaml:"colors,omitempty"     json:"colors,omitempty"`
	Extensions []NamedExtension `yaml:"extensions,omitempty" json:"extensions,omitempty"`

	// Remaining contains any remaining unknown fields.
	Remaining map[string]any `yaml:",inline,omitempty" json:",inline,omitempty"`
}

// Extension is a shim implementation of the Kubernetes [RawExtension].
// This contains an `apiVersion` and `kind`, along with unstructured data.
//
// [RawExtension]: https://github.com/kubernetes/apimachinery/blob/e8a77bd768fd1419e9b3b48a28dd2c6458733a20/pkg/runtime/types.go#L103
type Extension struct {
	ApiVersion *string `yaml:"apiVersion,omitempty" json:"apiVersion,omitempty"`
	Kind       *string `yaml:"kind,omitempty"       json:"kind,omitempty"`

	// Remaining contains any remaining unknown fields.
	Remaining map[string]any `yaml:",inline,omitempty" json:",inline,omitempty"`
}
