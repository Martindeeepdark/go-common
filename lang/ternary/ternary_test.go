package ternary

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestT(t *testing.T) {
	t.Run("condition true returns trueValue", func(t *testing.T) {
		result := T(true, "yes", "no")
		assert.Equal(t, "yes", result)
	})

	t.Run("condition false returns falseValue", func(t *testing.T) {
		result := T(false, "yes", "no")
		assert.Equal(t, "no", result)
	})

	t.Run("with integers", func(t *testing.T) {
		result := T(5 > 3, 100, 200)
		assert.Equal(t, 100, result)
	})

	t.Run("with slices", func(t *testing.T) {
		trueSlice := []int{1, 2, 3}
		falseSlice := []int{4, 5, 6}
		result := T(true, trueSlice, falseSlice)
		assert.Equal(t, trueSlice, result)
	})
}

func TestLazy(t *testing.T) {
	t.Run("condition true calls trueValueFunc", func(t *testing.T) {
		calls := 0
		trueFunc := func() int {
			calls++
			return 100
		}
		falseFunc := func() int {
			calls++
			return 200
		}

		result := Lazy(true, trueFunc, falseFunc)
		assert.Equal(t, 100, result)
		assert.Equal(t, 1, calls)
	})

	t.Run("condition false calls falseValueFunc", func(t *testing.T) {
		trueFunc := func() string {
			return "yes"
		}
		falseFunc := func() string {
			return "no"
		}

		result := Lazy(false, trueFunc, falseFunc)
		assert.Equal(t, "no", result)
	})

	t.Run("only one branch is evaluated", func(t *testing.T) {
		trueCalls := 0
		falseCalls := 0

		trueFunc := func() int {
			trueCalls++
			return 1
		}
		falseFunc := func() int {
			falseCalls++
			return 2
		}

		Lazy(true, trueFunc, falseFunc)
		assert.Equal(t, 1, trueCalls)
		assert.Equal(t, 0, falseCalls)

		trueCalls = 0
		falseCalls = 0

		Lazy(false, trueFunc, falseFunc)
		assert.Equal(t, 0, trueCalls)
		assert.Equal(t, 1, falseCalls)
	})
}

func TestOr(t *testing.T) {
	t.Run("first is not zero", func(t *testing.T) {
		result := Or(5, 10)
		assert.Equal(t, 5, result)
	})

	t.Run("first is zero returns second", func(t *testing.T) {
		result := Or(0, 10)
		assert.Equal(t, 10, result)
	})

	t.Run("with strings", func(t *testing.T) {
		result := Or("", "default")
		assert.Equal(t, "default", result)
	})

	t.Run("with pointers", func(t *testing.T) {
		a, b := 5, 10
		result := Or(&a, &b)
		assert.Equal(t, &a, result)
	})
}

func TestOrFunc(t *testing.T) {
	t.Run("first is not zero", func(t *testing.T) {
		calls := 0
		result := OrFunc(5, func() int {
			calls++
			return 10
		})
		assert.Equal(t, 5, result)
		assert.Equal(t, 0, calls)
	})

	t.Run("first is zero calls func", func(t *testing.T) {
		calls := 0
		result := OrFunc(0, func() int {
			calls++
			return 10
		})
		assert.Equal(t, 10, result)
		assert.Equal(t, 1, calls)
	})
}

func TestAnd(t *testing.T) {
	t.Run("condition true returns first", func(t *testing.T) {
		result := And(true, 5)
		assert.Equal(t, 5, result)
	})

	t.Run("condition false returns zero", func(t *testing.T) {
		result := And(false, 5)
		assert.Equal(t, 0, result)
	})
}

func TestCoalesce(t *testing.T) {
	t.Run("first non-zero value", func(t *testing.T) {
		result := Coalesce(0, 0, 5, 10)
		assert.Equal(t, 5, result)
	})

	t.Run("all zero returns zero", func(t *testing.T) {
		result := Coalesce(0, 0, 0)
		assert.Equal(t, 0, result)
	})

	t.Run("first value is non-zero", func(t *testing.T) {
		result := Coalesce(5, 10, 15)
		assert.Equal(t, 5, result)
	})

	t.Run("no args returns zero", func(t *testing.T) {
		result := Coalesce[int]()
		assert.Equal(t, 0, result)
	})

	t.Run("with strings", func(t *testing.T) {
		result := Coalesce("", "", "default", "backup")
		assert.Equal(t, "default", result)
	})
}

func TestIsZero(t *testing.T) {
	t.Run("zero values", func(t *testing.T) {
		assert.True(t, IsZero(0))
		assert.True(t, IsZero(""))
		assert.True(t, IsZero(false))
	})

	t.Run("non-zero values", func(t *testing.T) {
		assert.False(t, IsZero(5))
		assert.False(t, IsZero("hello"))
		assert.False(t, IsZero(true))
	})
}

func TestIsNotZero(t *testing.T) {
	t.Run("non-zero values", func(t *testing.T) {
		assert.True(t, IsNotZero(5))
		assert.True(t, IsNotZero("hello"))
		assert.True(t, IsNotZero(true))
	})

	t.Run("zero values", func(t *testing.T) {
		assert.False(t, IsNotZero(0))
		assert.False(t, IsNotZero(""))
		assert.False(t, IsNotZero(false))
	})
}

func TestIf(t *testing.T) {
	t.Run("condition true returns value", func(t *testing.T) {
		result := If(true, 42)
		assert.Equal(t, 42, result)
	})

	t.Run("condition false returns zero", func(t *testing.T) {
		result := If(false, 42)
		assert.Equal(t, 0, result)
	})
}

func TestUnless(t *testing.T) {
	t.Run("condition false returns value", func(t *testing.T) {
		result := Unless(false, 42)
		assert.Equal(t, 42, result)
	})

	t.Run("condition true returns zero", func(t *testing.T) {
		result := Unless(true, 42)
		assert.Equal(t, 0, result)
	})
}

func TestSwitch(t *testing.T) {
	t.Run("matching key", func(t *testing.T) {
		cases := map[int]string{
			1: "one",
			2: "two",
			3: "three",
		}
		result := Switch(2, cases, "default")
		assert.Equal(t, "two", result)
	})

	t.Run("no matching key returns default", func(t *testing.T) {
		cases := map[int]string{
			1: "one",
			2: "two",
		}
		result := Switch(3, cases, "default")
		assert.Equal(t, "default", result)
	})

	t.Run("empty cases returns default", func(t *testing.T) {
		cases := map[int]string{}
		result := Switch(1, cases, "default")
		assert.Equal(t, "default", result)
	})
}

func TestDefault(t *testing.T) {
	t.Run("value is not zero returns value", func(t *testing.T) {
		result := Default(5, 10)
		assert.Equal(t, 5, result)
	})

	t.Run("value is zero returns default", func(t *testing.T) {
		result := Default(0, 10)
		assert.Equal(t, 10, result)
	})

	t.Run("with strings", func(t *testing.T) {
		result := Default("", "default")
		assert.Equal(t, "default", result)
	})

	t.Run("with pointers", func(t *testing.T) {
		a := 5
		result := Default(&a, &a)
		assert.Equal(t, &a, result)
	})
}
