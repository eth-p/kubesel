package kubesel

import (
	"fmt"
	"os"

	"github.com/eth-p/kubesel/pkg/kubeconfig"
	"github.com/eth-p/kubesel/pkg/kubeconfig/kcutils"
	"gopkg.in/yaml.v3"
)

const (
	managedContextName   = "kubesel"
	managedExtensionName = "managed-by-kubesel"
)

// ManagedKubeconfig is a kubeconfig file managed by `kubesel`.
type ManagedKubeconfig struct {
	file  string
	owner Owner

	config  *kubeconfig.Config
	context *kubeconfig.Context
}

// Save writes the updated [ManagedKubeconfig] to disk, atomically replacing
// its prior contents.
func (s *ManagedKubeconfig) Save() error {
	file, err := os.OpenFile(s.file+".swp", os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}

	marshalled, err := yaml.Marshal(s.config)
	if err != nil {
		return fmt.Errorf("marshalling kubeconfig: %w", err)
	}

	_, err = file.Write(marshalled)
	if err != nil {
		return fmt.Errorf("writing to file: %w", err)
	}

	err = file.Close()
	if err != nil {
		return fmt.Errorf("closing file: %w", err)
	}

	// Rename over existing file for atomic save.
	err = os.Rename(file.Name(), s.file)
	if err != nil {
		return fmt.Errorf("replacing file: %w", err)
	}

	return nil
}

// Path returns the path of the managed kubeconfig file.
func (s *ManagedKubeconfig) Path() string {
	return s.file
}

// GetClusterName returns the name of the active [kubeconfig.Cluster] in
// the kubesel-managed kubeconfig.
func (s *ManagedKubeconfig) GetClusterName() string {
	return *s.context.Cluster
}

// GetAuthInfoName returns the name of the active [kubeconfig.AuthInfo] in
// the kubesel-managed kubeconfig.
func (s *ManagedKubeconfig) GetAuthInfoName() string {
	return *s.context.User
}

// GetNamespace returns the name of the active namespace in the kubesel-managed
// kubeconfig.
func (s *ManagedKubeconfig) GetNamespace() string {
	return *s.context.Namespace
}

// SetClusterName changes the active [kubeconfig.Cluster] in the kubesel-managed
// kubeconfig. To commit the change [ManagedKubeconfig.Save] should be called
// after.
func (s *ManagedKubeconfig) SetClusterName(name string) {
	*s.context.Cluster = name
}

// SetAuthInfoName changes the active [kubeconfig.AuthInfo] in the
// kubesel-managed kubeconfig. To commit the change [ManagedKubeconfig.Save]
// should be called after.
func (s *ManagedKubeconfig) SetAuthInfoName(name string) {
	*s.context.User = name
}

// SetNamespace changes the active namespace in the kubesel-managed kubeconfig.
// To commit the change [ManagedKubeconfig.Save] should be called after.
func (s *ManagedKubeconfig) SetNamespace(name string) {
	*s.context.Namespace = name
}

// IsManagedContext checks if the provided [kubeconfig.NamedContext] is managed
// by kubesel.
func IsManagedContext(kcNamedContext *kubeconfig.NamedContext) bool {
	if kcNamedContext.Name != nil && (*kcNamedContext.Name == managedContextName) {
		return true
	}

	return false
}

// IsManagedKubeconfig checks if the provided [kubeconfig.Config] is managed
// by kubesel.
func IsManagedKubeconfig(kc *kubeconfig.Config) bool {
	ext := kcutils.FindExtensionFrom(managedExtensionName, kc)
	if ext == nil {
		return false
	}

	return ext.Is(kcextApiVersion, kcextManagedByKubeselKind)
}
