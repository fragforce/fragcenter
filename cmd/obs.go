package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// obsCmd represents the obs command
var obsCmd = &cobra.Command{
	Use:   "obs",
	Short: "Start the OBS helper",
	Long:  `TODO: Fill this in`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("obs called")
	},
}

func init() {
	runCmd.AddCommand(obsCmd)
}
