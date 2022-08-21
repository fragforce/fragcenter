package logs

import "github.com/sirupsen/logrus"

// L returns a logger just for your local use
func L(log *logrus.Entry) *logrus.Entry {
	if log == nil {
		log = baseLog
	}
	return log
}
