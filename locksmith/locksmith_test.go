package locksmith

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Helper, builds absolute paths given the filenames of the various private keys.
func buildAbsPaths(names []string) ([]string, error) {
	var b strings.Builder
	var paths []string

	home_path := os.Getenv("HOME")
	infix := "/.ssh/"

	for _, name := range names {
		if !filepath.IsAbs(name) {
			b.WriteString(home_path)
			b.WriteString(infix)
			b.WriteString(name)

			if filepath.IsAbs(b.String()) {
				paths = append(paths, b.String())
				b.Reset()
			} else {
				err := fmt.Errorf("failed building absolute path for %s (partial result: %s)", name, b.String())
				return nil, err
			}
		}
	}

	return paths, nil
}

// Helper function to create some key files.
func createTestKeyFiles() []string {
	names := []string{"exp1", "exp2", "exp3"}

	paths, err := buildAbsPaths(names) // test using authentic path
	if err != nil {
		fmt.Println(err)
		return nil
	}

	for _, path := range paths {
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

	// check that files have been properly created
	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Println(err)
			return nil
		}

		if _, err := os.Stat(path + ".pub"); os.IsNotExist(err) {
			fmt.Println(err)
			return nil
		}
	}

	return paths
}

// Helper function, used to clean up the files created for testing.
func cleanupTestKeyFiles() {
	names := []string{"exp1", "exp2", "exp3"}

	paths, err := buildAbsPaths(names)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, path := range paths {
		err := os.Remove(path)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = os.Remove(path + ".pub")
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

// Tests the key changing functionality using RSA.
func Test_RotateKeys_RSA(t *testing.T) {
	oldData := make(map[string]time.Time)

	paths := createTestKeyFiles()
	if paths == nil {
		t.Fatal("paths is not supposed to be empty")
	}
	defer cleanupTestKeyFiles()

	for _, path := range paths {
		oldData[path] = time.Now()
	}

	newData, err := RotateKeys(paths, "rsa")
	if err != nil {
		t.Fatal(err)
	}

	for path := range newData {
		if newData[path].Sub(oldData[path]) < 0 {
			// oldCreationTime > newCreationTime
			t.Fatalf("the new creation time for %s cannot precede its old creation time: expected %q > %q", path, newData[path], oldData[path])
		}
	}
}

// Tests the key changing functionality using Ed25519.
func Test_RotateKeys_Ed25519(t *testing.T) {
	oldData := make(map[string]time.Time)

	paths := createTestKeyFiles()
	if paths == nil {
		t.Fatal("paths is not supposed to be empty")
	}
	defer cleanupTestKeyFiles()

	for _, path := range paths {
		oldData[path] = time.Now()
	}

	newData, err := RotateKeys(paths, "ed25519")
	if err != nil {
		t.Fatal(err)
	}

	for path := range newData {
		if newData[path].Sub(oldData[path]) < 0 {
			// oldCreationTime > newCreationTime
			t.Fatalf("the new creation time for %s cannot precede its old creation time: expected %q > %q", path, newData[path], oldData[path])
		}
	}
}
