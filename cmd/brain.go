package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// brainCmd represents the brain command
var brainCmd = &cobra.Command{
	Use:   "brain",
	Short: "Start the brain worker",
	Long:  `TODO: Fill this in`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("brain called")
	},
}

func init() {
	runCmd.AddCommand(brainCmd)
}
