package librarian

import (
	"encoding/gob"
	"fmt"
	"os"
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
