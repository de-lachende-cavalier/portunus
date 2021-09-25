package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/mowzhja/portunus/librarian"
	"github.com/mowzhja/portunus/locksmith"
)

func init() {
	rootCmd.AddCommand(rotateCmd)

	rotateCmd.Flags().StringP("cipher", "c", "ed25519",
		"Choose which cipher to use for key generation")
	rotateCmd.Flags().StringP("time", "t", "",
		"Specify for how much longer they key should be valid (format: -t <int><specifier>, where specifier is either s (seconds), m (minutes), h (hours) or d (days)")
	rotateCmd.Flags().StringP("password", "p", "",
		"Specifies the password to use with ssh-keygen (NOTE: this password is used for ALL the keys that are rotates, a compromise between security and user friendliness)")
	rotateCmd.Flags().StringSliceP("subset", "s", []string{},
		"Specify the subset of keys you want to act on")

	rotateCmd.MarkFlagRequired("time")
	rotateCmd.MarkFlagRequired("cipher")
	rotateCmd.MarkFlagRequired("password")
}

var rotateCmd = &cobra.Command{
	Use:   "rotate",
	Short: "rotate the SSH keys",
	Long: `If called without any flags, this command rotates ALL the keys in ~/.ssh/. 
  By 'rotates' I mean that the old keys are deleted and new ones are created (with the
  same name as the old ones), with new expiration dates. Once that's done, these keys
	are tracked for as long as they exist.`,
	Run: runRotateCmd,
}

// Helper function to use instead of the default anonymous function associated with Command.Run().
func runRotateCmd(cmd *cobra.Command, args []string) {
	fmt.Printf("[+] Rotating keys...\n")

	var paths []string

	configData := make(map[string][2]time.Time)
	partialData := make(map[string]time.Time)

	cipher, err := cmd.Flags().GetString("cipher")
	handleErr(err)

	delta_s, err := cmd.Flags().GetString("time")
	handleErr(err)
	delta_i, err := parseTime(delta_s)
	handleErr(err)

	passwd, err := cmd.Flags().GetString("password")
	handleErr(err)

	targetFiles, err := cmd.Flags().GetStringSlice("subset")
	handleErr(err)

	if len(targetFiles) > 0 {
		paths = buildPaths(targetFiles)
	} else {
		paths, err = librarian.GetAllKeys()
		handleErr(err)
	}

	partialData, err = locksmith.RotateKeys(paths, cipher, passwd)
	handleErr(err)

	configData = getCompleteConfig(partialData, delta_i)
	err = librarian.WriteConfig(configData)
	handleErr(err)

	fmt.Printf("[+] The keys have been succesfully rotated.\n")
}
