package cmd

import (
	"github.com/fragforce/fragcenter/lib/brain/core"
	"github.com/fragforce/fragcenter/lib/logs"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start worker",
	Run: func(cmd *cobra.Command, args []string) {
		log := logs.L(nil)

		brain, err := core.NewBrain(log)
		if err != nil {
			log.WithError(err).Panic("Problem with initial brain setup")
		}
		log = log.WithField("brain.obj", brain)

		if err := brain.Start(); err != nil {
			log.WithError(err).Panicln("Problem with running Brain")
		}

		log.Info("All done!")
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
