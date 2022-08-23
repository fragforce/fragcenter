package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// rtmpsideCmd represents the rtmpside command
var rtmpsideCmd = &cobra.Command{
	Use:   "rtmpside",
	Short: "Run the RTMP helper daemon",
	Long:  `TODO: Fill in`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("rtmpside called")
	},
}

func init() {
	runCmd.AddCommand(rtmpsideCmd)
}
