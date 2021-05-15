package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mowzhja/portunus/librarian"
)

func init() {
	rootCmd.AddCommand(checkCmd)
}

// Helper function to use instead of the default anonymous function associated with Command.Run().
func runCheckCmd(cmd *cobra.Command, args []string) {
	expiredPaths, err := librarian.GetExpiredKeys()
	handleErr(err)

	if len(expiredPaths) > 0 {
		// we have expired keys
		fmt.Println("The following keys have expired:", expiredPaths)
		fmt.Println("Either renew or rotate them!")
		return
	}
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "checks ~/.ssh/ for expired keys",
	Long: `This is the command that portunus should be run with then put
  in your bashrc/zshrc, etc. It checks whether you have expired keys by
  examining portunus' config file and notifies you of them`,
	Run: runCheckCmd,
}
