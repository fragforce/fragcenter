package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// srtsideCmd represents the srtside command
var srtsideCmd = &cobra.Command{
	Use:   "srtside",
	Short: "Start SRT helper daemon",
	Long:  `TODO: Fill in`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("srtside called")
	},
}

func init() {
	runCmd.AddCommand(srtsideCmd)
}
