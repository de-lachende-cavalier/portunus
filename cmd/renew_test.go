package cmd

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/de-lachende-cavalier/portunus/pkg/config"
	"github.com/de-lachende-cavalier/portunus/pkg/testutil"
	"github.com/spf13/cobra"
)

// TestRenewCmd tests key renewal.
func TestRenewCmd(t *testing.T) {
	// Set up test environment
	tempDir, configPath := setupTestEnvironment(t)
	sshDir := filepath.Join(tempDir, ".ssh")

	// Create test key files
	key1, _ := testutil.CreateTestKeyPair(t, sshDir, "id_ed25519")
	key2, _ := testutil.CreateTestKeyPair(t, sshDir, "id_rsa")

	// Initialize the config with keys
	now := time.Now()
	creationTime := now.Add(-24 * time.Hour)
	expirationTime := now.Add(-1 * time.Hour) // Expired 1 hour ago

	cfgFile = configPath
	appConfig = &config.Config{
		Keys: map[string]config.KeyConfig{
			key1: {
				CreatedAt: creationTime,
				ExpiresAt: expirationTime,
			},
			key2: {
				CreatedAt: creationTime,
				ExpiresAt: expirationTime,
			},
		},
	}

	// Save the config
	if err := appConfig.Save(configPath); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Set up command flags
	renewTime = "1h"
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
		// Check if the creation time is unchanged
		if !keyConfig1.CreatedAt.Equal(creationTime) {
			t.Errorf("Expected creation time to be unchanged: %v, got %v", creationTime, keyConfig1.CreatedAt)
		}

		// Check if the expiration time was updated correctly
		expectedExpiry := now.Add(1 * time.Hour).Round(time.Second)
		if !keyConfig1.ExpiresAt.Round(time.Second).Equal(expectedExpiry) {
			t.Errorf("Expected expiration time %v, got %v", expectedExpiry, keyConfig1.ExpiresAt)
		}
	}

	keyConfig2, ok := loadedConfig.Keys[key2]
	if !ok {
		t.Errorf("Expected key %s to be in config", key2)
	} else {
		// Check if the creation time is unchanged
		if !keyConfig2.CreatedAt.Equal(creationTime) {
			t.Errorf("Expected creation time to be unchanged: %v, got %v", creationTime, keyConfig2.CreatedAt)
		}

		// Check if the expiration time was updated correctly
		expectedExpiry := now.Add(1 * time.Hour).Round(time.Second)
		if !keyConfig2.ExpiresAt.Round(time.Second).Equal(expectedExpiry) {
			t.Errorf("Expected expiration time %v, got %v", expectedExpiry, keyConfig2.ExpiresAt)
		}
	}
}
