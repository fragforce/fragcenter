package core

import (
	"context"
	"github.com/fragforce/fragcenter/lib/brain/plugin"
	"os"
	"os/signal"
)

// Start begins execution
func (b Brain) Start() error {
	l := b.L(nil)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancelFunc := context.WithCancel(context.Background())

	l.Debug("Starting cells")
	for uid, cell := range b.cells {
		l := l.WithField("cell", cell)
		l.Trace("Starting cell")
		go func(uid string, cell plugin.Cell) {
			l := cell.L(nil).WithField("uid", uid)

			if err := cell.Run(
				context.WithValue(
					context.WithValue(
						ctx,
						"name",
						cell.Name(),
					),
					"uid",
					uid,
				),
			); err != nil {
				l.WithError(err).Error("Failed running cell")
			}
		}(uid.String(), cell)
	}

	for {
		l.Trace("In wait loop")
		select {
		case <-ctx.Done():
			l.Trace("CTX is done")
			cancelFunc()
			l.Trace("Cancel done")
			return nil
		case <-c:
			l.Trace("Calling cancel")
			cancelFunc()
			l.Trace("Cancel called")
		}
	}
}
