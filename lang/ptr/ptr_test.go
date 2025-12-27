package ptr

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToPtr(t *testing.T) {
	t.Run("basic types", func(t *testing.T) {
		p := ToPtr(42)
		assert.NotNil(t, p)
		assert.Equal(t, 42, *p)

		s := ToPtr("hello")
		assert.NotNil(t, s)
		assert.Equal(t, "hello", *s)

		b := ToPtr(true)
		assert.NotNil(t, b)
		assert.Equal(t, true, *b)
	})

	t.Run("zero values", func(t *testing.T) {
		zp := ToPtr(0)
		assert.NotNil(t, zp)
		assert.Equal(t, 0, *zp)

		zs := ToPtr("")
		assert.NotNil(t, zs)
		assert.Equal(t, "", *zs)

		zb := ToPtr(false)
		assert.NotNil(t, zb)
		assert.Equal(t, false, *zb)
	})

	t.Run("struct types", func(t *testing.T) {
		now := time.Now()
		tp := ToPtr(now)
		assert.NotNil(t, tp)
		assert.Equal(t, now, *tp)

		dp := ToPtr(time.Second)
		assert.NotNil(t, dp)
		assert.Equal(t, time.Second, *dp)
	})
}

func TestDeref(t *testing.T) {
	t.Run("non-nil pointer", func(t *testing.T) {
		x := 42
		assert.Equal(t, 42, Deref(&x))

		s := "hello"
		assert.Equal(t, "hello", Deref(&s))

		b := true
		assert.Equal(t, true, Deref(&b))
	})

	t.Run("nil pointer returns zero value", func(t *testing.T) {
		assert.Equal(t, 0, Deref((*int)(nil)))
		assert.Equal(t, "", Deref((*string)(nil)))
		assert.Equal(t, false, Deref((*bool)(nil)))
		assert.Equal(t, time.Time{}, Deref((*time.Time)(nil)))
	})
}

func TestDerefOrDefault(t *testing.T) {
	t.Run("non-nil pointer", func(t *testing.T) {
		x := 42
		assert.Equal(t, 42, DerefOrDefault(&x, 100))

		s := "hello"
		assert.Equal(t, "hello", DerefOrDefault(&s, "default"))
	})

	t.Run("nil pointer returns default value", func(t *testing.T) {
		assert.Equal(t, 100, DerefOrDefault((*int)(nil), 100))
		assert.Equal(t, "default", DerefOrDefault((*string)(nil), "default"))
		assert.Equal(t, true, DerefOrDefault((*bool)(nil), true))
	})
}

func TestToSlice(t *testing.T) {
	t.Run("convert slice to pointer slice", func(t *testing.T) {
		values := []int{1, 2, 3}
		ptrs := ToSlice(values)

		assert.Len(t, ptrs, 3)
		assert.Equal(t, 1, *ptrs[0])
		assert.Equal(t, 2, *ptrs[1])
		assert.Equal(t, 3, *ptrs[2])
	})

	t.Run("empty slice", func(t *testing.T) {
		ptrs := ToSlice([]int{})
		assert.Len(t, ptrs, 0)

		ptrs2 := ToSlice([]string(nil))
		assert.Len(t, ptrs2, 0)
	})

	t.Run("pointers point to original elements", func(t *testing.T) {
		values := []int{1, 2, 3}
		ptrs := ToSlice(values)

		// Pointers point to original slice elements
		values[0] = 999
		assert.Equal(t, 999, *ptrs[0]) // pointers reference the original elements
	})
}

func TestFromSlice(t *testing.T) {
	t.Run("convert pointer slice to value slice", func(t *testing.T) {
		a, b, c := 1, 2, 3
		ptrs := []*int{&a, &b, &c}
		values := FromSlice(ptrs)

		assert.Len(t, values, 3)
		assert.Equal(t, 1, values[0])
		assert.Equal(t, 2, values[1])
		assert.Equal(t, 3, values[2])
	})

	t.Run("skips nil pointers", func(t *testing.T) {
		a, c := 1, 3
		ptrs := []*int{&a, nil, &c}
		values := FromSlice(ptrs)

		assert.Len(t, values, 2)
		assert.Equal(t, 1, values[0])
		assert.Equal(t, 3, values[1])
	})

	t.Run("all nil pointers", func(t *testing.T) {
		ptrs := []*int{nil, nil, nil}
		values := FromSlice(ptrs)
		assert.Len(t, values, 0)
	})

	t.Run("empty slice", func(t *testing.T) {
		values := FromSlice([]*int{})
		assert.Len(t, values, 0)
	})
}

func TestTypeSpecificWrappers(t *testing.T) {
	t.Run("ToString", func(t *testing.T) {
		s := "test"
		p := ToString(s)
		assert.NotNil(t, p)
		assert.Equal(t, "test", *p)
	})

	t.Run("ToInt", func(t *testing.T) {
		p := ToInt(42)
		assert.NotNil(t, p)
		assert.Equal(t, 42, *p)
	})

	t.Run("ToInt8", func(t *testing.T) {
		var x int8 = 8
		p := ToInt8(x)
		assert.NotNil(t, p)
		assert.Equal(t, int8(8), *p)
	})

	t.Run("ToInt16", func(t *testing.T) {
		var x int16 = 16
		p := ToInt16(x)
		assert.NotNil(t, p)
		assert.Equal(t, int16(16), *p)
	})

	t.Run("ToInt32", func(t *testing.T) {
		var x int32 = 32
		p := ToInt32(x)
		assert.NotNil(t, p)
		assert.Equal(t, int32(32), *p)
	})

	t.Run("ToInt64", func(t *testing.T) {
		var x int64 = 64
		p := ToInt64(x)
		assert.NotNil(t, p)
		assert.Equal(t, int64(64), *p)
	})

	t.Run("ToUint", func(t *testing.T) {
		var x uint = 42
		p := ToUint(x)
		assert.NotNil(t, p)
		assert.Equal(t, uint(42), *p)
	})

	t.Run("ToUint8", func(t *testing.T) {
		var x uint8 = 8
		p := ToUint8(x)
		assert.NotNil(t, p)
		assert.Equal(t, uint8(8), *p)
	})

	t.Run("ToUint16", func(t *testing.T) {
		var x uint16 = 16
		p := ToUint16(x)
		assert.NotNil(t, p)
		assert.Equal(t, uint16(16), *p)
	})

	t.Run("ToUint32", func(t *testing.T) {
		var x uint32 = 32
		p := ToUint32(x)
		assert.NotNil(t, p)
		assert.Equal(t, uint32(32), *p)
	})

	t.Run("ToUint64", func(t *testing.T) {
		var x uint64 = 64
		p := ToUint64(x)
		assert.NotNil(t, p)
		assert.Equal(t, uint64(64), *p)
	})

	t.Run("ToFloat32", func(t *testing.T) {
		var x float32 = 3.14
		p := ToFloat32(x)
		assert.NotNil(t, p)
		assert.Equal(t, float32(3.14), *p)
	})

	t.Run("ToFloat64", func(t *testing.T) {
		var x float64 = 3.14159
		p := ToFloat64(x)
		assert.NotNil(t, p)
		assert.Equal(t, float64(3.14159), *p)
	})

	t.Run("ToBool", func(t *testing.T) {
		p := ToBool(true)
		assert.NotNil(t, p)
		assert.Equal(t, true, *p)

		p = ToBool(false)
		assert.NotNil(t, p)
		assert.Equal(t, false, *p)
	})

	t.Run("ToTime", func(t *testing.T) {
		now := time.Now()
		p := ToTime(now)
		assert.NotNil(t, p)
		assert.Equal(t, now, *p)
	})

	t.Run("ToDuration", func(t *testing.T) {
		d := time.Hour
		p := ToDuration(d)
		assert.NotNil(t, p)
		assert.Equal(t, time.Hour, *p)
	})
}
