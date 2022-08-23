package cmd

import (
	"github.com/fragforce/fragcenter/lib/brain/core"
	"github.com/fragforce/fragcenter/lib/logs"
	"github.com/spf13/cobra"
)

// brainCmd represents the brain command
var brainCmd = &cobra.Command{
	Use:   "brain",
	Short: "Start the brain worker",
	Long:  `TODO: Fill this in`,
	Run: func(cmd *cobra.Command, args []string) {
		log := logs.L(nil)

		brain, err := core.NewBrain(log)
		if err != nil {
			log.Panicln("Problem with initial brain setup")
		}
		log = log.WithField("brain.obj", brain)

		if err := brain.Start(); err != nil {
			log.WithError(err).Panicln("Problem with running Brain")
		}

		log.Info("All done!")
	},
}

func init() {
	runCmd.AddCommand(brainCmd)
}
