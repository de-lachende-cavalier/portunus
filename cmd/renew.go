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

	Run: func(cmd *cobra.Command, args []string) {
		expData := make(map[string][2]time.Time)

		expData, err := librarian.ReadConfig()
		if err != nil {
			fmt.Println(err)
			return
		}

		delta_s, err := cmd.Flags().GetString("time")

		delta_i, err := parseTime(delta_s)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, times := range expData {
			times[1] = times[1].Add(time.Second * time.Duration(delta_i))
		}

		err = librarian.WriteConfig(expData)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("The keys have been succesfully renewed")
	},
}
