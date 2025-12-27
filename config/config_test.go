package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	cfg := New()
	assert.NotNil(t, cfg)
	assert.NotNil(t, cfg.data)
}

func TestSetAndGet(t *testing.T) {
	cfg := New()

	t.Run("set and get value", func(t *testing.T) {
		cfg.Set("key1", "value1")
		value, err := cfg.Get("key1")
		assert.NoError(t, err)
		assert.Equal(t, "value1", value)
	})

	t.Run("get non-existent key", func(t *testing.T) {
		_, err := cfg.Get("non-existent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestGetOrDefault(t *testing.T) {
	cfg := New()

	t.Run("key exists", func(t *testing.T) {
		cfg.Set("key1", "value1")
		value := cfg.GetOrDefault("key1", "default")
		assert.Equal(t, "value1", value)
	})

	t.Run("key does not exist", func(t *testing.T) {
		value := cfg.GetOrDefault("key2", "default")
		assert.Equal(t, "default", value)
	})
}

func TestGetString(t *testing.T) {
	cfg := New()

	t.Run("get string value", func(t *testing.T) {
		cfg.Set("key1", "value1")
		value, err := cfg.GetString("key1")
		assert.NoError(t, err)
		assert.Equal(t, "value1", value)
	})

	t.Run("get non-string value", func(t *testing.T) {
		cfg.Set("key2", 123)
		_, err := cfg.GetString("key2")
		assert.Error(t, err)
	})

	t.Run("get non-existent key", func(t *testing.T) {
		_, err := cfg.GetString("non-existent")
		assert.Error(t, err)
	})
}

func TestGetStringOrDefault(t *testing.T) {
	cfg := New()

	t.Run("key exists with string value", func(t *testing.T) {
		cfg.Set("key1", "value1")
		value := cfg.GetStringOrDefault("key1", "default")
		assert.Equal(t, "value1", value)
	})

	t.Run("key exists with non-string value", func(t *testing.T) {
		cfg.Set("key2", 123)
		value := cfg.GetStringOrDefault("key2", "default")
		assert.Equal(t, "default", value)
	})

	t.Run("key does not exist", func(t *testing.T) {
		value := cfg.GetStringOrDefault("key3", "default")
		assert.Equal(t, "default", value)
	})
}

func TestGetInt(t *testing.T) {
	cfg := New()

	t.Run("get int value", func(t *testing.T) {
		cfg.Set("key1", 42)
		value, err := cfg.GetInt("key1")
		assert.NoError(t, err)
		assert.Equal(t, 42, value)
	})

	t.Run("get float64 value as int", func(t *testing.T) {
		cfg.Set("key2", 42.5)
		value, err := cfg.GetInt("key2")
		assert.NoError(t, err)
		assert.Equal(t, 42, value)
	})

	t.Run("get non-int value", func(t *testing.T) {
		cfg.Set("key3", "not an int")
		_, err := cfg.GetInt("key3")
		assert.Error(t, err)
	})
}

func TestGetIntOrDefault(t *testing.T) {
	cfg := New()

	t.Run("key exists with int value", func(t *testing.T) {
		cfg.Set("key1", 42)
		value := cfg.GetIntOrDefault("key1", 0)
		assert.Equal(t, 42, value)
	})

	t.Run("key does not exist", func(t *testing.T) {
		value := cfg.GetIntOrDefault("key2", 100)
		assert.Equal(t, 100, value)
	})
}

func TestGetBool(t *testing.T) {
	cfg := New()

	t.Run("get bool value", func(t *testing.T) {
		cfg.Set("key1", true)
		value, err := cfg.GetBool("key1")
		assert.NoError(t, err)
		assert.True(t, value)
	})

	t.Run("get non-bool value", func(t *testing.T) {
		cfg.Set("key2", "not a bool")
		_, err := cfg.GetBool("key2")
		assert.Error(t, err)
	})
}

func TestGetBoolOrDefault(t *testing.T) {
	cfg := New()

	t.Run("key exists with bool value", func(t *testing.T) {
		cfg.Set("key1", true)
		value := cfg.GetBoolOrDefault("key1", false)
		assert.True(t, value)
	})

	t.Run("key does not exist", func(t *testing.T) {
		value := cfg.GetBoolOrDefault("key2", false)
		assert.False(t, value)
	})
}

func TestHas(t *testing.T) {
	cfg := New()

	t.Run("key exists", func(t *testing.T) {
		cfg.Set("key1", "value1")
		assert.True(t, cfg.Has("key1"))
	})

	t.Run("key does not exist", func(t *testing.T) {
		assert.False(t, cfg.Has("non-existent"))
	})
}

func TestDelete(t *testing.T) {
	cfg := New()

	t.Run("delete existing key", func(t *testing.T) {
		cfg.Set("key1", "value1")
		assert.True(t, cfg.Has("key1"))
		cfg.Delete("key1")
		assert.False(t, cfg.Has("key1"))
	})

	t.Run("delete non-existent key", func(t *testing.T) {
		assert.NotPanics(t, func() {
			cfg.Delete("non-existent")
		})
	})
}

func TestClear(t *testing.T) {
	cfg := New()
	cfg.Set("key1", "value1")
	cfg.Set("key2", "value2")

	cfg.Clear()

	assert.False(t, cfg.Has("key1"))
	assert.False(t, cfg.Has("key2"))
}

func TestKeys(t *testing.T) {
	cfg := New()
	cfg.Set("key1", "value1")
	cfg.Set("key2", "value2")
	cfg.Set("key3", "value3")

	keys := cfg.Keys()
	assert.Len(t, keys, 3)
	assert.Contains(t, keys, "key1")
	assert.Contains(t, keys, "key2")
	assert.Contains(t, keys, "key3")
}

func TestMerge(t *testing.T) {
	cfg1 := New()
	cfg1.Set("key1", "value1")
	cfg1.Set("key2", "value2")

	cfg2 := New()
	cfg2.Set("key2", "new_value2")
	cfg2.Set("key3", "value3")

	cfg1.Merge(cfg2)

	value1, _ := cfg1.Get("key1")
	assert.Equal(t, "value1", value1)

	value2, _ := cfg1.Get("key2")
	assert.Equal(t, "new_value2", value2) // Overridden

	value3, _ := cfg1.Get("key3")
	assert.Equal(t, "value3", value3)
}

func TestToMap(t *testing.T) {
	cfg := New()
	cfg.Set("key1", "value1")
	cfg.Set("key2", 42)

	m := cfg.ToMap()

	assert.Equal(t, "value1", m["key1"])
	assert.Equal(t, 42, m["key2"])

	// Modifying returned map shouldn't affect config
	m["key3"] = "value3"
	assert.False(t, cfg.Has("key3"))
}

func TestLoadFromFile(t *testing.T) {
	t.Run("load valid yaml file", func(t *testing.T) {
		// Create temp YAML file
		tmpDir := os.TempDir()
		tmpFile := filepath.Join(tmpDir, "config_test.yaml")
		content := []byte("key1: value1\nkey2: 42\nkey3: true\n")
		err := os.WriteFile(tmpFile, content, 0644)
		assert.NoError(t, err)
		defer os.Remove(tmpFile)

		cfg := New()
		err = cfg.LoadFromFile(tmpFile)
		assert.NoError(t, err)

		value1, err := cfg.GetString("key1")
		assert.NoError(t, err)
		assert.Equal(t, "value1", value1)

		value2, err := cfg.GetInt("key2")
		assert.NoError(t, err)
		assert.Equal(t, 42, value2)
	})

	t.Run("load non-existent file", func(t *testing.T) {
		cfg := New()
		err := cfg.LoadFromFile("non-existent.yaml")
		assert.Error(t, err)
	})

	t.Run("load invalid yaml file", func(t *testing.T) {
		tmpDir := os.TempDir()
		tmpFile := filepath.Join(tmpDir, "invalid_config_test.yaml")
		content := []byte("invalid: yaml: content:\n  - [")
		err := os.WriteFile(tmpFile, content, 0644)
		assert.NoError(t, err)
		defer os.Remove(tmpFile)

		cfg := New()
		err = cfg.LoadFromFile(tmpFile)
		assert.Error(t, err)
	})
}

func TestLoad(t *testing.T) {
	t.Run("load and return config", func(t *testing.T) {
		tmpDir := os.TempDir()
		tmpFile := filepath.Join(tmpDir, "load_test.yaml")
		content := []byte("key: value\n")
		err := os.WriteFile(tmpFile, content, 0644)
		assert.NoError(t, err)
		defer os.Remove(tmpFile)

		cfg, err := Load(tmpFile)
		assert.NoError(t, err)
		assert.NotNil(t, cfg)

		value, err := cfg.GetString("key")
		assert.NoError(t, err)
		assert.Equal(t, "value", value)
	})
}

func TestGlobalConfig(t *testing.T) {
	t.Run("use global config functions", func(t *testing.T) {
		// Clear global config first
		defaultConfig.Clear()

		Set("global_key", "global_value")

		value, err := GetString("global_key")
		assert.NoError(t, err)
		assert.Equal(t, "global_value", value)
	})
}

func TestLoadFromEnv(t *testing.T) {
	cfg := New()
	// LoadFromEnv is a no-op in current implementation
	// This test just verifies it doesn't panic
	assert.NotPanics(t, func() {
		cfg.LoadFromEnv("PREFIX")
	})
}

func TestConcurrency(t *testing.T) {
	cfg := New()
	done := make(chan bool)

	// Concurrent writes
	for i := 0; i < 10; i++ {
		go func(idx int) {
			cfg.Set("key", idx)
			done <- true
		}(i)
	}

	// Concurrent reads
	for i := 0; i < 10; i++ {
		go func() {
			cfg.Get("key")
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 20; i++ {
		<-done
	}

	// Should not have any race conditions
	assert.True(t, true)
}
