// Package tracker contains all the functionality that allows tracking files in various directories.
package librarian

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Reads the config data relative to key expiration and properly decodes it.
func readConfig(configFile string) (map[string][2]time.Time, error) {
	Map := make(map[string][2]time.Time)

	f, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}

	d := gob.NewDecoder(f)
	err = d.Decode(&Map)
	if err != nil {
		return nil, err
	}

	return Map, nil
}

// Check key expiration, if the keys have expired return true.
func GetExpired() ([]string, error) {
	var expired []string

	configMap, err := readConfig("key_info.gob")
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

// Builds absolute paths given the filenames of the various private keys.
// TODO => what if we pass it absolute paths already? => this function should just
// return them as is
func BuildAbsPaths(names []string) ([]string, error) {
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

// Writes the public key bytes to the corresponding file.
func WritePubKey(bytes []byte, file string) error {
	err := os.WriteFile(file, bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

// Writes the private key bytes to the corresponding file.
func WritePrivKey(bytes []byte, file string) error {
	err := os.WriteFile(file, bytes, 0600)
	if err != nil {
		return err
	}

	return nil
}

// Writes config data to the proper file with the proper encoding (gob).
func WriteConfig(data map[string][2]time.Time, configFile string) error {
	file, err := os.Open(configFile)
	if err != nil {
		return err
	}
	defer file.Close()

	e := gob.NewEncoder(file)
	
	err = e.Encode(data)
	if err != nil {
		return err
	}

	return nil
}
