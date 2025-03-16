package cmd

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/de-lachende-cavalier/portunus/pkg/config"
	"github.com/de-lachende-cavalier/portunus/pkg/testutil"
	"github.com/spf13/cobra"
)

// Helper function, check whether the correct cipher has been used in overwriting the keys.
func usedCorrectCipher(keyFile string, cipher string) bool {
	bytes, err := os.ReadFile(keyFile + ".pub")
	if err != nil {
		return false
	}
	text := string(bytes)

	// ssh-<cipher> ...
	return strings.Contains(text, "ssh-"+cipher)
}

// TestRotateCmd_RSA tests key rotation for RSA.
func TestRotateCmd_RSA(t *testing.T) {
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
	rotateCipher = "rsa"
	rotateTime = "30m"
	rotatePassword = "test"
	rotateKeySubset = []string{key1, key2}

	// Initialize the context
	rootContext = context.Background()

	// Create a mock command for testing
	mockCmd := &cobra.Command{Use: "test"}

	// Run the rotate command
	runRotateCmd(mockCmd, nil)

	// Check if the config was updated
	loadedConfig, err := config.Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if len(loadedConfig.Keys) != 2 {
		t.Errorf("Expected 2 keys in config, got %d", len(loadedConfig.Keys))
	}

	// Check if the keys are in the config and have correct expiration times
	for key, keyConfig := range loadedConfig.Keys {
		// Check expiration time
		expectedExpiry := keyConfig.CreatedAt.Add(30 * time.Minute)
		if !keyConfig.ExpiresAt.Equal(expectedExpiry) {
			t.Errorf("Expected expiration time %v, got %v", expectedExpiry, keyConfig.ExpiresAt)
		}

		// Check cipher
		if !usedCorrectCipher(key, "rsa") {
			t.Errorf("Expected RSA cipher for key %s", key)
		}
	}
}

// TestRotateCmd_Ed25519 tests key rotation for Ed25519.
func TestRotateCmd_Ed25519(t *testing.T) {
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
	rotateTime = "30m"
	rotatePassword = "test"
	rotateKeySubset = []string{key1, key2}

	// Initialize the context
	rootContext = context.Background()

	// Create a mock command for testing
	mockCmd := &cobra.Command{Use: "test"}

	// Run the rotate command
	runRotateCmd(mockCmd, nil)

	// Check if the config was updated
	loadedConfig, err := config.Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if len(loadedConfig.Keys) != 2 {
		t.Errorf("Expected 2 keys in config, got %d", len(loadedConfig.Keys))
	}

	// Check if the keys are in the config and have correct expiration times
	for key, keyConfig := range loadedConfig.Keys {
		// Check expiration time
		expectedExpiry := keyConfig.CreatedAt.Add(30 * time.Minute)
		if !keyConfig.ExpiresAt.Equal(expectedExpiry) {
			t.Errorf("Expected expiration time %v, got %v", expectedExpiry, keyConfig.ExpiresAt)
		}

		// Check cipher
		if !usedCorrectCipher(key, "ed25519") {
			t.Errorf("Expected Ed25519 cipher for key %s", key)
		}
	}
}
