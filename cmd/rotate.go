package cmd

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/de-lachende-cavalier/portunus/pkg/keys"
	"github.com/de-lachende-cavalier/portunus/pkg/logger"
)

var (
	rotateCipher    string
	rotateTime      string
	rotatePassword  string
	rotateKeySubset []string
)

func init() {
	rootCmd.AddCommand(rotateCmd)

	rotateCmd.Flags().StringVarP(&rotateCipher, "cipher", "c", "ed25519",
		"specifies which cipher to use for key generation")
	rotateCmd.Flags().StringVarP(&rotateTime, "time", "t", "",
		"specifies for how much longer the key should be valid (format: <int><specifier>, where specifier is either s (seconds), m (minutes), h (hours) or d (days)")
	rotateCmd.Flags().StringVarP(&rotatePassword, "password", "p", "",
		"specifies the password to use with ssh-keygen (NOTE: this password is used for ALL the keys that are rotated)")
	rotateCmd.Flags().StringSliceVarP(&rotateKeySubset, "subset", "s", []string{},
		"specifies the subset of keys you want to act on (if empty, acts on all keys in ~/.ssh)")

	rotateCmd.MarkFlagRequired("time")
	rotateCmd.MarkFlagRequired("password")
}

var rotateCmd = &cobra.Command{
	Use:   "rotate",
	Short: "Rotate SSH keys",
	Long: `Rotate SSH keys by deleting old ones and creating new ones with the same name.
If called without the subset flag, this command rotates ALL the keys in ~/.ssh/.
The new keys will be tracked with their expiration dates.`,
	Run: runRotateCmd,
}

// runRotateCmd handles the rotation of SSH keys
func runRotateCmd(cmd *cobra.Command, args []string) {
	logger.Info("Rotating keys...")
	fmt.Println("[+] Rotating keys...")

	// Parse the time duration
	duration, err := parseDuration(rotateTime)
	if err != nil {
		logger.Fatal(err, "Failed to parse time duration")
	}

	// Create key manager
	keyManager, err := keys.NewManager()
	if err != nil {
		logger.Fatal(err, "Failed to create key manager")
	}

	// Get keys to rotate
	var keyPaths []string
	if len(rotateKeySubset) > 0 {
		// Use specified subset of keys
		homeDir, err := cmd.Flags().GetString("home")
		if err != nil {
			homeDir, err = getHomeDir()
			if err != nil {
				logger.Fatal(err, "Failed to get home directory")
			}
		}

		for _, key := range rotateKeySubset {
			// If the key doesn't have a path, assume it's in ~/.ssh/
			if !filepath.IsAbs(key) && !strings.HasPrefix(key, ".") {
				key = filepath.Join(homeDir, ".ssh", key)
			}
			keyPaths = append(keyPaths, key)
		}
	} else {
		// Get all keys
		keyPaths, err = keyManager.GetAllKeys(rootContext)
		if err != nil {
			logger.Fatal(err, "Failed to get SSH keys")
		}
	}

	if len(keyPaths) == 0 {
		logger.Info("No keys found to rotate")
		fmt.Println("[+] No keys found to rotate")
		return
	}

	// Rotate keys
	creationTimes, err := keyManager.RotateKeys(rootContext, keyPaths, rotateCipher, rotatePassword)
	if err != nil {
		logger.Fatal(err, "Failed to rotate keys")
	}

	// Update configuration
	for path, creationTime := range creationTimes {
		expirationTime := creationTime.Add(duration)
		appConfig.AddKey(path, creationTime, expirationTime)

		logger.Infof("Rotated key: %s (expires: %s)", path, expirationTime.Format(time.RFC3339))
		fmt.Printf("\t[+] %s rotated, expiration date: %s\n", path, expirationTime.Format(time.RFC3339))
	}

	// Save configuration
	if err := appConfig.Save(cfgFile); err != nil {
		logger.Fatal(err, "Failed to save configuration")
	}

	logger.Info("Keys have been successfully rotated")
	fmt.Println("[+] The keys have been successfully rotated")
}

// parseDuration parses a duration string in the format "<int><specifier>"
func parseDuration(s string) (time.Duration, error) {
	// Check if the duration ends with "d" for days
	if len(s) > 0 && s[len(s)-1] == 'd' {
		// Parse the number of days
		days, err := strconv.Atoi(s[:len(s)-1])
		if err != nil {
			return 0, fmt.Errorf("invalid days value: %w", err)
		}
		// Convert days to hours (24 hours per day)
		return time.Duration(days) * 24 * time.Hour, nil
	}

	// Use standard time.ParseDuration for other units
	return time.ParseDuration(s)
}

// getHomeDir returns the user's home directory
func getHomeDir() (string, error) {
	home, err := filepath.Abs(filepath.Dir(filepath.Join(filepath.Dir("~"), "..")))
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return home, nil
}
