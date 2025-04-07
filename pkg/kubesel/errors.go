package kubesel

import (
	"errors"
)

var (
	// ErrAlreadyManaged is returned when trying to create a new
	// [ManagedKubeconfig] with an owner that already associated with a
	// managed kubeconfig file.
	ErrAlreadyManaged = errors.New("already managing a kubeconfig file")

	// ErrManagedKubeconfigCorrupt is returned when a [ManagedKubeconfig] file
	// cannot be loaded.
	ErrManagedKubeconfigCorrupt = errors.New("managed kubeconfig file is invalid")

	// ErrUnmanaged is returned when a  [ManagedKubeconfig] does not exist.
	ErrUnmanaged = errors.New("no kubesel-managed kubeconfig file")

	// ErrOwnerProcessNotExist is returned when trying to create a
	// [ManagedKubeconfig] whose owner is not a living process.
	ErrOwnerProcessNotExist = errors.New("owner process does not exist")
)
