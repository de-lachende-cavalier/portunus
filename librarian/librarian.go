// Package librarian contains all the useful functionality for the interaction with the fs.
package librarian

import (
	"os"
	"strings"
	"time"
)

// Gets all keys stored in ~/.ssh/ and returns the paths corresponding to the private ones.
func GetAllKeys() ([]string, error) {
	var paths []string
	dir := os.Getenv("HOME") + "/.ssh/"

	isValid := func(s string) bool {
		if strings.Contains(s, "authorized_keys") ||
			strings.Contains(s, "known_hosts") ||
			strings.Contains(s, "config") ||
			strings.HasSuffix(s, ".pub") ||
			strings.HasPrefix(s, ".") {
			return false
		}

		return true
	}

	d, err := os.Open(dir)
	if err != nil {
		return nil, err
	}

	keyFiles, err := d.Readdir(-1) // read all the files
	if err != nil {
		return nil, err
	}
	d.Close()

	for _, file := range keyFiles {
		if isValid(file.Name()) {
			paths = append(paths, dir+file.Name())
		}
	}

	return paths, nil
}

// Check key expiration, if any have expired return the corresponding paths.
func GetExpiredKeys() ([]string, error) {
	var expired []string

	configMap, err := ReadConfig()
	if err != nil {
		return nil, err
	}

	current := time.Now()
	for keyFile, dates := range configMap {
		if current.Sub(dates[1]) >= 0 {
			// key has expired
			expired = append(expired, keyFile)
		}
	}

	return expired, nil
}

// Deletes the key files given as input (both public and private).
func DeleteKeyFiles(pathsToDelete []string) error {
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
