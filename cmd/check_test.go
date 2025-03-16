package cmd

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/de-lachende-cavalier/portunus/pkg/config"
	"github.com/de-lachende-cavalier/portunus/pkg/testutil"
	"github.com/spf13/cobra"
)

// captureOutput captures stdout for testing
func captureOutput(f func()) string {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// back to normal state
	w.Close()
	os.Stdout = old // restoring the real stdout
	return <-outC
}

// TestCheckCmd_AllExpired tests the check command when all keys have expired
func TestCheckCmd_AllExpired(t *testing.T) {
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

	// Initialize the context
	rootContext = context.Background()

	// Create a mock command for testing
	mockCmd := &cobra.Command{Use: "test"}

	// Capture output and run the check command
	output := captureOutput(func() {
		runCheckCmd(mockCmd, nil)
	})

	// Check if the output indicates expired keys
	if !strings.Contains(output, "The following keys have expired") {
		t.Errorf("Expected output to indicate expired keys, got: %s", output)
	}

	// Check if both keys are mentioned in the output
	if !strings.Contains(output, key1) || !strings.Contains(output, key2) {
		t.Errorf("Expected output to mention both expired keys, got: %s", output)
	}
}

// TestCheckCmd_NoneExpired tests the check command when no keys have expired
func TestCheckCmd_NoneExpired(t *testing.T) {
	// Set up test environment
	tempDir, configPath := setupTestEnvironment(t)
	sshDir := filepath.Join(tempDir, ".ssh")

	// Create test key files
	key1, _ := testutil.CreateTestKeyPair(t, sshDir, "id_ed25519")
	key2, _ := testutil.CreateTestKeyPair(t, sshDir, "id_rsa")

	// Initialize the config with non-expired keys
	now := time.Now()
	cfgFile = configPath
	appConfig = &config.Config{
		Keys: map[string]config.KeyConfig{
			key1: {
				CreatedAt: now.Add(-24 * time.Hour),
				ExpiresAt: now.Add(24 * time.Hour),
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

	// Capture output and run the check command
	output := captureOutput(func() {
		runCheckCmd(mockCmd, nil)
	})

	// Check if the output indicates no expired keys
	if !strings.Contains(output, "No expired keys found") {
		t.Errorf("Expected output to indicate no expired keys, got: %s", output)
	}
}

// TestCheckCmd_SomeExpired tests the check command when some keys have expired
func TestCheckCmd_SomeExpired(t *testing.T) {
	// Set up test environment
	tempDir, configPath := setupTestEnvironment(t)
	sshDir := filepath.Join(tempDir, ".ssh")

	// Create test key files
	key1, _ := testutil.CreateTestKeyPair(t, sshDir, "id_ed25519")
	key2, _ := testutil.CreateTestKeyPair(t, sshDir, "id_rsa")

	// Initialize the config with some expired keys
	now := time.Now()
	cfgFile = configPath
	appConfig = &config.Config{
		Keys: map[string]config.KeyConfig{
			key1: {
				CreatedAt: now.Add(-48 * time.Hour),
				ExpiresAt: now.Add(-24 * time.Hour), // Expired
			},
			key2: {
				CreatedAt: now.Add(-24 * time.Hour),
				ExpiresAt: now.Add(24 * time.Hour), // Not expired
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

	// Capture output and run the check command
	output := captureOutput(func() {
		runCheckCmd(mockCmd, nil)
	})

	// Check if the output indicates expired keys
	if !strings.Contains(output, "The following keys have expired") {
		t.Errorf("Expected output to indicate expired keys, got: %s", output)
	}

	// Check if only the expired key is mentioned in the output
	if !strings.Contains(output, key1) {
		t.Errorf("Expected output to mention expired key %s, got: %s", key1, output)
	}

	if strings.Contains(output, key2) {
		t.Errorf("Expected output to not mention non-expired key %s, got: %s", key2, output)
	}
}
