package locksmith

import (
	"fmt"
	"os"
	"testing"
)

// Helper function, creates three random files in /tmp.
func createTestFiles() []string {
	var paths []string
	// XXX for some reason the first file in this array doesn't get created
	// XXX WTF?
	names := []string{"gonomolo", "hyperion", "super_private"}
	base := "/tmp/"

	for _, name := range names {
		paths = append(paths, base+name, base+name+".pub")
	}

	for _, path := range paths {
		_, err := os.Create(path)
		if err != nil {
			fmt.Println(err)
			return nil
		}
	}

	// check if the files have actually been created
	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Printf("error creating %s", path)
			return nil
		}
	}

	return paths
}

func Test_buildAbsPaths(t *testing.T) {
	names := []string{"gonomolo", "hyperion", "super_private"}

	absPaths, err := buildAbsPaths(names)
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

func Test_deleteKeyFiles(t *testing.T) {
	testPaths := createTestFiles()
	// reminds me of C...
	if testPaths == nil {
		t.Fatal("testPaths is not supposed to be empty")
	}

	err := deleteKeyFiles(testPaths)
	if err != nil {
		t.Fatal(err)
	}

	for _, path := range testPaths {
		if _, err := os.Stat(path); err == nil {
			t.Fatalf("failed deleting file %s", path)
		}
	}
}

// No need to check writePubKey and writePrivKey seeing as they're fundamentaly just
// wrappers around std Go functions (which are probably thouroughly tested already)
