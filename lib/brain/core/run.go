package core

import (
	"github.com/fragforce/fragcenter/lib/brain/plugin"
	"os"
	"os/signal"
)

// Start begins execution
func (b Brain) Start() error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	exits := make(map[string]chan os.Signal, len(b.cells))

	for uid, cell := range b.cells {
		go func(uid string, cell plugin.Cell) {
			l := cell.L(nil).WithField("uid", uid)
			exits[uid] = make(chan os.Signal, 1)

			if err := cell.Run(exits[uid]); err != nil {
				l.WithError(err).Error("Failed running cell")
			}
		}(uid.String(), cell)
	}

	for {
		for sig := range c {
			l := b.L(nil).WithField("signal", sig)
			l.Warn("Received signal to exit")

			// Send cleanup notes
			for uid, ch := range exits {
				l := l.WithField("uid", uid)
				l.Debug("Sending exit signal to brain cell")
				go func(ch chan os.Signal) {
					ch <- sig
				}(ch)
			}
			l.Info("Exiting...")
			return nil
		}
	}
}
