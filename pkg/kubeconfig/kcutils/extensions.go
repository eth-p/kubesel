package kcutils

import (
	"fmt"

	"github.com/eth-p/kubesel/pkg/kubeconfig"
)

type HasExtension interface {
	kubeconfig.Config | kubeconfig.Cluster |
		kubeconfig.AuthInfo | kubeconfig.Context |
		kubeconfig.Preferences
}

// ExtensionsFrom returns the [kubeconfig.NamedExtension] attached to the
// provided kubeconfig struct type.
func ExtensionsFrom[T HasExtension](extensible *T) []kubeconfig.NamedExtension {
	switch refined := any(*extensible).(type) {
	case kubeconfig.Config:
		return refined.Extensions
	case kubeconfig.Cluster:
		return refined.Extensions
	case kubeconfig.AuthInfo:
		return refined.Extensions
	case kubeconfig.Context:
		return refined.Extensions
	case kubeconfig.Preferences:
		return refined.Extensions
	}

	panic(fmt.Sprintf("unsupported type: %T", extensible))
}
