package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	sc "strconv"
	s "strings"
	"time"

	"github.com/mowzhja/portunus/librarian"
)

// Parses the time spec given by the user and return the number of seconds corresponding to it.
func parseTime(usr_input string) (int, error) {
	specifier := usr_input[len(usr_input)-1:]
	value := s.Trim(usr_input, "smhd")

	switch specifier {
	case "s":
		return sc.Atoi(value)
	case "m":
		n, err := sc.Atoi(value)
		if err != nil {
			return 0, err
		}

		return n * 60, nil
	case "h":
		n, err := sc.Atoi(value)
		if err != nil {
			return 0, err
		}

		return n * 3600, nil
	case "d":
		n, err := sc.Atoi(value)
		if err != nil {
			return 0, err
		}

		return n * 86400, nil
	default:
		return 0, errors.New("error in specifer, either not recognized or not specified")
	}
}

// Creates the complete config data given the partial data received from the locksmith.
func getCompleteConfig(partialConfig map[string]time.Time, expirationDelta int) map[string][2]time.Time {
	completeConfig := make(map[string][2]time.Time)

	for keyFile, creationTime := range partialConfig {
		times := [2]time.Time{}
		times[0] = creationTime.Round(0)
		times[1] = creationTime.Add(time.Second * time.Duration(expirationDelta)).Round(0) // expiration time

		completeConfig[keyFile] = times
	}

	return completeConfig
}

// Builds the correct paths given the filenames specified with the --subset flag.
func buildPaths(fileNames []string) []string {
	var filePaths []string

	for _, name := range fileNames {
		if !filepath.IsAbs(name) {
			prefix := os.Getenv("HOME") + "/.ssh/"
			filePaths = append(filePaths, prefix+name)
		} else {
			filePaths = append(filePaths, name)
		}
	}

	return filePaths
}

// Helper for handling errors in a single line.
func handleErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Helper for tests that need to read some data from the portunus_data.gob file.
func writeTestConfig() (map[string][2]time.Time, error) {
	curConfig := make(map[string][2]time.Time)

	_, err := os.Create(os.Getenv("HOME") + "/.portunus_data.gob")
	if err != nil {
		return nil, err
	}

	curConfig["hello"] = [2]time.Time{time.Now().Round(0),
		time.Now().Add(1 * time.Second).Round(0)}
	curConfig["friend"] = [2]time.Time{time.Now().Round(0),
		time.Now().Add(5 * time.Second).Round(0)}
	curConfig["leave"] = [2]time.Time{time.Now().Round(0),
		time.Now().Add(4 * time.Second).Round(0)}

	err = librarian.WriteConfig(curConfig)
	if err != nil {
		return nil, err
	}

	return curConfig, nil
}

// Helper to cleanup the test data from the portunus_data.gob file.
func cleanupTestConfig() error {
	err := os.Remove(os.Getenv("HOME") + "/.portunus_data.gob")
	if err != nil {
		return err
	}

	return nil
}
