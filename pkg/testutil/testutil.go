// Package testutil provides utilities for testing
package testutil

import (
	"os"
	"path/filepath"
	"testing"
)

// TempDir creates a temporary directory for testing
func TempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "portunus-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	t.Cleanup(func() {
		os.RemoveAll(dir)
	})
	return dir
}

// CreateTestSSHDir creates a temporary SSH directory for testing
func CreateTestSSHDir(t *testing.T) string {
	t.Helper()
	baseDir := TempDir(t)
	sshDir := filepath.Join(baseDir, ".ssh")

	if err := os.MkdirAll(sshDir, 0700); err != nil {
		t.Fatalf("Failed to create test SSH dir: %v", err)
	}

	return sshDir
}

// CreateTestFile creates a test file with the given content
func CreateTestFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	return path
}

// CreateTestKeyFile creates a test key file
func CreateTestKeyFile(t *testing.T, sshDir, name string) string {
	t.Helper()
	content := `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACDFIzQzGAduJeaGNEX1Sf+T8EdVQEwmQUc7C/iHjdtDaQAAAJgwIvCYMCLw
mAAAAAtzc2gtZWQyNTUxOQAAACDFIzQzGAduJeaGNEX1Sf+T8EdVQEwmQUc7C/iHjdtDaQ
AAAEAIUUzFMUJ+vNJ57OwzxMW8+4ZJQIxGF2GYy3VTMFQTxsUjNDMYB24l5oY0RfVJ/5Pw
R1VATCZBRzsL+IeN20NpAAAAEXRlc3RAZXhhbXBsZS5jb20BAg==
-----END OPENSSH PRIVATE KEY-----`

	return CreateTestFile(t, sshDir, name, content)
}

// CreateTestPublicKeyFile creates a test public key file
func CreateTestPublicKeyFile(t *testing.T, sshDir, name string) string {
	t.Helper()
	content := "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIMUjNDMYB24l5oY0RfVJ/5PwR1VATCZBRzsL+IeN20Np test@example.com"

	return CreateTestFile(t, sshDir, name+".pub", content)
}

// CreateTestKeyPair creates a test key pair (private and public)
func CreateTestKeyPair(t *testing.T, sshDir, name string) (string, string) {
	t.Helper()
	privKey := CreateTestKeyFile(t, sshDir, name)
	pubKey := CreateTestPublicKeyFile(t, sshDir, name)
	return privKey, pubKey
}

// AssertFileExists checks if a file exists
func AssertFileExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("Expected file to exist: %s", path)
	}
}

// AssertFileNotExists checks if a file does not exist
func AssertFileNotExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("Expected file to not exist: %s", path)
	}
}
