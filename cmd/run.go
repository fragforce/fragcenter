package cmd

import (
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start workers",
}

func init() {
	rootCmd.AddCommand(runCmd)
}
