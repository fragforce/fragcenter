package logs

import "github.com/sirupsen/logrus"

type logAugFunc func(string) logrus.Fields

type LAble interface {
	L(log *logrus.Entry) *logrus.Entry
}

type BLog struct {
	bLog   *logrus.Entry         // Log object - doesn't have our custom log stuff added - that's done via L()
	logAug map[string]logAugFunc // Func to get anything to add to the loggers
}

// NewBLog creates a new log helper
func NewBLog(log *logrus.Entry) *BLog {
	return &BLog{
		bLog: log,
	}
}

// L returns a logger customized for BLog - if passed a logger it bases it on it rather than BLog's built in logger
func (b BLog) L(log *logrus.Entry) *logrus.Entry {
	if log == nil {
		log = b.bLog
	}

	for name, fn := range b.logAug {
		log = log.WithFields(fn(name))
	}

	return log
}
