package ternary

// T returns trueValue if condition is true, otherwise falseValue
func T[T any](condition bool, trueValue, falseValue T) T {
	if condition {
		return trueValue
	}
	return falseValue
}

// LazyT returns trueValueFunc() if condition is true, otherwise falseValueFunc()
// Useful for avoiding evaluation of both branches
func Lazy[T any](condition bool, trueValueFunc, falseValueFunc func() T) T {
	if condition {
		return trueValueFunc()
	}
	return falseValueFunc()
}

// Or returns first if it's not zero value, otherwise second
func Or[T comparable](first, second T) T {
	var zero T
	if first != zero {
		return first
	}
	return second
}

// OrFunc returns first if it's not zero value, otherwise calls second()
func OrFunc[T comparable](first T, second func() T) T {
	var zero T
	if first != zero {
		return first
	}
	return second()
}

// And returns second if condition is true and first is not zero value, otherwise zero value
func And[T comparable](condition bool, first T) T {
	var zero T
	if !condition {
		return zero
	}
	return first
}

// Coalesce returns the first non-zero value from args, or zero value if all are zero
func Coalesce[T comparable](args ...T) T {
	var zero T
	for _, arg := range args {
		if arg != zero {
			return arg
		}
	}
	return zero
}

// IsZero returns true if value is zero value
func IsZero[T comparable](value T) bool {
	var zero T
	return value == zero
}

// IsNotZero returns true if value is not zero value
func IsNotZero[T comparable](value T) bool {
	var zero T
	return value != zero
}

// If returns value if condition is true, otherwise zero value
func If[T any](condition bool, value T) T {
	if condition {
		return value
	}
	var zero T
	return zero
}

// Unless returns value if condition is false, otherwise zero value
func Unless[T any](condition bool, value T) T {
	if !condition {
		return value
	}
	var zero T
	return zero
}

// Switch returns the value from the first case where the key matches, or defaultValue if no match
func Switch[K comparable, V any](key K, cases map[K]V, defaultValue V) V {
	if value, ok := cases[key]; ok {
		return value
	}
	return defaultValue
}

// Default returns value if it's not zero value, otherwise defaultValue
func Default[T comparable](value, defaultValue T) T {
	var zero T
	if value == zero {
		return defaultValue
	}
	return value
}
