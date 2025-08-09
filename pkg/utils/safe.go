package utils

func Deref[T any](val *T) T {
	if val == nil {
		var zero T
		return zero
	}
	return *val
}
