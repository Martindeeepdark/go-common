package maps

// Keys returns all keys from a map
func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Values returns all values from a map
func Values[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

// Entries returns all key-value pairs from a map
func Entries[K comparable, V any](m map[K]V) []Entry[K, V] {
	entries := make([]Entry[K, V], 0, len(m))
	for k, v := range m {
		entries = append(entries, Entry[K, V]{Key: k, Value: v})
	}
	return entries
}

// Entry represents a key-value pair
type Entry[K comparable, V any] struct {
	Key   K
	Value V
}

// FromEntries creates a map from entries
func FromEntries[K comparable, V any](entries []Entry[K, V]) map[K]V {
	result := make(map[K]V, len(entries))
	for _, entry := range entries {
		result[entry.Key] = entry.Value
	}
	return result
}

// Clone creates a shallow copy of a map
func Clone[K comparable, V any](m map[K]V) map[K]V {
	result := make(map[K]V, len(m))
	for k, v := range m {
		result[k] = v
	}
	return result
}

// Merge merges multiple maps into one, later maps override earlier ones for duplicate keys
func Merge[K comparable, V any](maps ...map[K]V) map[K]V {
	result := make(map[K]V)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// MergeInto merges source map into destination map
func MergeInto[K comparable, V any](dest, src map[K]V) {
	for k, v := range src {
		dest[k] = v
	}
}

// HasKey checks if a map contains a key
func HasKey[K comparable, V any](m map[K]V, key K) bool {
	_, ok := m[key]
	return ok
}

// GetOrDefault returns the value for key, or defaultValue if key doesn't exist
func GetOrDefault[K comparable, V any](m map[K]V, key K, defaultValue V) V {
	if v, ok := m[key]; ok {
		return v
	}
	return defaultValue
}

// GetOrSet returns the value for key, or sets and returns defaultValue if key doesn't exist
func GetOrSet[K comparable, V any](m map[K]V, key K, defaultValue V) V {
	if v, ok := m[key]; ok {
		return v
	}
	m[key] = defaultValue
	return defaultValue
}

// DeleteKeys deletes multiple keys from a map
func DeleteKeys[K comparable, V any](m map[K]V, keys ...K) {
	for _, key := range keys {
		delete(m, key)
	}
}

// KeepOnly keeps only the specified keys in a map
func KeepOnly[K comparable, V any](m map[K]V, keys ...K) {
	keySet := make(map[K]struct{})
	for _, key := range keys {
		keySet[key] = struct{}{}
	}

	for k := range m {
		if _, ok := keySet[k]; !ok {
			delete(m, k)
		}
	}
}

// Clear removes all entries from a map
func Clear[K comparable, V any](m map[K]V) {
	for k := range m {
		delete(m, k)
	}
}

// IsEmpty checks if a map is empty
func IsEmpty[K comparable, V any](m map[K]V) bool {
	return len(m) == 0
}

// Size returns the number of entries in a map
func Size[K comparable, V any](m map[K]V) int {
	return len(m)
}

// Filter filters map entries by predicate
func Filter[K comparable, V any](m map[K]V, predicate func(K, V) bool) map[K]V {
	result := make(map[K]V)
	for k, v := range m {
		if predicate(k, v) {
			result[k] = v
		}
	}
	return result
}

// MapValues transforms map values using a mapper function
func MapValues[K comparable, V any, U any](m map[K]V, mapper func(K, V) U) map[K]U {
	result := make(map[K]U, len(m))
	for k, v := range m {
		result[k] = mapper(k, v)
	}
	return result
}

// MapKeys transforms map keys using a mapper function
func MapKeys[K comparable, L comparable, V any](m map[K]V, mapper func(K, V) L) map[L]V {
	result := make(map[L]V, len(m))
	for k, v := range m {
		newKey := mapper(k, v)
		result[newKey] = v
	}
	return result
}

// Invert inverts a map (key becomes value, value becomes key)
// Values must be unique, otherwise duplicate keys will be overwritten
func Invert[K comparable, V comparable](m map[K]V) map[V]K {
	result := make(map[V]K, len(m))
	for k, v := range m {
		result[v] = k
	}
	return result
}

// Each iterates over map entries and calls the function for each entry
func Each[K comparable, V any](m map[K]V, fn func(K, V)) {
	for k, v := range m {
		fn(k, v)
	}
}

// Any returns true if any entry satisfies the predicate
func Any[K comparable, V any](m map[K]V, predicate func(K, V) bool) bool {
	for k, v := range m {
		if predicate(k, v) {
			return true
		}
	}
	return false
}

// All returns true if all entries satisfy the predicate
func All[K comparable, V any](m map[K]V, predicate func(K, V) bool) bool {
	for k, v := range m {
		if !predicate(k, v) {
			return false
		}
	}
	return true
}

// Equal checks if two maps are equal (same keys and values)
func Equal[K comparable, V comparable](a, b map[K]V) bool {
	if len(a) != len(b) {
		return false
	}

	for k, va := range a {
		if vb, ok := b[k]; !ok || va != vb {
			return false
		}
	}

	return true
}

// Diff returns entries that are in a but not in b
func Diff[K comparable, V comparable](a, b map[K]V) map[K]V {
	result := make(map[K]V)
	for k, va := range a {
		if vb, ok := b[k]; !ok || va != vb {
			result[k] = va
		}
	}
	return result
}

// Intersect returns entries that are in both maps
func Intersect[K comparable, V comparable](a, b map[K]V) map[K]V {
	result := make(map[K]V)
	for k, va := range a {
		if vb, ok := b[k]; ok && va == vb {
			result[k] = va
		}
	}
	return result
}

// Union returns the union of two maps (b's values override a's for duplicate keys)
func Union[K comparable, V any](a, b map[K]V) map[K]V {
	return Merge(a, b)
}
