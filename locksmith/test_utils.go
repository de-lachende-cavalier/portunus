package locksmith

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Cleans up the file created for testing purposes.
func cleanupPaths(paths []string) error {
	for _, path := range paths {
		err := os.Remove(path)
		if err != nil {
			return err
		}

		err = os.Remove(path + ".pub")
		if err != nil {
			return err
		}
	}

	return nil
}

// Checks that the files have been correctly created.
func checkPathsExist(paths []string) error {
	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return err
		}
		if _, err := os.Stat(path + ".pub"); os.IsNotExist(err) {
			return err
		}
	}

	return nil
}

// Builds absolute paths given the filenames of the various private keys (for locksmith_test)
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

// Helper function to create some key files (for locksmith_test).
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
