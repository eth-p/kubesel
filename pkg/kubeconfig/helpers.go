package kubeconfig

import (
	"fmt"

	"github.com/tiendc/go-deepcopy"
)

type cloneInto[T any] interface {
	CloneInto(target *T)
}

// firstNonNil returns the first non-nil pointer.
func firstNonNil[T any](first, other *T) *T {
	if first != nil {
		return first
	} else {
		return other
	}
}

// shallowCopy creates a shallow copy of the pointer contents.
func shallowCopy[T any](value *T) *T {
	if value == nil {
		return nil
	}

	var clone T = *value

	return &clone
}

// shallowCopy creates a shallow copy of the slice.
func shallowCopySlice[T any](slice []T) []T {
	if slice == nil {
		return nil
	}

	clone := make([]T, len(slice))
	copy(clone, slice)

	return clone
}

// cloneSlice creates a deep copy of the slice by calling the [CloneInto]
// method on its elements.
func cloneSlice[T cloneInto[T]](slice []T) []T {
	if slice == nil {
		return nil
	}

	clone := make([]T, len(slice))
	for i, item := range slice {
		item.CloneInto(&clone[i])
	}

	return clone
}

// cloneMap creates a deep copy of a map using reflection.
func cloneMap[T map[K]V, K comparable, V any](dict T) T {
	clone := make(T)
	err := deepcopy.Copy(&clone, dict)
	if err != nil {
		panic(fmt.Errorf("error cloning: %w", err))
	}

	return clone
}
