package plugins

import (
	"github.com/fragforce/fragcenter/lib/brain/plugin"
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

// Drinks is a plugin to add a web api/ui
type Drinks struct {
	*plugin.BrainCell
	Gin *gin.Engine
}

func NewDrinks(log *logrus.Entry) (plugin.Cell, error) {
	bc, err := plugin.NewBrainCell(log, "drinks")
	if err != nil {
		return nil, err
	}
	return &Drinks{
		BrainCell: bc,
		Gin:       gin.Default(),
	}, nil
}

func init() {
	// Register the plugin creator - it gets created when system is ready for it
	plugin.Register(NewDrinks)
}

func (d *Drinks) Run() error {
	return d.Gin.Run(d.Viper().GetString("listen"))
}
