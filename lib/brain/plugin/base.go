package plugin

import (
	"github.com/fragforce/fragcenter/lib/logs"
	"github.com/sirupsen/logrus"
)

// BrainCells is the base for Brain Plugins
type BrainCells struct {
	logs.BLog
}

func NewBaics(log *logrus.Entry) (*BrainCells, error) {
	b := BrainCells{
		BLog: *logs.NewBLog(log),
	}

	return &b, nil
}
