package cmd

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/de-lachende-cavalier/portunus/pkg/testutil"
)

// Test parsing of the input times using days.
func Test_parseDuration_Days(t *testing.T) {
	duration, err := parseDuration("1d")
	if err != nil {
		t.Fatalf("Failed to parse duration: %v", err)
	}
	if duration != 24*time.Hour {
		t.Fatalf("Expected 24h, got %v", duration)
	}

	duration, err = parseDuration("365d")
	if err != nil {
		t.Fatalf("Failed to parse duration: %v", err)
	}
	if duration != 365*24*time.Hour {
		t.Fatalf("Expected 8760h, got %v", duration)
	}
}

// Test parsing of the input times using hours.
func Test_parseDuration_Hours(t *testing.T) {
	duration, err := parseDuration("24h")
	if err != nil {
		t.Fatalf("Failed to parse duration: %v", err)
	}
	if duration != 24*time.Hour {
		t.Fatalf("Expected 24h, got %v", duration)
	}

	duration, err = parseDuration("1h")
	if err != nil {
		t.Fatalf("Failed to parse duration: %v", err)
	}
	if duration != time.Hour {
		t.Fatalf("Expected 1h, got %v", duration)
	}
}

// Test parsing of the input times using minutes.
func Test_parseDuration_Minutes(t *testing.T) {
	duration, err := parseDuration("60m")
	if err != nil {
		t.Fatalf("Failed to parse duration: %v", err)
	}
	if duration != time.Hour {
		t.Fatalf("Expected 1h, got %v", duration)
	}

	duration, err = parseDuration("1m")
	if err != nil {
		t.Fatalf("Failed to parse duration: %v", err)
	}
	if duration != time.Minute {
		t.Fatalf("Expected 1m, got %v", duration)
	}
}

// Test parsing of the input times using seconds.
func Test_parseDuration_Seconds(t *testing.T) {
	duration, err := parseDuration("3600s")
	if err != nil {
		t.Fatalf("Failed to parse duration: %v", err)
	}
	if duration != time.Hour {
		t.Fatalf("Expected 1h, got %v", duration)
	}
}

// Test parsing of the input containing an invalid specifier/no specifer.
func Test_parseDuration_InvalidSpecifier(t *testing.T) {
	_, err := parseDuration("80l")
	if err == nil {
		t.Fatal("Expected error for invalid specifier, got nil")
	}

	_, err = parseDuration("80")
	if err == nil {
		t.Fatal("Expected error for missing unit, got nil")
	}
}

// Test expandPath function
func Test_expandPath(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := testutil.TempDir(t)

	// Set the HOME environment variable to the test directory
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tempDir)

	// Test with path starting with ~
	path, err := expandPath("~/test/path")
	if err != nil {
		t.Fatalf("Failed to expand path: %v", err)
	}
	expected := filepath.Join(tempDir, "test/path")
	if path != expected {
		t.Errorf("Expected path %s, got %s", expected, path)
	}

	// Test with absolute path
	absPath := "/absolute/path"
	path, err = expandPath(absPath)
	if err != nil {
		t.Fatalf("Failed to expand path: %v", err)
	}
	if path != absPath {
		t.Errorf("Expected path %s, got %s", absPath, path)
	}

	// Test with empty path
	_, err = expandPath("")
	if err == nil {
		t.Error("Expected error for empty path, got nil")
	}
}

// Test fileExists function
func Test_fileExists(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := testutil.TempDir(t)

	// Create a test file
	testFile := filepath.Join(tempDir, "test_file")
	if err := os.WriteFile(testFile, []byte("test"), 0600); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test with existing file
	if !fileExists(testFile) {
		t.Errorf("Expected file %s to exist", testFile)
	}

	// Test with non-existent file
	nonExistentFile := filepath.Join(tempDir, "non_existent_file")
	if fileExists(nonExistentFile) {
		t.Errorf("Expected file %s to not exist", nonExistentFile)
	}

	// Test with directory
	if fileExists(tempDir) {
		t.Errorf("Expected directory %s to not be reported as a file", tempDir)
	}
}
