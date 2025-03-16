package cmd

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/de-lachende-cavalier/portunus/pkg/config"
	"github.com/de-lachende-cavalier/portunus/pkg/testutil"
	"github.com/spf13/cobra"
)

// setupTestEnvironment sets up a test environment for the commands
func setupTestEnvironment(t *testing.T) (string, string) {
	// Create a test directory
	tempDir := testutil.TempDir(t)

	// Create a test SSH directory
	sshDir := filepath.Join(tempDir, ".ssh")
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		t.Fatalf("Failed to create test SSH directory: %v", err)
	}

	// Create a test config file
	configPath := filepath.Join(tempDir, ".portunus.json")

	// Set environment variables
	originalHome := os.Getenv("HOME")
	t.Cleanup(func() {
		os.Setenv("HOME", originalHome)
	})
	os.Setenv("HOME", tempDir)

	return tempDir, configPath
}

// TestRotateCommand tests the rotate command
func TestRotateCommand(t *testing.T) {
	// Skip this test if ssh-keygen is not available
	if _, err := os.Stat("/usr/bin/ssh-keygen"); os.IsNotExist(err) {
		t.Skip("ssh-keygen not available, skipping test")
	}

	// Set up test environment
	tempDir, configPath := setupTestEnvironment(t)
	sshDir := filepath.Join(tempDir, ".ssh")

	// Create test key files
	key1, _ := testutil.CreateTestKeyPair(t, sshDir, "id_ed25519")
	key2, _ := testutil.CreateTestKeyPair(t, sshDir, "id_rsa")

	// Initialize the config
	cfgFile = configPath
	appConfig = &config.Config{
		Keys: make(map[string]config.KeyConfig),
	}

	// Set up command flags
	rotateCipher = "ed25519"
	rotateTime = "24h"
	rotatePassword = "test"
	rotateKeySubset = []string{key1, key2}

	// Initialize the context
	rootContext = context.Background()

	// Create a mock command for testing
	mockCmd := &cobra.Command{Use: "test"}

	// Run the rotate command
	runRotateCmd(mockCmd, nil)

	// Check if the keys were rotated
	testutil.AssertFileExists(t, key1)
	testutil.AssertFileExists(t, key1+".pub")
	testutil.AssertFileExists(t, key2)
	testutil.AssertFileExists(t, key2+".pub")

	// Check if the config was updated
	loadedConfig, err := config.Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if len(loadedConfig.Keys) != 2 {
		t.Errorf("Expected 2 keys in config, got %d", len(loadedConfig.Keys))
	}

	// Check if the keys are in the config
	_, ok := loadedConfig.Keys[key1]
	if !ok {
		t.Errorf("Expected key %s to be in config", key1)
	}

	_, ok = loadedConfig.Keys[key2]
	if !ok {
		t.Errorf("Expected key %s to be in config", key2)
	}
}

// TestRenewCommand tests the renew command
func TestRenewCommand(t *testing.T) {
	// Set up test environment
	tempDir, configPath := setupTestEnvironment(t)
	sshDir := filepath.Join(tempDir, ".ssh")

	// Create test key files
	key1, _ := testutil.CreateTestKeyPair(t, sshDir, "id_ed25519")
	key2, _ := testutil.CreateTestKeyPair(t, sshDir, "id_rsa")

	// Initialize the config with expired keys
	now := time.Now()
	cfgFile = configPath
	appConfig = &config.Config{
		Keys: map[string]config.KeyConfig{
			key1: {
				CreatedAt: now.Add(-48 * time.Hour),
				ExpiresAt: now.Add(-24 * time.Hour),
			},
			key2: {
				CreatedAt: now.Add(-48 * time.Hour),
				ExpiresAt: now.Add(-24 * time.Hour),
			},
		},
	}

	// Save the config
	if err := appConfig.Save(configPath); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Set up command flags
	renewTime = "48h"
	renewKeySubset = []string{}

	// Initialize the context
	rootContext = context.Background()

	// Create a mock command for testing
	mockCmd := &cobra.Command{Use: "test"}

	// Run the renew command
	runRenewCmd(mockCmd, nil)

	// Check if the config was updated
	loadedConfig, err := config.Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Check if the keys are in the config
	keyConfig1, ok := loadedConfig.Keys[key1]
	if !ok {
		t.Errorf("Expected key %s to be in config", key1)
	} else {
		// Check if the expiration time was updated
		if keyConfig1.ExpiresAt.Before(now) {
			t.Errorf("Expected key %s to have future expiration time, got %v", key1, keyConfig1.ExpiresAt)
		}
	}

	keyConfig2, ok := loadedConfig.Keys[key2]
	if !ok {
		t.Errorf("Expected key %s to be in config", key2)
	} else {
		// Check if the expiration time was updated
		if keyConfig2.ExpiresAt.Before(now) {
			t.Errorf("Expected key %s to have future expiration time, got %v", key2, keyConfig2.ExpiresAt)
		}
	}
}

// TestCheckCommand tests the check command
func TestCheckCommand(t *testing.T) {
	// Set up test environment
	tempDir, configPath := setupTestEnvironment(t)
	sshDir := filepath.Join(tempDir, ".ssh")

	// Create test key files
	key1, _ := testutil.CreateTestKeyPair(t, sshDir, "id_ed25519")
	key2, _ := testutil.CreateTestKeyPair(t, sshDir, "id_rsa")

	// Initialize the config with expired and non-expired keys
	now := time.Now()
	cfgFile = configPath
	appConfig = &config.Config{
		Keys: map[string]config.KeyConfig{
			key1: {
				CreatedAt: now.Add(-48 * time.Hour),
				ExpiresAt: now.Add(-24 * time.Hour),
			},
			key2: {
				CreatedAt: now.Add(-24 * time.Hour),
				ExpiresAt: now.Add(24 * time.Hour),
			},
		},
	}

	// Save the config
	if err := appConfig.Save(configPath); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Initialize the context
	rootContext = context.Background()

	// Create a mock command for testing
	mockCmd := &cobra.Command{Use: "test"}

	// Run the check command
	// Note: We can't easily test the output, but at least ensure it doesn't crash
	runCheckCmd(mockCmd, nil)
}
