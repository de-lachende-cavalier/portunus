package librarian

import (
	"os"
	s "strings"
	"time"

	"testing"
)

// Tests fetching all the private key files in ~/.ssh/.
func Test_getAllKeys(t *testing.T) {
	var paths []string

	prefix, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	prefix += "/.ssh/"

	names := []string{"no", "no.pub", "authorized_keys", ".dontreadme", "yes", "yes.pub"}
	for _, name := range names {
		paths = append(paths, prefix+name)
	}

	for _, path := range paths {
		_, err := os.Create(path)
		if err != nil {
			t.Fatal(err)
		}
	}

	readPaths, err := GetAllKeys()
	if err != nil {
		t.Fatal(err)
	}

	count := 0
	for _, rp := range readPaths {
		if s.Contains(rp, "authorized_keys") || s.Contains(rp, ".dontreadme") {
			t.Fatalf("unexpected file has been read: %s", rp)
		} else if s.Contains(rp, "no") || s.Contains(rp, "yes") {
			count += 1
		}

		// clean up
		if !s.Contains(rp, "id_ed25") {
			err := os.Remove(rp)
			if err != nil {
				t.Fatal(err)
			}

			os.Remove(rp + ".pub") // remove the pub keys as well, throw away errors
		}
	}

	if count < 2 {
		t.Fatalf("expected both no and yes to be among the paths returned, got %s", readPaths)
	}
}

// Tests the fetching of expired keys.
func Test_getExpiredKeys(t *testing.T) {
	curConfig := make(map[string][2]time.Time)

	curConfig["hello"] = [2]time.Time{time.Now().Round(0),
		time.Now().Add(1 * time.Second).Round(0)}
	curConfig["friend"] = [2]time.Time{time.Now().Round(0),
		time.Now().Add(22 * time.Hour).Round(0)}
	curConfig["leave"] = [2]time.Time{time.Now().Round(0),
		time.Now().Add(2 * time.Minute).Round(0)}

	err := WriteConfig(curConfig)
	if err != nil {
		t.Fatal(err)
	}

	// to make sure only "hello" expires (so only "hello" is returned by GetExpiredKeys())
	time.Sleep(2 * time.Second)

	expired, err := GetExpiredKeys()
	if err != nil {
		t.Fatal(err)
	}

	if len(expired) != 1 {
		t.Fatalf("wrong length for returned array: expected 1, got %d", len(expired))
	}

	if expired[0] != "hello" {
		t.Fatalf("wrong file name: expected hello, got %s", expired[0])
	}

	// clean up
	err = os.Remove(configFile)
	if err != nil {
		t.Fatal(err)
	}
}

// Tests standard deletion of key files.
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
