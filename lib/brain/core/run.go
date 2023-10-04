package core

import "github.com/fragforce/fragcenter/lib/brain/plugin"

// Start begins execution
func (b Brain) Start() error {
	for uid, cell := range b.cells {
		go func(uid string, cell plugin.Cell) {
			l := cell.L(nil).WithField("uid", uid)
			if err := cell.Run(); err != nil {
				l.WithError(err).Error("Failed running cell")
			}
		}(uid.String(), cell)
	}
	return nil
}
