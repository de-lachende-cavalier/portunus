package locksmith

import (
	"testing"
	"time"
)

// Tests the key changing functionality using RSA.
func Test_RotateKeys_RSA(t *testing.T) {
	passwd := "rsa_rotate_keys"
	paths := createTestKeyFiles()
	if paths == nil {
		t.Fatal("the paths are not supposed to be nil")
	}

	beforeCreation := time.Now()
	upData, err := RotateKeys(paths, "rsa", passwd)

	// check that the files have been correctly recreated after deletion
	err = checkPathsExist(paths)
	if err != nil {
		t.Fatal(err)
	}

	for path, creationTime := range upData {
		// check equality with (10) second accuracy
		// (to account for buffering, scheduling et alia)
		if creationTime.Round(10*time.Second) != beforeCreation.Round(10*time.Second) {
			t.Fatalf("the creation time for %s has been set incorrectly: expected %q, got %q",
				path, beforeCreation.Round(10*time.Second), creationTime.Round(10*time.Second))
		}
	}

	err = cleanupPaths(paths)
	if err != nil {
		t.Fatal(err)
	}
}

// Tests the key changing functionality using Ed25519.
func Test_RotateKeys_Ed25519(t *testing.T) {
	passwd := "ed25519_rotate_keys"
	paths := createTestKeyFiles()
	if paths == nil {
		t.Fatal("the paths are not supposed to be nil")
	}

	beforeCreation := time.Now()
	upData, err := RotateKeys(paths, "ed25519", passwd)

	// check that the files have been correctly recreated after deletion
	err = checkPathsExist(paths)
	if err != nil {
		t.Fatal(err)
	}

	for path, creationTime := range upData {
		// check equality with (10) second accuracy
		// (to account for buffering, scheduling et alia)
		if creationTime.Round(10*time.Second) != beforeCreation.Round(10*time.Second) {
			t.Fatalf("the creation time for %s has been set incorrectly: expected %q, got %q",
				path, beforeCreation.Round(10*time.Second), creationTime.Round(10*time.Second))
		}
	}

	err = cleanupPaths(paths)
	if err != nil {
		t.Fatal(err)
	}
}
