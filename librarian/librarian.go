// Package librarian contains all the useful functionality for the interaction with the fs.
package librarian

import (
	"encoding/gob"
	"fmt"
	"os"
	"strings"
	"time"
)

var configFile = os.Getenv("HOME") + "/.portunus_data.gob"

// Tests the data to write in the config file for validity.
func checkConfigData(data map[string][2]time.Time) error {
	for file, times := range data {
		if times[0].Sub(times[1]) >= 0 {
			err := fmt.Errorf("creation time >= expiration time for file %s!", file)
			return err
		}
	}

	return nil
}

// Reads the config data relative to key expiration and properly decodes it.
func ReadConfig() (map[string][2]time.Time, error) {
	Map := make(map[string][2]time.Time)

	f, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	d := gob.NewDecoder(f)

	err = d.Decode(&Map)
	if err != nil {
		return nil, err
	}

	return Map, nil
}

// Writes config data to the proper file with the proper encoding (gob).
//
// If the file doesn't exist already, it gets created. Data is checked for validity
// before being written.
func WriteConfig(data map[string][2]time.Time) error {
	if err := checkConfigData(data); err != nil {
		return err
	}

	// create file if it doesn't already exist
	// truncate it before writing
	file, err := os.OpenFile(configFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
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
