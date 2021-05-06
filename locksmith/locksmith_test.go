package locksmith

import (
	"fmt"
	"os"
	"testing"
	"time"
)

// Helper function to create some key files.
func createTestKeyFiles() []string {
	var paths []string
	names := []string{"exp1", "exp2", "exp3"}
	base := "/tmp/"

	for _, name := range names {
		_, err := os.Create(base + name)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		_, err = os.Create(base + name + ".pub")
		if err != nil {
			fmt.Println(err)
			return nil
		}

		paths = append(paths, base+name, base+name+".pub")
	}

	return paths
}

// Helper function, used to clean up the files created for testing.
func cleanupTestKeyFiles() {
	names := []string{"exp1", "exp2", "exp3"}
	base := "/tmp/"

	for _, name := range names {
		err := os.Remove(base + name)
		if err != nil {
			fmt.Println(err)
		}
		err = os.Remove(base + name + ".pub")
		if err != nil {
			fmt.Println(err)
		}
	}
}

// Tests the key changing functionality using RSA.
func Test_ChangeKeys_RSA(t *testing.T) {
	oldData := make(map[string]time.Time)

	paths := createTestKeyFiles()
	if paths == nil {
		t.Fatal("paths is not supposed to be empty")
	}
	defer cleanupTestKeyFiles()

	for _, path := range paths {
		oldData[path] = time.Now()
	}

	newData, err := ChangeKeys(paths, "rsa")
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
func Test_ChangeKeys_Ed25519(t *testing.T) {
	oldData := make(map[string]time.Time)

	paths := createTestKeyFiles()
	if paths == nil {
		t.Fatal("paths is not supposed to be empty")
	}
	defer cleanupTestKeyFiles()

	for _, path := range paths {
		oldData[path] = time.Now()
	}

	newData, err := ChangeKeys(paths, "ed25519")
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
