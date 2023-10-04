package plugins

import (
	"github.com/fragforce/fragcenter/lib/brain/plugin"
	"github.com/sirupsen/logrus"
)

// MessageLogger is a simple plugin to log all ingress messages
type MessageLogger struct {
	*plugin.BrainCell
}

func NewMessageLogger(log *logrus.Entry) (plugin.Cell, error) {
	bc, err := plugin.NewBrainCell(log, "message-logger")
	if err != nil {
		return nil, err
	}
	return &MessageLogger{
		BrainCell: bc,
	}, nil
}

func init() {
	// Register the plugin creator so it gets created when system is ready for it
	plugin.Register(NewMessageLogger)
}

func (l *MessageLogger) Run() error {
	return nil
}
