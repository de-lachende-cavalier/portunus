package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	sc "strconv"
	s "strings"
	"time"
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
		err := fmt.Errorf("Wrong format, %s not recognized.", specifier)

		return 0, err
	}
}

// Creates the complete config data given the partial data received from the locksmith.
func getCompleteConfig(partialConfig map[string]time.Time, expirationDelta int) map[string][2]time.Time {
	completeConfig := make(map[string][2]time.Time)

	for keyFile, creationTime := range partialConfig {
		times := [2]time.Time{}

		times[0] = creationTime
		times[1] = creationTime.Add(time.Second * time.Duration(expirationDelta)) // expiration time

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
