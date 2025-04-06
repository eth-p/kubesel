package kubeconfig

// Clone creates a deep copy.
func (c *Config) Clone() *Config {
	if c == nil {
		return nil
	}

	target := new(Config)
	c.CloneInto(target)

	return target
}

// CloneInto replaces the value of another [Config] with a deep copy of this
// [Config]'s values.
func (c Config) CloneInto(target *Config) {
	target.ApiVersion = shallowCopy(c.ApiVersion)
	target.Kind = shallowCopy(c.Kind)
	target.CurrentContext = shallowCopy(c.CurrentContext)
	target.Preferences = c.Preferences.Clone()
	target.Clusters = cloneSlice(c.Clusters)
	target.Contexts = cloneSlice(c.Contexts)
	target.AuthInfos = cloneSlice(c.AuthInfos)
	target.Extensions = cloneSlice(c.Extensions)
	target.Remaining = cloneMap(c.Remaining)
}

// Clone creates a deep copy.
func (n *NamedCluster) Clone() *NamedCluster {
	if n == nil {
		return nil
	}

	target := new(NamedCluster)
	n.CloneInto(target)

	return target
}

// CloneInto replaces the value of another [NamedCluster] with a deep copy of
// this [NamedCluster]'s values.
func (n NamedCluster) CloneInto(target *NamedCluster) {
	target.Name = shallowCopy(n.Name)
	target.Cluster = n.Cluster.Clone()
	target.Remaining = cloneMap(n.Remaining)
}

// Clone creates a deep copy.
func (c *Cluster) Clone() *Cluster {
	if c == nil {
		return nil
	}

	target := new(Cluster)
	c.CloneInto(target)

	return target
}

// CloneInto replaces the value of another [Cluster] with a deep copy of
// this [Cluster]'s values.
func (c Cluster) CloneInto(target *Cluster) {
	target.Server = shallowCopy(c.Server)
	target.TLSServerName = shallowCopy(c.TLSServerName)
	target.InsecureSkipTLSVerify = shallowCopy(c.InsecureSkipTLSVerify)
	target.CertificateAuthorityFile = shallowCopy(c.CertificateAuthorityFile)
	target.CertificateAuthorityData = shallowCopy(c.CertificateAuthorityData)
	target.ProxyURL = shallowCopy(c.ProxyURL)
	target.DisableCompression = shallowCopy(c.DisableCompression)
	target.Extensions = cloneSlice(c.Extensions)
	target.Remaining = cloneMap(c.Remaining)
}

// Clone creates a deep copy.
func (n *NamedContext) Clone() *NamedContext {
	if n == nil {
		return nil
	}

	target := new(NamedContext)
	n.CloneInto(target)

	return target
}

// CloneInto replaces the value of another [NamedContext] with a deep copy of
// this [NamedContext]'s values.
func (n NamedContext) CloneInto(target *NamedContext) {
	target.Name = shallowCopy(n.Name)
	target.Context = n.Context.Clone()
	target.Remaining = cloneMap(n.Remaining)
}

// Clone creates a deep copy.
func (c *Context) Clone() *Context {
	if c == nil {
		return nil
	}

	target := new(Context)
	c.CloneInto(target)

	return target
}

// CloneInto replaces the value of another [Context] with a deep copy of
// this [Context]'s values.
func (c Context) CloneInto(target *Context) {
	target.Cluster = shallowCopy(c.Cluster)
	target.User = shallowCopy(c.User)
	target.Namespace = shallowCopy(c.Namespace)
	target.Extensions = cloneSlice(c.Extensions)
	target.Remaining = cloneMap(c.Remaining)
}

// Clone creates a deep copy.
func (n *NamedAuthInfo) Clone() *NamedAuthInfo {
	if n == nil {
		return nil
	}

	target := new(NamedAuthInfo)
	n.CloneInto(target)

	return target
}

// CloneInto replaces the value of another [NamedAuthInfo] with a deep copy of
// this [NamedAuthInfo]'s values.
func (n NamedAuthInfo) CloneInto(target *NamedAuthInfo) {
	target.Name = shallowCopy(n.Name)
	target.User = n.User.Clone()
	target.Remaining = cloneMap(n.Remaining)
}

// Clone creates a deep copy.
func (a *AuthInfo) Clone() *AuthInfo {
	if a == nil {
		return nil
	}

	target := new(AuthInfo)
	a.CloneInto(target)

	return target
}

// CloneInto replaces the value of another [AuthInfo] with a deep copy of
// this [AuthInfo]'s values.
func (a AuthInfo) CloneInto(target *AuthInfo) {
	target.ClientCertificateFile = shallowCopy(a.ClientCertificateFile)
	target.ClientCertificateData = shallowCopy(a.ClientCertificateData)
	target.ClientKeyFile = shallowCopy(a.ClientKeyFile)
	target.ClientKeyData = shallowCopy(a.ClientKeyData)
	target.TokenFile = shallowCopy(a.TokenFile)
	target.Token = shallowCopy(a.Token)
	target.As = shallowCopy(a.As)
	target.AsUID = shallowCopy(a.AsUID)
	target.AsGroups = shallowCopySlice(a.AsGroups)
	target.AsUserExtra = cloneMap(a.AsUserExtra)
	target.Username = shallowCopy(a.Username)
	target.Password = shallowCopy(a.Password)
	target.AuthProvider = a.AuthProvider.Clone()
	target.Exec = a.Exec.Clone()
	target.Extensions = cloneSlice(a.Extensions)
	target.Remaining = cloneMap(a.Remaining)
}

// Clone creates a deep copy.
func (a *AuthProviderConfig) Clone() *AuthProviderConfig {
	if a == nil {
		return nil
	}

	target := new(AuthProviderConfig)
	a.CloneInto(target)

	return target
}

// CloneInto replaces the value of another [AuthProviderConfig] with a deep copy
// of this [AuthProviderConfig]'s values.
func (a AuthProviderConfig) CloneInto(target *AuthProviderConfig) {
	target.Name = shallowCopy(a.Name)
	target.Config = cloneMap(a.Config)
	target.Remaining = cloneMap(a.Remaining)
}

// Clone creates a deep copy.
func (e *ExecConfig) Clone() *ExecConfig {
	if e == nil {
		return nil
	}

	target := new(ExecConfig)
	e.CloneInto(target)

	return target
}

// CloneInto replaces the value of another [ExecConfig] with a deep copy
// of this [ExecConfig]'s values.
func (e ExecConfig) CloneInto(target *ExecConfig) {
	target.Command = shallowCopy(e.Command)
	target.Args = shallowCopySlice(e.Args)
	target.Env = cloneSlice(e.Env)
	target.ApiVersion = shallowCopy(e.ApiVersion)
	target.InstallHint = shallowCopy(e.InstallHint)
	target.ProvideClusterInfo = shallowCopy(e.ProvideClusterInfo)
	target.InteractiveMode = shallowCopy(e.InteractiveMode)
	target.Remaining = cloneMap(e.Remaining)
}

// Clone creates a deep copy.
func (e *ExecEnvVar) Clone() *ExecEnvVar {
	if e == nil {
		return nil
	}

	target := new(ExecEnvVar)
	e.CloneInto(target)

	return target
}

// CloneInto replaces the value of another [ExecEnvVar] with a deep copy
// of this [ExecEnvVar]'s values.
func (e ExecEnvVar) CloneInto(target *ExecEnvVar) {
	target.Name = shallowCopy(e.Name)
	target.Value = shallowCopy(e.Value)
	target.Remaining = cloneMap(e.Remaining)
}

// Clone creates a deep copy.
func (n *NamedExtension) Clone() *NamedExtension {
	if n == nil {
		return nil
	}

	target := new(NamedExtension)
	n.CloneInto(target)

	return target
}

// CloneInto replaces the value of another [NamedExtension] with a deep copy
// of this [NamedExtension]'s values.
func (n NamedExtension) CloneInto(target *NamedExtension) {
	target.Name = shallowCopy(n.Name)
	target.Extension = n.Extension.Clone()
	target.Remaining = cloneMap(n.Remaining)
}

// Clone creates a deep copy.
func (p *Preferences) Clone() *Preferences {
	if p == nil {
		return nil
	}

	target := new(Preferences)
	p.CloneInto(target)

	return target
}

// CloneInto replaces the value of another [Preferences] with a deep copy
// of this [Preferences]'s values.
func (p Preferences) CloneInto(target *Preferences) {
	target.Colors = shallowCopy(p.Colors)
	target.Extensions = cloneSlice(p.Extensions)
	target.Remaining = cloneMap(p.Remaining)
}

// Clone creates a deep copy.
func (e *Extension) Clone() *Extension {
	if e == nil {
		return nil
	}

	target := new(Extension)
	e.CloneInto(target)

	return target
}

// CloneInto replaces the value of another [Extension] with a deep copy
// of this [Preferences]'s values.
func (e Extension) CloneInto(target *Extension) {
	target.ApiVersion = shallowCopy(e.ApiVersion)
	target.Kind = shallowCopy(e.Kind)
	target.Remaining = cloneMap(e.Remaining)
}
