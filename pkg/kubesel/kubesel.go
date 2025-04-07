package kubesel

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/eth-p/kubesel/pkg/kubeconfig"
	"github.com/eth-p/kubesel/pkg/kubeconfig/loader"
)

type Kubesel struct {
	kubeconfigs *loader.LoadedKubeconfigCollection

	dataDir    string
	sessionDir string

	lazyCurrentSession func() (*Session, error)
	lazyClusterNames   func() []string
	lazyAuthInfoNames  func() []string
	lazyContextNames   func() []string
}

// NewKubesel reads the kubectl configuration files and sets up this instance
// of kubesel.
func NewKubesel() (*Kubesel, error) {
	dataDir := filepath.Join(findDataHomeDir(), "kubesel")
	sessionDir := filepath.Join(dataDir, "sessions")

	// Load the kubeconfig files.
	kcFiles, err := loader.FindKubeConfigFiles()
	if err != nil {
		return nil, fmt.Errorf("error finding kubeconfig files: %w", err)
	}

	kubeconfigs := loader.LoadMultipleFiles(kcFiles)

	// Create the Kubesel instance.
	kubesel := &Kubesel{
		kubeconfigs: kubeconfigs,
		sessionDir:  sessionDir,
		dataDir:     dataDir,
	}

	kubesel.lazyCurrentSession = sync.OnceValues(kubesel.findCurrentSession)
	kubesel.lazyClusterNames = sync.OnceValue(kubesel.findClusterNames)
	kubesel.lazyAuthInfoNames = sync.OnceValue(kubesel.findAuthInfoNames)
	kubesel.lazyContextNames = sync.OnceValue(kubesel.findContextNames)
	return kubesel, nil
}

// GetMergedKubeconfig returns merged contents of the files specified by the
// `KUBECONFIG` environment variable.
func (k *Kubesel) GetMergedKubeconfig() *kubeconfig.Config {
	return k.kubeconfigs.Merged
}

// CurrentSession returns the current kubesel [Session], if one exists.
// If one does not exist, this returns [ErrNoSession].
//
// The current session the considered to be the first kubeconfig file inside
// the `KUBECONFIG` environment variable that is located inside kubesel's
// sessions directory.
func (k *Kubesel) CurrentSession() (*Session, error) {
	return k.lazyCurrentSession()
}

// GetClusterNames returns the list of known [kubeconfig.NamedCluster] names
// inside the merged kubeconfig.
func (k *Kubesel) GetClusterNames() []string {
	return k.lazyClusterNames()
}

// GetAuthInfoNames returns the list of known [kubeconfig.NamedAuthInfo] names
// inside the merged kubeconfig.
func (k *Kubesel) GetAuthInfoNames() []string {
	return k.lazyAuthInfoNames()
}

// GetContextNames returns the list of known [kubeconfig.NamedContext] names
// inside the merged kubeconfig.
func (k *Kubesel) GetContextNames() []string {
	return k.lazyContextNames()
}

// CreateSession creates a new kubesel [Session] for the given [SessionOwner].
//
// If there is already a session for the given owner, this will return an
// [ErrSessionExists] error.
func (k *Kubesel) CreateSession(owner SessionOwner) (*Session, error) {
	if err := k.ensureSessionsDirExists(); err != nil {
		return nil, err
	}

	// Check if the session file already exists.
	sessionFile := filepath.Join(k.sessionDir, owner.fileName())
	if _, err := os.Stat(sessionFile); !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("%w: owner pid %d", ErrSessionExists, owner.pid)
	}

	// If it does not, create it.
	session, err := newSessionForOwner(sessionFile, owner)
	if err != nil {
		return nil, err
	}

	err = session.Save()
	if err != nil {
		return nil, fmt.Errorf("error saving session: %w", err)
	}

	return session, nil
}

// findCurrentSession looks for the first kubesel session file found within
// the loaded kubeconfig files.
func (k *Kubesel) findCurrentSession() (*Session, error) {
	for _, kc := range k.kubeconfigs.Configs {
		rel, err := filepath.Rel(k.sessionDir, kc.Path)
		if err != nil || strings.HasPrefix(rel, "..") {
			continue
		}

		return newSessionFromLoadedKubeconfig(kc)
	}

	return nil, ErrNoSession
}

// ensureSessionDirExists creates the directory containing kubesel session
// files if it does not already exist.
func (k *Kubesel) ensureSessionsDirExists() error {
	err := os.MkdirAll(k.sessionDir, 0o700)
	if err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	return nil
}

// findClusterNames returns all the cluster names in the merged kubeconfig.
func (k *Kubesel) findClusterNames() []string {
	names := make([]string, 0, len(k.kubeconfigs.Merged.Clusters))
	for _, kcCluster := range k.kubeconfigs.Merged.Clusters {
		if kcCluster.Name != nil {
			names = append(names, *kcCluster.Name)
		}
	}

	return names
}

// findAuthInfoNames returns all the authinfo names in the merged kubeconfig.
func (k *Kubesel) findAuthInfoNames() []string {
	names := make([]string, 0, len(k.kubeconfigs.Merged.AuthInfos))
	for _, kcAuthInfo := range k.kubeconfigs.Merged.AuthInfos {
		if kcAuthInfo.Name != nil {
			names = append(names, *kcAuthInfo.Name)
		}
	}

	return names
}

// findContextNames returns all the context names in the merged kubeconfig.
func (k *Kubesel) findContextNames() []string {
	names := make([]string, 0, len(k.kubeconfigs.Merged.Contexts))
	for _, kcContext := range k.kubeconfigs.Merged.Contexts {
		if kcContext.Name != nil && !IsManagedContext(&kcContext) {
			names = append(names, *kcContext.Name)
		}
	}

	return names
}
