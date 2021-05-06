package librarian

import (
	"fmt"
	"os"
	"testing"
)

// Helper function, creates three random files in /tmp.
func createTestFiles() []string {
	var privPaths []string
	names := []string{"gonomolo", "hyperion", "super_private"}
	base := "/tmp/"

	for _, name := range names {
		privPaths = append(privPaths, base+name)
	}

	for _, path := range privPaths {
		_, err := os.Create(path)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		_, err = os.Create(path + ".pub")
		if err != nil {
			fmt.Println(err)
			return nil
		}
	}

	// check if the files have actually been created
	for _, path := range privPaths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Printf("error creating %s", path)
			return nil
		}

		if _, err := os.Stat(path + ".pub"); os.IsNotExist(err) {
			fmt.Printf("error creating %s", path)
			return nil
		}
	}

	return privPaths
}

// Test the building of absolute paths given a slice of names.
func Test_buildAbsPaths(t *testing.T) {
	names := []string{"gonomolo", "hyperion", "super_private"}

	absPaths, err := BuildAbsPaths(names)
	if err != nil {
		t.Fatal(err)
	}

	base, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	base += "/.ssh/"

	for i, absPath := range absPaths {
		if base+names[i] != absPath {
			t.Fatalf("failed creating correct path: expected %s, got %s", base+names[i], absPath)
		}
	}
}

// Tests standard delition of key files.
func Test_deleteKeyFiles(t *testing.T) {
	testPaths := createTestFiles()
	// reminds me of C...
	if testPaths == nil {
		t.Fatal("testPaths is not supposed to be empty")
	}

	err := DeleteKeyFiles(testPaths)
	if err != nil {
		t.Fatal(err)
	}

	for _, path := range testPaths {
		if _, err := os.Stat(path); err == nil {
			t.Fatalf("failed deleting file %s", path)
		}
	}
}

// Tests deletion of nonexisting key files.
func Test_deleteKeyFiles_NonExisting(t *testing.T) {
	testPaths := []string{"/tmp/laiuwetyo93745g", "/tmp/nnnnnnnnnnnnnnnn"}

	err := DeleteKeyFiles(testPaths)
	if err == nil {
		t.Fatal("error should not be nil when deleteKeyFiles() is called with nonexisting files")
	}
}

// Tests deletion of key files where we either only have a private key file or
// only a public key file.
func Test_deleteKeyFiles_NoPairing(t *testing.T) {
	var testPaths []string

	path := "/tmp/onlyprivate"
	_, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	testPaths = append(testPaths, path)

	err = DeleteKeyFiles(testPaths)
	if err == nil {
		t.Fatal("error should not be nil when deleteKeyFiles() is called one a single private key file without corresponding pub key file")
	}
	os.Remove(path)

	testPaths = nil // reset testPaths

	path = "/tmp/onlypublic.pub"
	_, err = os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	testPaths = append(testPaths, path)

	err = DeleteKeyFiles(testPaths)
	if err == nil {
		t.Fatal("error should not be nil when deleteKeyFiles() is called one a single public key file without corresponding private key file")
	}
	os.Remove(path)
}

// No need to check writePubKey and writePrivKey seeing as they're fundamentally
// wrappers around std Go functions (which are probably thouroughly tested alreay)
