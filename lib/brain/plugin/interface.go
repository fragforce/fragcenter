package plugin

import (
	"github.com/fragforce/fragcenter/lib/logs"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type NewCellFunc func(log *logrus.Entry) (Cell, error)

type Cell interface {
	logs.LAble                // Add .L()
	GUID() (uuid.UUID, error) // Returns the GUID
	ViperKeyFragment() string // The viper key fragment to use for config info
	Viper() *viper.Viper      // Viper sub to our config
	Name() string             // Cell name
}
