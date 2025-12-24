package slices

import (
	"math/rand"
	"sort"
)

// Contains checks if a slice contains a value
func Contains[T comparable](slice []T, value T) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// ContainsAny checks if a slice contains any of the given values
func ContainsAny[T comparable](slice []T, values ...T) bool {
	for _, v := range slice {
		for _, val := range values {
			if v == val {
				return true
			}
		}
	}
	return false
}

// ContainsAll checks if a slice contains all of the given values
func ContainsAll[T comparable](slice []T, values ...T) bool {
	for _, val := range values {
		if !Contains(slice, val) {
			return false
		}
	}
	return true
}

// Filter returns a new slice containing only the elements that satisfy the predicate
func Filter[T any](slice []T, predicate func(T) bool) []T {
	result := make([]T, 0, len(slice))
	for _, v := range slice {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

// Map transforms a slice using a mapper function
func Map[T any, U any](slice []T, mapper func(T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = mapper(v)
	}
	return result
}

// Reduce reduces a slice to a single value
func Reduce[T any, U any](slice []T, initial U, reducer func(U, T) U) U {
	result := initial
	for _, v := range slice {
		result = reducer(result, v)
	}
	return result
}

// Find returns the first element that satisfies the predicate, or zero value if not found
func Find[T any](slice []T, predicate func(T) bool) (T, bool) {
	for _, v := range slice {
		if predicate(v) {
			return v, true
		}
	}
	var zero T
	return zero, false
}

// FindIndex returns the index of the first element that satisfies the predicate, or -1 if not found
func FindIndex[T any](slice []T, predicate func(T) bool) int {
	for i, v := range slice {
		if predicate(v) {
			return i
		}
	}
	return -1
}

// IndexOf returns the index of the first occurrence of value, or -1 if not found
func IndexOf[T comparable](slice []T, value T) int {
	for i, v := range slice {
		if v == value {
			return i
		}
	}
	return -1
}

// LastIndexOf returns the index of the last occurrence of value, or -1 if not found
func LastIndexOf[T comparable](slice []T, value T) int {
	for i := len(slice) - 1; i >= 0; i-- {
		if slice[i] == value {
			return i
		}
	}
	return -1
}

// Unique removes duplicate elements from a slice
func Unique[T comparable](slice []T) []T {
	seen := make(map[T]struct{})
	result := make([]T, 0, len(slice))
	for _, v := range slice {
		if _, exists := seen[v]; !exists {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

// Reverse reverses a slice
func Reverse[T any](slice []T) []T {
	result := make([]T, len(slice))
	for i, v := range slice {
		result[len(slice)-1-i] = v
	}
	return result
}

// Chunk splits a slice into chunks of the given size
func Chunk[T any](slice []T, size int) [][]T {
	if size <= 0 {
		return [][]T{}
	}

	var chunks [][]T
	for i := 0; i < len(slice); i += size {
		end := i + size
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}

// Flatten flattens a 2D slice into a 1D slice
func Flatten[T any](slices [][]T) []T {
	var result []T
	for _, slice := range slices {
		result = append(result, slice...)
	}
	return result
}

// Shuffle returns a new shuffled slice
func Shuffle[T any](slice []T) []T {
	result := make([]T, len(slice))
	copy(result, slice)

	// Fisher-Yates shuffle
	for i := len(result) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		result[i], result[j] = result[j], result[i]
	}
	return result
}

// ToMap converts a slice to a map using a key selector
func ToMap[T any, K comparable](slice []T, keySelector func(T) K) map[K]T {
	result := make(map[K]T, len(slice))
	for _, v := range slice {
		result[keySelector(v)] = v
	}
	return result
}

// GroupBy groups elements of a slice by a key selector
func GroupBy[T any, K comparable](slice []T, keySelector func(T) K) map[K][]T {
	result := make(map[K][]T)
	for _, v := range slice {
		key := keySelector(v)
		result[key] = append(result[key], v)
	}
	return result
}

// Delete removes elements at the given indices
func Delete[T any](slice []T, indices ...int) []T {
	// Sort indices in descending order to avoid index shifting
	sort.Sort(sort.Reverse(sort.IntSlice(indices)))

	result := slice
	for _, idx := range indices {
		if idx >= 0 && idx < len(result) {
			result = append(result[:idx], result[idx+1:]...)
		}
	}
	return result
}

// Insert inserts elements at the given index
func Insert[T any](slice []T, index int, elements ...T) []T {
	if index < 0 {
		index = 0
	}
	if index > len(slice) {
		index = len(slice)
	}

	result := make([]T, 0, len(slice)+len(elements))
	result = append(result, slice[:index]...)
	result = append(result, elements...)
	result = append(result, slice[index:]...)
	return result
}

// Concat concatenates multiple slices
func Concat[T any](slices ...[]T) []T {
	var totalLen int
	for _, slice := range slices {
		totalLen += len(slice)
	}

	result := make([]T, 0, totalLen)
	for _, slice := range slices {
		result = append(result, slice...)
	}
	return result
}

// Intersect returns the intersection of two slices
func Intersect[T comparable](a, b []T) []T {
	set := make(map[T]struct{})
	for _, v := range b {
		set[v] = struct{}{}
	}

	result := make([]T, 0)
	for _, v := range a {
		if _, exists := set[v]; exists {
			result = append(result, v)
		}
	}
	return Unique(result)
}

// Difference returns elements in a that are not in b
func Difference[T comparable](a, b []T) []T {
	set := make(map[T]struct{})
	for _, v := range b {
		set[v] = struct{}{}
	}

	result := make([]T, 0)
	for _, v := range a {
		if _, exists := set[v]; !exists {
			result = append(result, v)
		}
	}
	return result
}

// Each iterates over a slice and calls the function for each element
func Each[T any](slice []T, fn func(T)) {
	for _, v := range slice {
		fn(v)
	}
}

// EachIndex iterates over a slice and calls the function with index and element
func EachIndex[T any](slice []T, fn func(int, T)) {
	for i, v := range slice {
		fn(i, v)
	}
}

// Any returns true if any element satisfies the predicate
func Any[T any](slice []T, predicate func(T) bool) bool {
	for _, v := range slice {
		if predicate(v) {
			return true
		}
	}
	return false
}

// All returns true if all elements satisfy the predicate
func All[T any](slice []T, predicate func(T) bool) bool {
	for _, v := range slice {
		if !predicate(v) {
			return false
		}
	}
	return true
}

// None returns true if no element satisfies the predicate
func None[T any](slice []T, predicate func(T) bool) bool {
	return !Any(slice, predicate)
}

// Count counts the number of elements that satisfy the predicate
func Count[T any](slice []T, predicate func(T) bool) int {
	count := 0
	for _, v := range slice {
		if predicate(v) {
			count++
		}
	}
	return count
}

// Max returns the maximum element in a slice
func Max[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64](slice []T) T {
	if len(slice) == 0 {
		var zero T
		return zero
	}

	max := slice[0]
	for _, v := range slice[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

// Min returns the minimum element in a slice
func Min[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64](slice []T) T {
	if len(slice) == 0 {
		var zero T
		return zero
	}

	min := slice[0]
	for _, v := range slice[1:] {
		if v < min {
			min = v
		}
	}
	return min
}

// Sum returns the sum of all elements in a slice
func Sum[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64](slice []T) T {
	var sum T
	for _, v := range slice {
		sum += v
	}
	return sum
}

// Average returns the average of all elements in a slice
func Average[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64](slice []T) float64 {
	if len(slice) == 0 {
		return 0
	}

	sum := Sum(slice)
	return float64(sum) / float64(len(slice))
}
