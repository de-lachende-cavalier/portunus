package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

// expandPath expands a path with ~ to the user's home directory
func expandPath(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("empty path")
	}

	if path[0] != '~' {
		return path, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	return filepath.Join(homeDir, path[1:]), nil
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
