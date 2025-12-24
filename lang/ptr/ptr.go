package ptr

import "time"

// ToPtr returns a pointer to the given value
func ToPtr[T any](v T) *T {
	return &v
}

// Deref returns the value pointed to by the pointer, or zero value if nil
func Deref[T any](v *T) T {
	if v == nil {
		var zero T
		return zero
	}
	return *v
}

// DerefOrDefault returns the value pointed to by the pointer, or the default value if nil
func DerefOrDefault[T any](v *T, defaultValue T) T {
	if v == nil {
		return defaultValue
	}
	return *v
}

// ToSlice converts a slice of values to a slice of pointers
func ToSlice[T any](vs []T) []*T {
	result := make([]*T, len(vs))
	for i := range vs {
		result[i] = &vs[i]
	}
	return result
}

// FromSlice converts a slice of pointers to a slice of values, skipping nil pointers
func FromSlice[T any](vs []*T) []T {
	result := make([]T, 0, len(vs))
	for _, v := range vs {
		if v != nil {
			result = append(result, *v)
		}
	}
	return result
}

// Common type-specific wrappers for convenience

// ToString returns a pointer to the given string
func ToString(v string) *string {
	return &v
}

// ToInt returns a pointer to the given int
func ToInt(v int) *int {
	return &v
}

// ToInt8 returns a pointer to the given int8
func ToInt8(v int8) *int8 {
	return &v
}

// ToInt16 returns a pointer to the given int16
func ToInt16(v int16) *int16 {
	return &v
}

// ToInt32 returns a pointer to the given int32
func ToInt32(v int32) *int32 {
	return &v
}

// ToInt64 returns a pointer to the given int64
func ToInt64(v int64) *int64 {
	return &v
}

// ToUint returns a pointer to the given uint
func ToUint(v uint) *uint {
	return &v
}

// ToUint8 returns a pointer to the given uint8
func ToUint8(v uint8) *uint8 {
	return &v
}

// ToUint16 returns a pointer to the given uint16
func ToUint16(v uint16) *uint16 {
	return &v
}

// ToUint32 returns a pointer to the given uint32
func ToUint32(v uint32) *uint32 {
	return &v
}

// ToUint64 returns a pointer to the given uint64
func ToUint64(v uint64) *uint64 {
	return &v
}

// ToFloat32 returns a pointer to the given float32
func ToFloat32(v float32) *float32 {
	return &v
}

// ToFloat64 returns a pointer to the given float64
func ToFloat64(v float64) *float64 {
	return &v
}

// ToBool returns a pointer to the given bool
func ToBool(v bool) *bool {
	return &v
}

// ToTime returns a pointer to the given time.Time
func ToTime(v time.Time) *time.Time {
	return &v
}

// ToDuration returns a pointer to the given time.Duration
func ToDuration(v time.Duration) *time.Duration {
	return &v
}
