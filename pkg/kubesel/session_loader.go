package kubesel

import (
	"errors"
	"fmt"

	"github.com/eth-p/kubesel/pkg/kubeconfig"
	"github.com/eth-p/kubesel/pkg/kubeconfig/kcutils"
	"github.com/eth-p/kubesel/pkg/kubeconfig/loader"
)

func newSessionFromLoadedKubeconfig(kc *loader.LoadedKubeconfig) (*Session, error) {
	if len(kc.Errors) > 0 {
		return nil, fmt.Errorf(
			"%w: %s: %w",
			ErrSessionCorrupt,
			kc.Path,
			errors.Join(kc.Errors...),
		)
	}

	// Ensure the current-context is set to the kubsel-managed context.
	kcCurrentContext := ""
	if kc.Config.CurrentContext != nil {
		kcCurrentContext = *kc.Config.CurrentContext
	}

	if kcCurrentContext == "" {
		return nil, fmt.Errorf(
			"%w: the current-context is unset",
			ErrSessionCorrupt,
		)
	}

	if kcCurrentContext != managedContextName {
		return nil, fmt.Errorf(
			"%w: the current-context is not managed by kubesel",
			ErrSessionCorrupt,
		)
	}

	// Find the context.
	kcContext := kcutils.FindContext(kcCurrentContext, &kc.Config)
	if kcContext == nil {
		return nil, fmt.Errorf(
			"%w: the %q context is missing",
			ErrSessionCorrupt,
			kcCurrentContext,
		)
	}

	// TODO: Implement me
	return &Session{
		file:    kc.Path,
		config:  &kc.Config,
		context: kcContext,
	}, nil
}

func newSessionForOwner(sessionFile string, owner SessionOwner) (*Session, error) {
	var kcContextName = managedContextName
	kcContext := &kubeconfig.Context{
		Cluster:   new(string),
		User:      new(string),
		Namespace: new(string),
	}

	return &Session{
		file:    sessionFile,
		owner:   owner,
		context: kcContext,
		config: &kubeconfig.Config{
			CurrentContext: &kcContextName,
			Contexts: []kubeconfig.NamedContext{
				{
					Name:    &kcContextName,
					Context: kcContext,
				},
			},
		},
	}, nil
}
