package kubesel

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"path/filepath"
	"slices"

	"github.com/eth-p/kubesel/pkg/kubeconfig/loader"
)

type GarbageCollectOptions struct {
	MaxFilesToCheck  int
	MaxFilesToDelete int
}

type GarbageCollectResult struct {
	FilesDeleted []string
	FilesChecked []string
	Errors       []error
}

// GarbageCollect removes kubesel-managed files belonging to processes which
// are no longer alive.
func (k *Kubesel) GarbageCollect(opts *GarbageCollectOptions) (*GarbageCollectResult, error) {
	err := k.ensureSessionsDirExists()
	if err != nil {
		return nil, fmt.Errorf("error ensuring session directory exists: %w", err)
	}

	// Get the files in the session directory.
	sessionDir := k.sessionDir
	entries, err := os.ReadDir(sessionDir)
	if err != nil {
		return nil, fmt.Errorf("error listing session directory: %w", err)
	}

	// Filter the files to only ones (likely) created by kubesel.
	// Consider any file ending in `.yaml` to be one.
	entries = slices.DeleteFunc(entries, func(entry os.DirEntry) bool {
		return filepath.Ext(entry.Name()) != ".yaml"
	})

	// Prepare.
	var filesChecked []string
	maxChecks := opts.MaxFilesToCheck
	if maxChecks == 0 {
		maxChecks = math.MaxInt
	}

	var filesDeleted []string
	maxDeletes := opts.MaxFilesToDelete
	if maxDeletes == 0 {
		maxDeletes = math.MaxInt
	}

	var errs []error

	// Check the files in a nondeterministic order.
	// If either "MaxFiles{Checked,Deleted}" limit is reached, stop.
	order := rand.Perm(len(entries))
	for _, index := range order {
		path := filepath.Join(sessionDir, entries[index].Name())
		canDelete, err := k.canGarbageCollect(path)

		// If it's ok to delete the file, try to do it.
		if err == nil && canDelete {
			err = os.Remove(path)
		}

		// Record the results.
		if err != nil {
			errs = append(errs, fmt.Errorf("gc error: %s: %w", path, err))
			continue
		}

		filesChecked = append(filesChecked, path)
		if len(filesChecked) > maxChecks {
			break
		}

		if canDelete {
			filesDeleted = append(filesDeleted, path)
			if len(filesDeleted) > maxDeletes {
				break
			}
		}
	}

	return &GarbageCollectResult{
		FilesDeleted: filesDeleted,
		FilesChecked: filesChecked,
		Errors:       errs,
	}, nil
}

func (k *Kubesel) canGarbageCollect(path string) (bool, error) {
	kc := loader.LoadFromFile(path)

	// If the file can't be parsed as a kubeconfig file, don't touch it.
	if len(kc.Errors) > 0 {
		return false, errors.Join(kc.Errors...)
	}

	// Try to convert it into a managed kubeconfig file.
	// If it's corrupt, it can be deleted.
	managedKc, err := newManagedKubeconfigFromExistingKubeconfig(kc)
	if errors.Is(err, ErrManagedKubeconfigCorrupt) {
		return true, nil
	}

	// If it can't be loaded, don't touch it.
	if err != nil {
		return false, err
	}

	// If the owner isn't alive, it can be deleted.
	isAlive, err := managedKc.owner.IsAlive()
	if err != nil {
		return false, err
	}

	return !isAlive, nil
}
