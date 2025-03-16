// Package keys provides SSH key management functionality
package keys

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/de-lachende-cavalier/portunus/pkg/logger"
)

// SupportedCiphers is a list of supported ciphers
var SupportedCiphers = map[string]bool{
	"ed25519": true,
	"rsa":     true,
	"ecdsa":   true,
}

// Manager handles SSH key operations
type Manager struct {
	sshDir string
}

// NewManager creates a new key manager
func NewManager() (*Manager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	sshDir := filepath.Join(homeDir, ".ssh")
	if err := ensureDir(sshDir); err != nil {
		return nil, fmt.Errorf("failed to ensure SSH directory exists: %w", err)
	}

	return &Manager{
		sshDir: sshDir,
	}, nil
}

// ensureDir ensures that the directory exists
func ensureDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0700)
	}
	return nil
}

// GetAllKeys returns all private keys in the SSH directory
func (m *Manager) GetAllKeys(ctx context.Context) ([]string, error) {
	var keys []string

	entries, err := os.ReadDir(m.sshDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read SSH directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if isPrivateKey(name) {
			keys = append(keys, filepath.Join(m.sshDir, name))
		}
	}

	return keys, nil
}

// isPrivateKey checks if a file is likely a private key
func isPrivateKey(name string) bool {
	// Skip common non-key files
	if strings.Contains(name, "authorized_keys") ||
		strings.Contains(name, "known_hosts") ||
		strings.Contains(name, "config") ||
		strings.HasSuffix(name, ".pub") ||
		strings.HasPrefix(name, ".") {
		return false
	}
	return true
}

// DeleteKeyPair deletes both the private and public key files
func (m *Manager) DeleteKeyPair(ctx context.Context, path string) error {
	// Delete private key
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete private key %s: %w", path, err)
	}

	// Delete public key
	pubPath := path + ".pub"
	if err := os.Remove(pubPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete public key %s: %w", pubPath, err)
	}

	return nil
}

// GenerateKeyPair generates a new SSH key pair
func (m *Manager) GenerateKeyPair(ctx context.Context, path, cipher, password string) error {
	if !SupportedCiphers[cipher] {
		return fmt.Errorf("unsupported cipher: %s", cipher)
	}

	// Remove existing keys to avoid ssh-keygen prompts
	_ = os.Remove(path)
	_ = os.Remove(path + ".pub")

	args := []string{"-q", "-t", cipher, "-N", password, "-f", path, "-a", "100"}

	// Add specific options based on cipher type
	switch cipher {
	case "rsa":
		args = append(args, "-b", "4096")
	case "ecdsa":
		args = append(args, "-b", "521")
	}

	cmd := exec.CommandContext(ctx, "ssh-keygen", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ssh-keygen failed: %w, output: %s", err, string(output))
	}

	// Verify the keys were created
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("private key was not created at %s", path)
	}
	if _, err := os.Stat(path + ".pub"); os.IsNotExist(err) {
		return fmt.Errorf("public key was not created at %s.pub", path)
	}

	logger.Infof("Generated new key pair: %s", path)
	return nil
}

// RotateKeys rotates the specified keys
func (m *Manager) RotateKeys(ctx context.Context, paths []string, cipher, password string) (map[string]time.Time, error) {
	creationTimes := make(map[string]time.Time)

	for _, path := range paths {
		if err := m.DeleteKeyPair(ctx, path); err != nil {
			return nil, err
		}

		if err := m.GenerateKeyPair(ctx, path, cipher, password); err != nil {
			return nil, err
		}

		creationTimes[path] = time.Now()
	}

	return creationTimes, nil
}
