package kcutils

func PointerFor[T any](value T) *T {
	return &value
}
