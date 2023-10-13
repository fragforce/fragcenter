package plugins

import (
	"context"
	"github.com/fragforce/fragcenter/lib/brain/plugin"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// Drinks is a plugin to add a web api/ui
type Drinks struct {
	*plugin.BrainCell
	Gin *gin.Engine
	Srv *http.Server
}

func NewDrinks(log *logrus.Entry) (plugin.Cell, error) {
	bc, err := plugin.NewBrainCell(log, "drinks")
	if err != nil {
		return nil, err
	}

	v := bc.Viper()
	v.SetDefault("listen", "0.0.0.0:80")

	d := Drinks{
		BrainCell: bc,
		Gin:       gin.Default(),
	}
	d.Srv = &http.Server{
		Addr:    v.GetString("listen"),
		Handler: d.Gin,
	}

	return &d, nil
}

func init() {
	// Register the plugin creator - it gets created when system is ready for it
	plugin.Register(NewDrinks)
}

func (d *Drinks) Run(exitC chan os.Signal) error {
	return d.Srv.ListenAndServe()
}

func (d *Drinks) CleanupIsDone() error {
	return d.Srv.Shutdown(context.Background())
}
