package cmd

import (
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/mowzhja/portunus/librarian"
)

// Helper function to use instead of the default anonymous function associated with Command.Run().
func runRootCmd(cmd *cobra.Command, args []string) {
	newConfig := make(map[string][2]time.Time)

	curConfig, err := librarian.ReadConfig()
	handleErr(err)

	for keyFile, times := range curConfig {
		if _, err := os.Stat(keyFile); err == nil {
			// keyFile still exists
			newConfig[keyFile] = times
		}
	}

	// override old config
	err = librarian.WriteConfig(newConfig)
	handleErr(err)
}

var rootCmd = &cobra.Command{
	Use:   "portunus",
	Short: "portunus is a utility for managing SSH keys expiration",
	Long: `portunus acts as middleware, handling the keys for you without the need for
   ssh-keygen and keeping track of their expiration dates (as specified by you), which
   ssh-keygen cannot do. Once the keys have expired, portunus will notify you and prompt 
   to either change them (delete the old ones and make new ones) or to renew them (delay 
   their expiration by some amount you specify).`,
	Run: runRootCmd,
}

// Entry point for main.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		handleErr(err)
	}
}
