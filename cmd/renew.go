package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/mowzhja/portunus/librarian"
)

func init() {
	rootCmd.AddCommand(renewCmd)

	renewCmd.Flags().StringP("time", "t", "",
		"Specify for how much longer they key should be valid (format: -t <int><specifier>, where specifier is either s (seconds), m (minutes), h (hours) or d (days)")

	renewCmd.MarkFlagRequired("time")
}

var renewCmd = &cobra.Command{
	Use:   "renew",
	Short: "renews the expiry on all the tracked keys",
	Long: `There's no need to specify which keys
  to renew: it automatically renews all the ones it tracks.`,
	Run: runRenewCmd,
}

// Helper function to use instead of the default anonymous function associated with Command.Run().
func runRenewCmd(cmd *cobra.Command, args []string) {
	expData := make(map[string][2]time.Time)

	oldData, err := librarian.ReadConfig()
	handleErr(err)

	delta_s, err := cmd.Flags().GetString("time")
	handleErr(err)
	delta_i, err := parseTime(delta_s)
	handleErr(err)

	for keyFile, times := range oldData {
		n_times := [2]time.Time{}
		n_times[0] = times[0].Round(0)
		n_times[1] = times[1].Add(time.Second * time.Duration(delta_i)).Round(0)

		expData[keyFile] = n_times
	}

	err = librarian.WriteConfig(expData)
	handleErr(err)

	fmt.Println("The keys have been succesfully renewed")
}
