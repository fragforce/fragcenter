package logs

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

var (
	rootLog *logrus.Logger
	baseLog *logrus.Entry
)

// InitLogging sets up logging - should only be called from the rootCmd's initConfig
// Warning: Panics on error!
func InitLogging(rootCmd *cobra.Command) {
	// FIXME: Add a check to make sure viper has been started
	rootLog = logrus.New()
	var log *logrus.Entry

	lvl, err := logrus.ParseLevel(viper.GetString("log.level"))
	if err != nil {
		panic("Bad log level: " + err.Error())
	}
	if viper.GetBool("debug") && lvl != logrus.TraceLevel {
		lvl = logrus.DebugLevel
	}
	rootLog.SetLevel(lvl)

	rootLog.SetFormatter(&logrus.JSONFormatter{
		DisableTimestamp:  false,
		DisableHTMLEscape: false,
	})
	rootLog.SetReportCaller(true)

	log = rootLog.WithFields(logrus.Fields{
		"app":               rootCmd.Name(),
		"log.level.current": lvl.String(),
	})

	for k, v := range viper.GetStringMapString("runtime") {
		if strings.Contains(strings.ToLower(k), "_cert") || strings.Contains(strings.ToLower(v), "certificate") {
			// has cert data or kafka info
			continue
		}
		log = log.WithField(k, v)
	}

	// Set global logger
	baseLog = log

	log.Info("Init'ed logging")
}
