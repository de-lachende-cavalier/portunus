package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/de-lachende-cavalier/portunus/pkg/config"
	"github.com/de-lachende-cavalier/portunus/pkg/logger"
)

var (
	cfgFile     string
	logLevel    string
	prettyLogs  bool
	appConfig   *config.Config
	rootContext context.Context
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "portunus",
	Short: "portunus is a simple utility for managing SSH key expiration",
	Long: `portunus acts as middleware, handling the keys for you through ssh-keygen 
and keeping track of their specified expiration dates, which ssh-keygen cannot do. 
Once the keys have expired, portunus will prompt you to either rotate them 
(delete the old ones and make new ones) or to renew them 
(postpone their expiration date by some specified amount).`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Initialize context
		rootContext = context.Background()

		// Initialize logger
		logger.Init(logLevel, prettyLogs)

		// Load configuration
		var err error
		appConfig, err = config.Load(cfgFile)
		if err != nil {
			logger.Fatal(err, "Failed to load configuration")
		}

		// Clean up non-existent keys from config
		appConfig.CleanNonExistentKeys()
		if err := appConfig.Save(cfgFile); err != nil {
			logger.Error(err, "Failed to save configuration after cleanup")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		// The root command doesn't do anything by itself
		_ = cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Error(err, "Command execution failed")
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.portunus.json)")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "log level (debug, info, warn, error)")
	rootCmd.PersistentFlags().BoolVar(&prettyLogs, "pretty-logs", true, "enable pretty logging")
}
