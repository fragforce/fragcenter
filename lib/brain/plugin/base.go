package plugin

import (
	"github.com/fragforce/fragcenter/lib/logs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// BrainCell is the base for Brain Plugins
type BrainCell struct {
	logs.BLog
	name string
}

func NewBrainCell(log *logrus.Entry, name string) (*BrainCell, error) {
	// TODO: Make sure 'name' doesn't have anything but letters, numbers, and dashes in it
	// Plus is viper key part valid

	b := BrainCell{
		BLog: *logs.NewBLog(log),
		name: name,
	}

	// TODO: Move out of closures to improve mem usage and overall perf
	b.AddLogAugment("brain.cell.name", func(fname string) logrus.Fields {
		return logrus.Fields{
			fname: name,
		}
	})

	return &b, nil
}

// Viper sub to our config
func (c *BrainCell) Viper() *viper.Viper {
	return viper.Sub("brain.plugin." + c.ViperKeyFragment())
}

// ViperKeyFragment returns the viper key fragment to use for config info
func (c *BrainCell) ViperKeyFragment() string {
	return c.name
}

// Name is just the name of the plugin/cell - Must be unique for all plugins
func (c *BrainCell) Name() string {
	return c.name
}
