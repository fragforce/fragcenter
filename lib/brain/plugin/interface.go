package plugin

import (
	"github.com/fragforce/fragcenter/lib/logs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

type NewCellFunc func(log *logrus.Entry) (Cell, error)

type Cell interface {
	logs.LAble                      // Add .L()
	ViperKeyFragment() string       // The viper key fragment to use for config info
	Viper() *viper.Viper            // Viper sub to our config
	Name() string                   // Cell name
	Run(exitC chan os.Signal) error // Execute - Can block, return nil if not need
	CleanupIsDone() error           // Should block until this cell is done cleaning up
}
