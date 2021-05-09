package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/mowzhja/portunus/librarian"
)

func init() {
	rootCmd.AddCommand(renewCmd)
}

var renewCmd = &cobra.Command{
	Use:   "renew",
	Short: "renews the expiry on all the tracked keys",
	Long: `The format it expects is exactly the same as the one
  used for the 'set' command. There's no need to specify which keys
  to renew: it automatically renews all the ones it tracks.`,
	// TODO maybe allow users to specify which keys to renew?

	Run: func(cmd *cobra.Command, args []string) {
		expData := make(map[string][2]time.Time)

		expData, err := librarian.ReadConfig()
		if err != nil {
			fmt.Println(err)
			return
		}

		toAdd, err := parseTime(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, times := range expData {
			times[1] = times[1].Add(time.Second * time.Duration(toAdd))
		}

		err = librarian.WriteConfig(expData)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("The keys have been succesfully renewed")
	},
}
