package conv

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToString(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		assert.Equal(t, "hello", ToString("hello"))
	})

	t.Run("integers", func(t *testing.T) {
		assert.Equal(t, "42", ToString(42))
		assert.Equal(t, "-42", ToString(-42))
		assert.Equal(t, "8", ToString(int8(8)))
		assert.Equal(t, "16", ToString(int16(16)))
		assert.Equal(t, "32", ToString(int32(32)))
		assert.Equal(t, "64", ToString(int64(64)))
	})

	t.Run("unsigned integers", func(t *testing.T) {
		assert.Equal(t, "42", ToString(uint(42)))
		assert.Equal(t, "8", ToString(uint8(8)))
		assert.Equal(t, "16", ToString(uint16(16)))
		assert.Equal(t, "32", ToString(uint32(32)))
		assert.Equal(t, "64", ToString(uint64(64)))
	})

	t.Run("floats", func(t *testing.T) {
		assert.Equal(t, "3.14", ToString(3.14))
		assert.Equal(t, "3.14159", ToString(float32(3.14159)))
	})

	t.Run("bool", func(t *testing.T) {
		assert.Equal(t, "true", ToString(true))
		assert.Equal(t, "false", ToString(false))
	})

	t.Run("time.Time", func(t *testing.T) {
		now := time.Now()
		result := ToString(now)
		assert.NotEmpty(t, result)
		assert.Contains(t, result, now.Format("2006"))
	})

	t.Run("nil", func(t *testing.T) {
		assert.Equal(t, "", ToString(nil))
	})

	t.Run("other types", func(t *testing.T) {
		assert.Equal(t, "[1 2 3]", ToString([]int{1, 2, 3}))
	})
}

func TestToInt(t *testing.T) {
	t.Run("integers", func(t *testing.T) {
		result, err := ToInt(42)
		assert.NoError(t, err)
		assert.Equal(t, int64(42), result)

		result, err = ToInt(-42)
		assert.NoError(t, err)
		assert.Equal(t, int64(-42), result)
	})

	t.Run("unsigned integers", func(t *testing.T) {
		result, err := ToInt(uint(42))
		assert.NoError(t, err)
		assert.Equal(t, int64(42), result)
	})

	t.Run("floats", func(t *testing.T) {
		result, err := ToInt(3.14)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), result)
	})

	t.Run("string", func(t *testing.T) {
		result, err := ToInt("42")
		assert.NoError(t, err)
		assert.Equal(t, int64(42), result)

		result, err = ToInt("-42")
		assert.NoError(t, err)
		assert.Equal(t, int64(-42), result)

		_, err = ToInt("invalid")
		assert.Error(t, err)
	})

	t.Run("bool", func(t *testing.T) {
		result, err := ToInt(true)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), result)

		result, err = ToInt(false)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), result)
	})

	t.Run("nil", func(t *testing.T) {
		result, err := ToInt(nil)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), result)
	})

	t.Run("unsupported type", func(t *testing.T) {
		_, err := ToInt([]int{1, 2, 3})
		assert.Error(t, err)
	})
}

func TestToInt64(t *testing.T) {
	t.Run("valid conversion", func(t *testing.T) {
		assert.Equal(t, int64(42), ToInt64(42))
		assert.Equal(t, int64(42), ToInt64("42"))
	})

	t.Run("invalid conversion returns 0", func(t *testing.T) {
		assert.Equal(t, int64(0), ToInt64("invalid"))
		assert.Equal(t, int64(0), ToInt64([]int{}))
	})
}

func TestToFloat64(t *testing.T) {
	t.Run("integers", func(t *testing.T) {
		result, err := ToFloat64(42)
		assert.NoError(t, err)
		assert.Equal(t, 42.0, result)
	})

	t.Run("floats", func(t *testing.T) {
		result, err := ToFloat64(3.14)
		assert.NoError(t, err)
		assert.Equal(t, 3.14, result)

		result, err = ToFloat64(float32(3.14))
		assert.NoError(t, err)
		assert.InDelta(t, 3.14, result, 0.01)
	})

	t.Run("string", func(t *testing.T) {
		result, err := ToFloat64("3.14")
		assert.NoError(t, err)
		assert.Equal(t, 3.14, result)

		_, err = ToFloat64("invalid")
		assert.Error(t, err)
	})

	t.Run("bool", func(t *testing.T) {
		result, err := ToFloat64(true)
		assert.NoError(t, err)
		assert.Equal(t, 1.0, result)

		result, err = ToFloat64(false)
		assert.NoError(t, err)
		assert.Equal(t, 0.0, result)
	})

	t.Run("nil", func(t *testing.T) {
		result, err := ToFloat64(nil)
		assert.NoError(t, err)
		assert.Equal(t, 0.0, result)
	})

	t.Run("unsupported type", func(t *testing.T) {
		_, err := ToFloat64([]int{1, 2, 3})
		assert.Error(t, err)
	})
}

func TestToFloat64Default(t *testing.T) {
	t.Run("valid conversion", func(t *testing.T) {
		assert.Equal(t, 3.14, ToFloat64Default(3.14))
		assert.Equal(t, 3.14, ToFloat64Default("3.14"))
	})

	t.Run("invalid conversion returns 0", func(t *testing.T) {
		assert.Equal(t, 0.0, ToFloat64Default("invalid"))
		assert.Equal(t, 0.0, ToFloat64Default([]int{}))
	})
}

func TestToBool(t *testing.T) {
	t.Run("bool", func(t *testing.T) {
		result, err := ToBool(true)
		assert.NoError(t, err)
		assert.True(t, result)

		result, err = ToBool(false)
		assert.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("string", func(t *testing.T) {
		result, err := ToBool("true")
		assert.NoError(t, err)
		assert.True(t, result)

		result, err = ToBool("false")
		assert.NoError(t, err)
		assert.False(t, result)

		result, err = ToBool("1")
		assert.NoError(t, err)
		assert.True(t, result)

		result, err = ToBool("0")
		assert.NoError(t, err)
		assert.False(t, result)

		_, err = ToBool("invalid")
		assert.Error(t, err)
	})

	t.Run("integers", func(t *testing.T) {
		result, err := ToBool(1)
		assert.NoError(t, err)
		assert.True(t, result)

		result, err = ToBool(0)
		assert.NoError(t, err)
		assert.False(t, result)

		result, err = ToBool(-1)
		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("floats", func(t *testing.T) {
		result, err := ToBool(1.5)
		assert.NoError(t, err)
		assert.True(t, result)

		result, err = ToBool(0.0)
		assert.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("nil", func(t *testing.T) {
		result, err := ToBool(nil)
		assert.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("unsupported type", func(t *testing.T) {
		_, err := ToBool([]int{1, 2, 3})
		assert.Error(t, err)
	})
}

func TestToBoolDefault(t *testing.T) {
	t.Run("valid conversion", func(t *testing.T) {
		assert.True(t, ToBoolDefault(true))
		assert.True(t, ToBoolDefault("true"))
		assert.True(t, ToBoolDefault(1))
	})

	t.Run("invalid conversion returns false", func(t *testing.T) {
		assert.False(t, ToBoolDefault("invalid"))
		assert.False(t, ToBoolDefault([]int{}))
	})
}

func TestToTime(t *testing.T) {
	t.Run("time.Time", func(t *testing.T) {
		now := time.Now()
		result, err := ToTime(now)
		assert.NoError(t, err)
		assert.Equal(t, now, result)

		ptr := &now
		result, err = ToTime(ptr)
		assert.NoError(t, err)
		assert.Equal(t, now, result)

		var nilPtr *time.Time
		result, err = ToTime(nilPtr)
		assert.NoError(t, err)
		assert.Equal(t, time.Time{}, result)
	})

	t.Run("RFC3339 format", func(t *testing.T) {
		str := "2023-12-25T10:30:00Z"
		result, err := ToTime(str)
		assert.NoError(t, err)
		assert.False(t, result.IsZero())
	})

	t.Run("common formats", func(t *testing.T) {
		str := "2023-12-25 10:30:00"
		result, err := ToTime(str)
		assert.NoError(t, err)
		assert.False(t, result.IsZero())

		str = "2023-12-25"
		result, err = ToTime(str)
		assert.NoError(t, err)
		assert.False(t, result.IsZero())
	})

	t.Run("unix timestamp", func(t *testing.T) {
		timestamp := int64(1703505600)
		result, err := ToTime(timestamp)
		assert.NoError(t, err)
		assert.False(t, result.IsZero())
	})

	t.Run("float timestamp", func(t *testing.T) {
		timestamp := 1703505600.5
		result, err := ToTime(timestamp)
		assert.NoError(t, err)
		assert.False(t, result.IsZero())
	})

	t.Run("invalid string", func(t *testing.T) {
		_, err := ToTime("invalid")
		assert.Error(t, err)
	})

	t.Run("nil", func(t *testing.T) {
		result, err := ToTime(nil)
		assert.NoError(t, err)
		assert.Equal(t, time.Time{}, result)
	})

	t.Run("unsupported type", func(t *testing.T) {
		_, err := ToTime([]int{1, 2, 3})
		assert.Error(t, err)
	})
}

func TestToTimeDefault(t *testing.T) {
	t.Run("valid conversion", func(t *testing.T) {
		now := time.Now()
		assert.Equal(t, now, ToTimeDefault(now))
		assert.Equal(t, now, ToTimeDefault(&now))
	})

	t.Run("invalid conversion returns zero time", func(t *testing.T) {
		assert.Equal(t, time.Time{}, ToTimeDefault("invalid"))
		assert.Equal(t, time.Time{}, ToTimeDefault([]int{}))
	})
}
