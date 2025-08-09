package utils

func Deref[T any](val *T) T {
	if val == nil {
		var zero T
		return zero
	}
	return *val
}

// Ptr 返回指向给定值的指针
func Ptr[T any](val T) *T {
	return &val
}
