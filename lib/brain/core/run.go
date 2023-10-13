package core

import (
	"context"
	"github.com/fragforce/fragcenter/lib/brain/plugin"
	"os"
	"os/signal"
)

// Start begins execution
func (b Brain) Start() error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancelFunc := context.WithCancel(context.Background())

	for uid, cell := range b.cells {
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
		select {
		case <-ctx.Done():
			cancelFunc()
			return nil
		case <-c:
			cancelFunc()
		}
	}
}
