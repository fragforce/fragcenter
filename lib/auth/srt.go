package auth

import (
	"github.com/fragforce/fragcenter/lib/brain/plugin"
	"github.com/voc/srtrelay/stream"
)

type SRTAuth struct {
	bc *plugin.BrainCell // Not an actual plugin - just a helper struct

}

func NewSRTAuth(bc *plugin.BrainCell) SRTAuth {
	return SRTAuth{
		bc: bc,
	}
}

func (a SRTAuth) Authenticate(sid stream.StreamID) bool {
	l := a.bc.L(nil).WithField("stream-id", sid)
	l.Debug("Request to auth new SRT stream")

	// FIXME: Add stream auth in here!
	return true
}
