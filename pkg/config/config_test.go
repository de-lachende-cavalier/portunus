package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/de-lachende-cavalier/portunus/pkg/testutil"
)

func TestConfig_SaveAndLoad(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := testutil.TempDir(t)
	configPath := filepath.Join(tempDir, "config.json")

	// Create a test config
	now := time.Now().Round(time.Second)
	expiry := now.Add(24 * time.Hour).Round(time.Second)

	cfg := &Config{
		Keys: map[string]KeyConfig{
			"/home/user/.ssh/id_ed25519": {
				CreatedAt: now,
				ExpiresAt: expiry,
			},
		},
	}

	// Save the config
	err := cfg.Save(configPath)
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Load the config
	loadedCfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Check if the loaded config matches the original
	if len(loadedCfg.Keys) != len(cfg.Keys) {
		t.Errorf("Expected %d keys, got %d", len(cfg.Keys), len(loadedCfg.Keys))
	}

	for key, keyConfig := range cfg.Keys {
		loadedKeyConfig, ok := loadedCfg.Keys[key]
		if !ok {
			t.Errorf("Expected key %s to exist in loaded config", key)
			continue
		}

		if !keyConfig.CreatedAt.Equal(loadedKeyConfig.CreatedAt) {
			t.Errorf("Expected CreatedAt %v, got %v", keyConfig.CreatedAt, loadedKeyConfig.CreatedAt)
		}

		if !keyConfig.ExpiresAt.Equal(loadedKeyConfig.ExpiresAt) {
			t.Errorf("Expected ExpiresAt %v, got %v", keyConfig.ExpiresAt, loadedKeyConfig.ExpiresAt)
		}
	}
}

func TestConfig_AddAndRemoveKey(t *testing.T) {
	// Create a test config
	cfg := &Config{
		Keys: make(map[string]KeyConfig),
	}

	// Add a key
	now := time.Now().Round(time.Second)
	expiry := now.Add(24 * time.Hour).Round(time.Second)
	keyPath := "/home/user/.ssh/id_ed25519"

	cfg.AddKey(keyPath, now, expiry)

	// Check if the key was added
	if len(cfg.Keys) != 1 {
		t.Errorf("Expected 1 key, got %d", len(cfg.Keys))
	}

	keyConfig, ok := cfg.Keys[keyPath]
	if !ok {
		t.Errorf("Expected key %s to exist", keyPath)
	} else {
		if !keyConfig.CreatedAt.Equal(now) {
			t.Errorf("Expected CreatedAt %v, got %v", now, keyConfig.CreatedAt)
		}

		if !keyConfig.ExpiresAt.Equal(expiry) {
			t.Errorf("Expected ExpiresAt %v, got %v", expiry, keyConfig.ExpiresAt)
		}
	}

	// Remove the key
	cfg.RemoveKey(keyPath)

	// Check if the key was removed
	if len(cfg.Keys) != 0 {
		t.Errorf("Expected 0 keys, got %d", len(cfg.Keys))
	}

	_, ok = cfg.Keys[keyPath]
	if ok {
		t.Errorf("Expected key %s to not exist", keyPath)
	}
}

func TestConfig_GetExpiredKeys(t *testing.T) {
	// Create a test config
	cfg := &Config{
		Keys: make(map[string]KeyConfig),
	}

	now := time.Now()

	// Add an expired key
	expiredKeyPath := "/home/user/.ssh/expired_key"
	cfg.AddKey(expiredKeyPath, now.Add(-48*time.Hour), now.Add(-24*time.Hour))

	// Add a non-expired key
	validKeyPath := "/home/user/.ssh/valid_key"
	cfg.AddKey(validKeyPath, now.Add(-24*time.Hour), now.Add(24*time.Hour))

	// Get expired keys
	expiredKeys := cfg.GetExpiredKeys()

	// Check if only the expired key is returned
	if len(expiredKeys) != 1 {
		t.Errorf("Expected 1 expired key, got %d", len(expiredKeys))
	}

	if len(expiredKeys) > 0 && expiredKeys[0] != expiredKeyPath {
		t.Errorf("Expected expired key %s, got %s", expiredKeyPath, expiredKeys[0])
	}
}

func TestConfig_CleanNonExistentKeys(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := testutil.TempDir(t)

	// Create a test key file
	keyName := "test_key"
	keyPath, _ := testutil.CreateTestKeyPair(t, tempDir, keyName)

	// Create a test config with existing and non-existent keys
	cfg := &Config{
		Keys: map[string]KeyConfig{
			keyPath: {
				CreatedAt: time.Now().Add(-24 * time.Hour),
				ExpiresAt: time.Now().Add(24 * time.Hour),
			},
			filepath.Join(tempDir, "non_existent_key"): {
				CreatedAt: time.Now().Add(-24 * time.Hour),
				ExpiresAt: time.Now().Add(24 * time.Hour),
			},
		},
	}

	// Clean non-existent keys
	cfg.CleanNonExistentKeys()

	// Check if only the existing key remains
	if len(cfg.Keys) != 1 {
		t.Errorf("Expected 1 key, got %d", len(cfg.Keys))
	}

	_, ok := cfg.Keys[keyPath]
	if !ok {
		t.Errorf("Expected key %s to exist", keyPath)
	}

	_, ok = cfg.Keys[filepath.Join(tempDir, "non_existent_key")]
	if ok {
		t.Errorf("Expected non-existent key to be removed")
	}
}

func TestDefaultConfigPath(t *testing.T) {
	path := DefaultConfigPath()

	// Check if the path is not empty
	if path == "" {
		t.Error("Expected non-empty default config path")
	}

	// Check if the path contains .portunus.json
	if filepath.Base(path) != ".portunus.json" {
		t.Errorf("Expected filename to be .portunus.json, got %s", filepath.Base(path))
	}
}

func TestLoad_EmptyFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := testutil.TempDir(t)
	configPath := filepath.Join(tempDir, "empty_config.json")

	// Create an empty file
	err := os.WriteFile(configPath, []byte{}, 0600)
	if err != nil {
		t.Fatalf("Failed to create empty config file: %v", err)
	}

	// Load the config
	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load empty config: %v", err)
	}

	// Check if the config is initialized with empty keys
	if cfg.Keys == nil {
		t.Error("Expected Keys to be initialized")
	}

	if len(cfg.Keys) != 0 {
		t.Errorf("Expected 0 keys, got %d", len(cfg.Keys))
	}
}

func TestLoad_NonExistentFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := testutil.TempDir(t)
	configPath := filepath.Join(tempDir, "non_existent_config.json")

	// Load the config
	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load non-existent config: %v", err)
	}

	// Check if the config is initialized with empty keys
	if cfg.Keys == nil {
		t.Error("Expected Keys to be initialized")
	}

	if len(cfg.Keys) != 0 {
		t.Errorf("Expected 0 keys, got %d", len(cfg.Keys))
	}
}
