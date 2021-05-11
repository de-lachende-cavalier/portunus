package cmd

import (
	"fmt"
	sc "strconv"
	s "strings"
	"time"
	"path/filepath"
	"os"

	"github.com/spf13/cobra"

	"github.com/mowzhja/portunus/librarian"
	"github.com/mowzhja/portunus/locksmith"
)

func init() {
	rootCmd.AddCommand(rotateCmd)

	rotateCmd.Flags().StringP("cipher", "c", "ed25519", "Choose which cipher to use for key generation (default is Ed25519)")
	rotateCmd.Flags().StringP("time", "t", "",
		"Specify for how much longer they key should be valid (format: -t <int><specifier>, where specifier is either s (seconds), m (minutes), h (hours) or d (days)")
	rotateCmd.Flags().StringSliceP("subset", "s", []string{}, "Specify the subset of keys you want to act on")

	rotateCmd.MarkFlagRequired("time")
}

var rotateCmd = &cobra.Command{
	Use:   "rotate",
	Short: "rotate the SSH keys",
	Long: `If called without any flags, this command rotates ALL the keys in ~/.ssh/. 
  By 'rotates' I mean that the old keys are deleted and new ones are created (with the
  same name as the old ones), with new expiration dates. Once that's done, these keys
	are tracked for as long as they exist.`,

	Run: func(cmd *cobra.Command, args []string) {
		configData := make(map[string][2]time.Time)
		partialData := make(map[string]time.Time)

		cipher, err := cmd.Flags().GetString("cipher")
		if err != nil {
			fmt.Println(err)
			return
		}

		delta_s, err := cmd.Flags().GetString("time")
		if err != nil {
			fmt.Println(err)
			return
		}

		delta_i, err := parseTime(delta_s)
		if err != nil {
			fmt.Println(err)
			return
		}

		targetFiles, err := cmd.Flags().GetStringSlice("subset")
		if err != nil {
			fmt.Println(err)
			return
		}

		if len(targetFiles) > 0 {
			targetPaths := buildPaths(targetFiles)
			partialData, err = locksmith.RotateKeys(targetPaths, cipher)
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			paths, err := librarian.GetAllKeys()
			if err != nil {
				fmt.Println(err)
				return
			}

			partialData, err = locksmith.RotateKeys(paths, cipher)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		configData = getCompleteConfig(partialData, delta_i)

		err = librarian.WriteConfig(configData)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

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
			filePaths = append(filePaths, prefix + name)
		}
	}

	return filePaths
}
