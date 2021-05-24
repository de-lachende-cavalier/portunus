package cmd

import (
	"os"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/mowzhja/portunus/librarian"
)

// Check the integrity of the config file each time portunus is called.
func init() {
	newConfig := make(map[string][2]time.Time)

	curConfig, err := librarian.ReadConfig()
	if err != nil {
		fmt.Printf("[+] Config file either missing or empty.\n")
		return
	}

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
	Run: func(cmd *cobra.Command, args []string) {
	},
}

// Entry point for main.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		handleErr(err)
	}
}
