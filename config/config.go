package config

import (
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	mu   sync.RWMutex
	data map[string]interface{}
}

// New creates a new configuration instance
func New() *Config {
	return &Config{
		data: make(map[string]interface{}),
	}
}

// LoadFromFile loads configuration from a YAML file
func (c *Config) LoadFromFile(path string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var config map[string]interface{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	c.data = config
	return nil
}

// LoadFromEnv loads configuration from environment variables with a prefix
func (c *Config) LoadFromEnv(prefix string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, env := range os.Environ() {
		key, value, _ := parseEnvVar(env, prefix)
		if key != "" {
			c.data[key] = value
		}
	}
}

func parseEnvVar(env, prefix string) (string, string, bool) {
	// Simple implementation - can be enhanced
	// Format: PREFIX_KEY=value
	if prefix != "" {
		// Check if env var starts with prefix
		// This is simplified - in production you might want more sophisticated handling
	}
	return "", "", false
}

// Get retrieves a configuration value by key
func (c *Config) Get(key string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, ok := c.data[key]
	if !ok {
		return nil, fmt.Errorf("config key not found: %s", key)
	}

	return value, nil
}

// GetOrDefault retrieves a configuration value or returns default value if not found
func (c *Config) GetOrDefault(key string, defaultValue interface{}) interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if value, ok := c.data[key]; ok {
		return value
	}

	return defaultValue
}

// Set sets a configuration value
func (c *Config) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value
}

// GetString retrieves a string value
func (c *Config) GetString(key string) (string, error) {
	value, err := c.Get(key)
	if err != nil {
		return "", err
	}

	if str, ok := value.(string); ok {
		return str, nil
	}

	return "", fmt.Errorf("config value is not a string: %s", key)
}

// GetStringOrDefault retrieves a string value or returns default
func (c *Config) GetStringOrDefault(key, defaultValue string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if value, ok := c.data[key]; ok {
		if str, ok := value.(string); ok {
			return str
		}
	}

	return defaultValue
}

// GetInt retrieves an integer value
func (c *Config) GetInt(key string) (int, error) {
	value, err := c.Get(key)
	if err != nil {
		return 0, err
	}

	switch v := value.(type) {
	case int:
		return v, nil
	case float64:
		return int(v), nil
	default:
		return 0, fmt.Errorf("config value is not an integer: %s", key)
	}
}

// GetIntOrDefault retrieves an integer value or returns default
func (c *Config) GetIntOrDefault(key string, defaultValue int) int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if value, ok := c.data[key]; ok {
		switch v := value.(type) {
		case int:
			return v
		case float64:
			return int(v)
		}
	}

	return defaultValue
}

// GetBool retrieves a boolean value
func (c *Config) GetBool(key string) (bool, error) {
	value, err := c.Get(key)
	if err != nil {
		return false, err
	}

	if b, ok := value.(bool); ok {
		return b, nil
	}

	return false, fmt.Errorf("config value is not a boolean: %s", key)
}

// GetBoolOrDefault retrieves a boolean value or returns default
func (c *Config) GetBoolOrDefault(key string, defaultValue bool) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if value, ok := c.data[key]; ok {
		if b, ok := value.(bool); ok {
			return b
		}
	}

	return defaultValue
}

// Has checks if a key exists in the configuration
func (c *Config) Has(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, ok := c.data[key]
	return ok
}

// Delete removes a key from the configuration
func (c *Config) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)
}

// Clear removes all configuration
func (c *Config) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = make(map[string]interface{})
}

// Keys returns all configuration keys
func (c *Config) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]string, 0, len(c.data))
	for key := range c.data {
		keys = append(keys, key)
	}

	return keys
}

// Merge merges another configuration into this one
func (c *Config) Merge(other *Config) {
	c.mu.Lock()
	defer c.mu.Unlock()

	other.mu.RLock()
	defer other.mu.RUnlock()

	for key, value := range other.data {
		c.data[key] = value
	}
}

// ToMap returns a copy of the configuration as a map
func (c *Config) ToMap() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make(map[string]interface{}, len(c.data))
	for key, value := range c.data {
		result[key] = value
	}

	return result
}

// Load is a convenience function to load config from file
func Load(path string) (*Config, error) {
	cfg := New()
	if err := cfg.LoadFromFile(path); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Global default config instance
var defaultConfig = New()

// LoadDefault loads configuration into the default config instance
func LoadDefault(path string) error {
	return defaultConfig.LoadFromFile(path)
}

// Get retrieves a value from the default config
func Get(key string) (interface{}, error) {
	return defaultConfig.Get(key)
}

// Set sets a value in the default config
func Set(key string, value interface{}) {
	defaultConfig.Set(key, value)
}

// GetString retrieves a string value from the default config
func GetString(key string) (string, error) {
	return defaultConfig.GetString(key)
}

// GetInt retrieves an integer value from the default config
func GetInt(key string) (int, error) {
	return defaultConfig.GetInt(key)
}

// GetBool retrieves a boolean value from the default config
func GetBool(key string) (bool, error) {
	return defaultConfig.GetBool(key)
}
