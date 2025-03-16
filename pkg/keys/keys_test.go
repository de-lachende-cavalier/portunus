package keys

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/de-lachende-cavalier/portunus/pkg/testutil"
)

// TestNewManager tests the NewManager function
func TestNewManager(t *testing.T) {
	// Create a test SSH directory
	sshDir := testutil.CreateTestSSHDir(t)

	// Set the HOME environment variable to the test directory
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	os.Setenv("HOME", filepath.Dir(sshDir))

	// Create a new manager
	manager, err := NewManager()
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	// Check if the manager has the correct SSH directory
	if manager.sshDir != sshDir {
		t.Errorf("Expected SSH directory %s, got %s", sshDir, manager.sshDir)
	}
}

// TestManager_GetAllKeys tests the GetAllKeys method
func TestManager_GetAllKeys(t *testing.T) {
	// Create a test SSH directory
	sshDir := testutil.CreateTestSSHDir(t)

	// Create test key files
	key1, _ := testutil.CreateTestKeyPair(t, sshDir, "id_ed25519")
	key2, _ := testutil.CreateTestKeyPair(t, sshDir, "id_rsa")

	// Create some non-key files
	testutil.CreateTestFile(t, sshDir, "authorized_keys", "test")
	testutil.CreateTestFile(t, sshDir, "known_hosts", "test")
	testutil.CreateTestFile(t, sshDir, "config", "test")

	// Create a manager with the test SSH directory
	manager := &Manager{
		sshDir: sshDir,
	}

	// Get all keys
	keys, err := manager.GetAllKeys(context.Background())
	if err != nil {
		t.Fatalf("Failed to get all keys: %v", err)
	}

	// Check if the correct keys are returned
	if len(keys) != 2 {
		t.Errorf("Expected 2 keys, got %d", len(keys))
	}

	// Check if the keys are in the result
	foundKey1 := false
	foundKey2 := false

	for _, key := range keys {
		if key == key1 {
			foundKey1 = true
		} else if key == key2 {
			foundKey2 = true
		}
	}

	if !foundKey1 {
		t.Errorf("Expected key %s to be in the result", key1)
	}

	if !foundKey2 {
		t.Errorf("Expected key %s to be in the result", key2)
	}
}

// TestManager_DeleteKeyPair tests the DeleteKeyPair method
func TestManager_DeleteKeyPair(t *testing.T) {
	// Create a test SSH directory
	sshDir := testutil.CreateTestSSHDir(t)

	// Create a test key pair
	keyPath, pubKeyPath := testutil.CreateTestKeyPair(t, sshDir, "id_ed25519")

	// Create a manager with the test SSH directory
	manager := &Manager{
		sshDir: sshDir,
	}

	// Delete the key pair
	err := manager.DeleteKeyPair(context.Background(), keyPath)
	if err != nil {
		t.Fatalf("Failed to delete key pair: %v", err)
	}

	// Check if the key files are deleted
	testutil.AssertFileNotExists(t, keyPath)
	testutil.AssertFileNotExists(t, pubKeyPath)
}

// TestManager_GenerateKeyPair tests the GenerateKeyPair method
func TestManager_GenerateKeyPair(t *testing.T) {
	// Skip this test if ssh-keygen is not available
	if _, err := os.Stat("/usr/bin/ssh-keygen"); os.IsNotExist(err) {
		t.Skip("ssh-keygen not available, skipping test")
	}

	// Create a test SSH directory
	sshDir := testutil.CreateTestSSHDir(t)

	// Create a manager with the test SSH directory
	manager := &Manager{
		sshDir: sshDir,
	}

	// Generate a key pair
	keyPath := filepath.Join(sshDir, "test_key")
	err := manager.GenerateKeyPair(context.Background(), keyPath, "ed25519", "")
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	// Check if the key files are created
	testutil.AssertFileExists(t, keyPath)
	testutil.AssertFileExists(t, keyPath+".pub")
}

// TestManager_RotateKeys tests the RotateKeys method
func TestManager_RotateKeys(t *testing.T) {
	// Skip this test if ssh-keygen is not available
	if _, err := os.Stat("/usr/bin/ssh-keygen"); os.IsNotExist(err) {
		t.Skip("ssh-keygen not available, skipping test")
	}

	// Create a test SSH directory
	sshDir := testutil.CreateTestSSHDir(t)

	// Create test key files
	key1, _ := testutil.CreateTestKeyPair(t, sshDir, "id_ed25519")
	key2, _ := testutil.CreateTestKeyPair(t, sshDir, "id_rsa")

	// Create a manager with the test SSH directory
	manager := &Manager{
		sshDir: sshDir,
	}

	// Rotate the keys
	creationTimes, err := manager.RotateKeys(context.Background(), []string{key1, key2}, "ed25519", "")
	if err != nil {
		t.Fatalf("Failed to rotate keys: %v", err)
	}

	// Check if the creation times are returned
	if len(creationTimes) != 2 {
		t.Errorf("Expected 2 creation times, got %d", len(creationTimes))
	}

	// Check if the key files still exist (they should be recreated)
	testutil.AssertFileExists(t, key1)
	testutil.AssertFileExists(t, key1+".pub")
	testutil.AssertFileExists(t, key2)
	testutil.AssertFileExists(t, key2+".pub")
}

// TestIsPrivateKey tests the isPrivateKey function
func TestIsPrivateKey(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     bool
	}{
		{"ValidKey", "id_ed25519", true},
		{"ValidKey2", "id_rsa", true},
		{"PublicKey", "id_ed25519.pub", false},
		{"AuthorizedKeys", "authorized_keys", false},
		{"KnownHosts", "known_hosts", false},
		{"Config", "config", false},
		{"HiddenFile", ".hidden", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isPrivateKey(tt.filename); got != tt.want {
				t.Errorf("isPrivateKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
