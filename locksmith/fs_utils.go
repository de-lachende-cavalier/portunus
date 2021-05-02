package locksmith

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Builds absolute paths given the filenames of the various private keys.
func buildAbsPaths(names []string) ([]string, error) {
	var b strings.Builder
	var paths []string
	home_path := os.Getenv("HOME")
	infix := "/.ssh/"

	for _, name := range names {
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

	return paths, nil
}

// Deletes the key files given as input (both public and private).
func deleteKeyFiles(pathsToDelete []string) error {
	for _, path := range pathsToDelete {
		err := os.Remove(path) // delete private key
		if err != nil {
			return err
		}

		err = os.Remove(path + ".pub") // delete public key
		if err != nil {
			return err
		}
	}

	return nil
}

// Writes the public key bytes to the corresponding file.
func writePubKey(bytes []byte, file string) error {
	err := os.WriteFile(file+".pub", bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

// Writes the private key bytes to the corresponding file.
func writePrivKey(bytes []byte, file string) error {
	err := os.WriteFile(file, bytes, 0600)
	if err != nil {
		return err
	}

	return nil
}
