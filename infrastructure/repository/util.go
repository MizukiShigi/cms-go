package repository

func ToNullable[T comparable, N any](value T, isZero func(T) bool, toNull func(T) N) N {
	if isZero(value) {
		var zero N
		return zero
	}
	return toNull(value)
}
