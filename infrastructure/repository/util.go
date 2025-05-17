package repository

// 値が nil または型のゼロ値の場合、NULL として扱う
func ToNullable[T comparable, N any](value *T, isZero func(T) bool, toNull func(T) N) N {
	if value == nil || isZero(*value) {
		var zero N
		return zero
	}
	return toNull(*value)
}
