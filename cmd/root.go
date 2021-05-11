package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "portunus",
	Short: "portunus is a utility for managing SSH keys expiration",
	Long: `portunus acts as middleware, handling the keys for you without the need for
   ssh-keygen and keeping track of their expiration dates (as specified by you), which
   ssh-keygen cannot do. Once the keys have expired, portunus will notify you and prompt 
   to either change them (delete the old ones and make new ones) or to renew them (delay 
   their expiration by some amount you specify).`,

	// TODO make a run here that checks whether some keys were deleted from .ssh, and if
	// so remove those keys from the ones tracked
	Run: func(cmd *cobra.Command, args []string) {
	},
}
