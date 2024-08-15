package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/de-lachende-cavalier/portunus/librarian"
)

func init() {
	rootCmd.AddCommand(checkCmd)
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "checks ~/.ssh for expired keys",
	Long:  `This is the command that portunus should be run with when put in your bashrc/zshrc. It checks whether you have expired keys by examining portunus' config file and notifies you.`,
	Run:   runCheckCmd,
}

// Helper function to use instead of the default anonymous function associated with Command.Run().
func runCheckCmd(cmd *cobra.Command, args []string) {
	fmt.Printf("[+] Checking for expired keys...\n")

	expiredPaths, err := librarian.GetExpiredKeys()
	handleErr(err)

	if len(expiredPaths) > 0 {
		// we have expired keys
		fmt.Printf("[+] The following keys have expired: \n\t%s\n", expiredPaths)
		fmt.Printf("[+] Either renew or rotate them.\n")
		return
	}

	fmt.Printf("[+] The keys are still fresh!\n")
}
