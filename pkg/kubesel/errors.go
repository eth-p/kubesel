package kubesel

import (
	"errors"
)

var (
	// ErrSessionExists is returned when trying to create a new session
	// with the same owner as an existing session.
	ErrSessionExists = errors.New("kubesel session already exists")

	// ErrSessionCorrupt is returned when a session file cannot be loaded
	// as kubectl configuration.
	ErrSessionCorrupt = errors.New("kubesel session file is invalid")

	// ErrNoSession is returned when a session does not exist.
	ErrNoSession = errors.New("no kubesel session")

	// ErrOwnerProcessNotExist is returned when trying to create a session
	// whose owner is not a living process.
	ErrOwnerProcessNotExist = errors.New("kubesel session owner process does not exist")
)
