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

	// Decode the ownership information.
	rawExt := kcutils.FindExtensionFrom(managedExtensionName, &kc.Config)
	if rawExt == nil {
		return nil, fmt.Errorf(
			"%w: the %q extension is missing",
			ErrSessionCorrupt,
			managedExtensionName,
		)
	}

	if !rawExt.Is(kcextApiVersion, kcextManagedByKubeselKind) {
		return nil, fmt.Errorf(
			"%w: the %q extension has the wrong apiVersion or kind",
			ErrSessionCorrupt,
			managedExtensionName,
		)
	}

	var ext kcextManagedByKubesel
	err := kcutils.DecodeExtension(rawExt, &ext)
	if err != nil {
		return nil, fmt.Errorf(
			"%w: could not decode %s: %w",
			ErrSessionCorrupt,
			kcextManagedByKubeselKind,
			err,
		)
	}

	return &Session{
		file:    kc.Path,
		config:  &kc.Config,
		context: kcContext,
		owner:   ext.SessionOwner,
	}, nil
}

func newSessionForOwner(sessionFile string, owner SessionOwner) (*Session, error) {
	kcContext := &kubeconfig.Context{
		Cluster:   new(string),
		User:      new(string),
		Namespace: new(string),
	}

	// Create the ManagedByKubsel extension for the kubeconfig.
	ext := kcextManagedByKubesel{
		SessionOwner: owner,
	}

	extRaw := &kubeconfig.Extension{
		ApiVersion: kcutils.PointerFor(kcextApiVersion),
		Kind:       kcutils.PointerFor(kcextManagedByKubeselKind),
	}

	err := kcutils.EncodeExtension(&ext, extRaw)
	if err != nil {
		panic(fmt.Errorf(
			"failed to encode ManagedByKubesel extension: %w",
			err,
		))
	}

	// Create the kubeconfig.
	kc := &kubeconfig.Config{
		CurrentContext: kcutils.PointerFor(managedContextName),
		Contexts: []kubeconfig.NamedContext{
			{
				Name:    kcutils.PointerFor(managedContextName),
				Context: kcContext,
			},
		},
		Extensions: []kubeconfig.NamedExtension{
			{
				Name:      kcutils.PointerFor(managedExtensionName),
				Extension: extRaw,
			},
		},
	}

	// Return the
	return &Session{
		file:    sessionFile,
		owner:   owner,
		context: kcContext,
		config:  kc,
	}, nil
}
