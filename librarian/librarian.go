// Package librarian contains all the useful functionality for the interaction with the fs.
package librarian

import (
	"encoding/gob"
	"os"
	"time"
)

var configFile = os.Getenv("HOME") + ".portunus_data.gob"

// Reads the config data relative to key expiration and properly decodes it.
func ReadConfig() (map[string][2]time.Time, error) {
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

// Gets all keys stored in ~/.ssh/.
func GetAllKeys() {
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

// Writes config data to the proper file with the proper encoding (gob).
func WriteConfig(data map[string][2]time.Time) error {
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
