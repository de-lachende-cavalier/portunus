package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/de-lachende-cavalier/portunus/pkg/logger"
)

func init() {
	rootCmd.AddCommand(checkCmd)
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for expired SSH keys",
	Long:  `Check if any SSH keys have expired and need to be rotated or renewed.`,
	Run:   runCheckCmd,
}

// runCheckCmd checks for expired SSH keys
func runCheckCmd(cmd *cobra.Command, args []string) {
	logger.Info("Checking for expired keys...")

	// Get expired keys from config
	expiredKeys := appConfig.GetExpiredKeys()

	if len(expiredKeys) == 0 {
		logger.Info("No expired keys found")
		fmt.Println("[+] No expired keys found")
		return
	}

	// Display expired keys
	logger.Info("The following keys have expired:")
	fmt.Println("[+] The following keys have expired:")

	for _, key := range expiredKeys {
		keyConfig, exists := appConfig.Keys[key]
		if !exists {
			continue
		}

		expiredFor := time.Since(keyConfig.ExpiresAt).Round(time.Second)
		logger.Infof("- %s (expired %s ago)", key, expiredFor)
		fmt.Printf("\t[+] %s (expired %s ago)\n", key, expiredFor)
	}

	// Provide instructions
	fmt.Println("\n[+] To rotate expired keys, run:")
	fmt.Println("\tportunus rotate -t <duration> -p <password>")
	fmt.Println("\n[+] To renew expired keys, run:")
	fmt.Println("\tportunus renew -t <duration>")
}
