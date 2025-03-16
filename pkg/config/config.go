// Package config provides configuration management for the application
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// KeyConfig represents the configuration for a key
type KeyConfig struct {
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// Config represents the application configuration
type Config struct {
	Keys map[string]KeyConfig `json:"keys"`
}

// DefaultConfigPath returns the default path for the config file
func DefaultConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ".portunus.json"
	}
	return filepath.Join(homeDir, ".portunus.json")
}

// Load loads the configuration from the specified path
func Load(path string) (*Config, error) {
	if path == "" {
		path = DefaultConfigPath()
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{
				Keys: make(map[string]KeyConfig),
			}, nil
		}
		return nil, err
	}

	var config Config
	if len(data) == 0 {
		return &Config{
			Keys: make(map[string]KeyConfig),
		}, nil
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// Save saves the configuration to the specified path
func (c *Config) Save(path string) error {
	if path == "" {
		path = DefaultConfigPath()
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

// AddKey adds a key to the configuration
func (c *Config) AddKey(path string, createdAt, expiresAt time.Time) {
	c.Keys[path] = KeyConfig{
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
	}
}

// RemoveKey removes a key from the configuration
func (c *Config) RemoveKey(path string) {
	delete(c.Keys, path)
}

// GetExpiredKeys returns a list of expired keys
func (c *Config) GetExpiredKeys() []string {
	var expired []string
	now := time.Now()

	for path, keyConfig := range c.Keys {
		if now.After(keyConfig.ExpiresAt) {
			expired = append(expired, path)
		}
	}

	return expired
}

// CleanNonExistentKeys removes keys that no longer exist from the configuration
func (c *Config) CleanNonExistentKeys() {
	for path := range c.Keys {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			delete(c.Keys, path)
		}
	}
}
