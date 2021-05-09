package cmd

import (
	"fmt"
	sc "strconv"
	s "strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/mowzhja/portunus/librarian"
	"github.com/mowzhja/portunus/locksmith"
)

func init() {
	rootCmd.AddCommand(rotateCmd)
}

var rotateCmd = &cobra.Command{
	Use:   "rotate",
	Short: "rotate the SSH keys",
	Long: `If called without any flags, this command rotates ALL the keys in ~/.ssh/. 
  By 'rotates' I mean that the old keys are deleted and new ones are created (with the
  same name as the old ones), with new expiration dates.
  The new expiration date is specified in deltas, aka you specify how much the key 
  should live for from the time of creation.
  Multiple time formats are used for convenience: s (seconds), m (minutes), h (hours), 
  d (days). The value is expected in the form <int><format>, where <format> takes 
  values from the options above`,
	// TODO add a flag to allow users to set expiry on a subset of keys in ~/.ssh

	Run: func(cmd *cobra.Command, args []string) {
		configData := make(map[string][2]time.Time)

		// more or less something like this:
		paths := librarian.GetAllKeys()
		name_creation := locksmith.RotateKeys(paths, "ed25519")

		configData = makeProperMapToStoreInConfig(name_creation)
		err := librarian.WriteConfig(configData)
		if err != nil {
			// handle it
		}
	},
}

// Parses the time spec given by the user and return the number of seconds corresponding to it.
func parseTime(usr_input string) (int, error) {
	qualifier := usr_input[len(usr_input)-1:]
	value := s.Trim(usr_input, "smhd")

	switch qualifier {
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
		err := fmt.Errorf("Wrong format, %s not recognized.", qualifier)

		return 0, err
	}
}
