package testutil

func PtrFrom[T any](v T) *T {
	return &v
}
