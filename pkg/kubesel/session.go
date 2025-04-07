package kubesel

import (
	"fmt"
	"os"

	"github.com/eth-p/kubesel/pkg/kubeconfig"
	"gopkg.in/yaml.v3"
)

const (
	managedContextName   = "kubesel"
	managedExtensionName = "managed-by-kubesel"
)

// Session is a kubeconfig file managed by `kubesel`.
type Session struct {
	file    string
	config  *kubeconfig.Config
	context *kubeconfig.Context
	owner   SessionOwner
}

// Save writes the updated [Session] to disk, atomically replacing the its
// prior contents.
func (s *Session) Save() error {
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

// Path returns the path to the session file.
func (s *Session) Path() string {
	return s.file
}

func (s *Session) GetClusterName() string {
	return *s.context.Cluster
}

func (s *Session) GetAuthInfoName() string {
	return *s.context.User
}

func (s *Session) GetNamespace() string {
	return *s.context.Namespace
}

func (s *Session) SetClusterName(name string) {
	*s.context.Cluster = name
}

func (s *Session) SetAuthInfoName(name string) {
	*s.context.User = name
}

func (s *Session) SetNamespace(name string) {
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
