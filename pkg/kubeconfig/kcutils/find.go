package kcutils

import "github.com/eth-p/kubesel/pkg/kubeconfig"

// FindContext returns a pointer to the [kubeconfig.Context] with the given
// name inside the provided [kubeconfig.Config] struct.
//
// If the configuration does not contain the context, nil will be returned.
func FindContext(name string, config *kubeconfig.Config) *kubeconfig.Context {
	for _, namedContext := range config.Contexts {
		if namedContext.Name != nil && *namedContext.Name == name {
			return namedContext.Context
		}
	}

	return nil
}

// FindContext returns a pointer to the [kubeconfig.Cluster] with the given
// name inside the provided [kubeconfig.Config] struct.
//
// If the configuration does not contain the cluster, nil will be returned.
func FindCluster(name string, config *kubeconfig.Config) *kubeconfig.Cluster {
	for _, namedCluster := range config.Clusters {
		if namedCluster.Name != nil && *namedCluster.Name == name {
			return namedCluster.Cluster
		}
	}

	return nil
}

// FindAuthInfo returns a pointer to the [kubeconfig.AuthInfo] with the given
// name inside the provided [kubeconfig.Config] struct.
//
// If the configuration does not contain the cluster, nil will be returned.
func FindAuthInfo(name string, config *kubeconfig.Config) *kubeconfig.AuthInfo {
	for _, namedAuthInfo := range config.AuthInfos {
		if namedAuthInfo.Name != nil && *namedAuthInfo.Name == name {
			return namedAuthInfo.User
		}
	}

	return nil
}

// FindExtension returns a pointer to the [kubeconfig.Extension] with the
// given name inside the provided list of [kubeconfig.NamedExtension]s.
//
// If the list does not contain the cluster, nil will be returned.
func FindExtension(name string, extensions []kubeconfig.NamedExtension) *kubeconfig.Extension {
	for _, namedExtension := range extensions {
		if namedExtension.Name != nil && *namedExtension.Name == name {
			return namedExtension.Extension
		}
	}

	return nil
}

// FindExtension returns a pointer to the [kubeconfig.Extension] with the given
// apiVersion and kind inside the provided list of [kubeconfig.NamedExtension]s.
//
// If the list does not contain the cluster, nil will be returned.
func FindExtensionsByKind(apiVersion, kind string, extensions []kubeconfig.NamedExtension) []*kubeconfig.Extension {
	var found []*kubeconfig.Extension
	for _, namedExtension := range extensions {
		if namedExtension.Name == nil || namedExtension.Extension == nil {
			continue
		}

		ext := namedExtension.Extension
		if ext.ApiVersion == nil || ext.Kind == nil {
			continue
		}

		if (*ext.ApiVersion == apiVersion) && (*ext.Kind == kind) {
			found = append(found, ext)
		}
	}

	return found
}

// FindExtensionFrom returns a pointer to the [kubeconfig.Extension] with the
// given name inside the provided kubeconfig struct.
//
// If the struct does not contain the cluster, nil will be returned.
func FindExtensionFrom[T HasExtension](name string, extensible *T) *kubeconfig.Extension {
	extensions := ExtensionsFrom(extensible)
	return FindExtension(name, extensions)
}

// FindAuthInfo returns a pointer to the [kubeconfig.Extension] with the given
// apiVersion and kind inside the provided [kubeconfig.Config] struct.
//
// If the configuration does not contain the cluster, nil will be returned.
func FindExtensionsByKindFrom[T HasExtension](apiVersion, kind string, extensible *T) []*kubeconfig.Extension {
	extensions := ExtensionsFrom(extensible)
	return FindExtensionsByKind(apiVersion, kind, extensions)
}
