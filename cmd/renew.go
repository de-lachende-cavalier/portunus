package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/de-lachende-cavalier/portunus/pkg/config"
	"github.com/de-lachende-cavalier/portunus/pkg/logger"
)

var (
	renewTime      string
	renewKeySubset []string
)

func init() {
	rootCmd.AddCommand(renewCmd)

	renewCmd.Flags().StringVarP(&renewTime, "time", "t", "",
		"specifies for how much longer the key should be valid (format: <int><specifier>, where specifier is either s (seconds), m (minutes), h (hours) or d (days)")
	renewCmd.Flags().StringSliceVarP(&renewKeySubset, "subset", "s", []string{},
		"specifies the subset of keys you want to act on (if empty, acts on all expired keys)")

	renewCmd.MarkFlagRequired("time")
}

var renewCmd = &cobra.Command{
	Use:   "renew",
	Short: "Renew expired SSH keys",
	Long:  `Renew expired SSH keys by extending their expiration date.`,
	Run:   runRenewCmd,
}

// runRenewCmd handles the renewal of SSH keys
func runRenewCmd(cmd *cobra.Command, args []string) {
	logger.Info("Renewing keys...")
	fmt.Println("[+] Renewing keys...")

	// Parse the time duration
	duration, err := parseDuration(renewTime)
	if err != nil {
		logger.Fatal(err, "Failed to parse time duration")
	}

	// Get keys to renew
	var keysToRenew []string
	if len(renewKeySubset) > 0 {
		// Use specified subset of keys
		keysToRenew = renewKeySubset
	} else {
		// Get all expired keys
		keysToRenew = appConfig.GetExpiredKeys()
	}

	if len(keysToRenew) == 0 {
		logger.Info("No keys found to renew")
		fmt.Println("[+] No keys found to renew")
		return
	}

	// Renew keys
	now := time.Now()
	renewedCount := 0

	for _, key := range keysToRenew {
		keyConfig, exists := appConfig.Keys[key]
		if !exists {
			logger.Infof("Key not found in configuration: %s", key)
			fmt.Printf("[+] Key not found in configuration: %s\n", key)
			continue
		}

		// Update expiration time by creating a new KeyConfig
		newKeyConfig := config.KeyConfig{
			CreatedAt: keyConfig.CreatedAt,
			ExpiresAt: now.Add(duration),
		}
		appConfig.Keys[key] = newKeyConfig

		logger.Infof("Renewed key: %s (new expiration: %s)", key, newKeyConfig.ExpiresAt.Format(time.RFC3339))
		fmt.Printf("\t[+] %s renewed, new expiration date: %s\n", key, newKeyConfig.ExpiresAt.Format(time.RFC3339))
		renewedCount++
	}

	// Save configuration
	if err := appConfig.Save(cfgFile); err != nil {
		logger.Fatal(err, "Failed to save configuration")
	}

	logger.Infof("Successfully renewed %d keys", renewedCount)
	fmt.Printf("[+] The keys have been successfully renewed\n")
}
