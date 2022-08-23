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

}
