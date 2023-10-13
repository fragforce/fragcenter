package auth

import (
	"github.com/fragforce/fragcenter/lib/brain/plugin"
	"github.com/voc/srtrelay/stream"
)

type SRTAuth struct {
	*plugin.BrainCell
}

func NewSRTAuth(bc *plugin.BrainCell) SRTAuth {
	return SRTAuth{
		BrainCell: bc,
	}
}

func (a SRTAuth) Authenticate(stream.StreamID) bool {
	// FIXME: Add stream auth in here!
	return true
}
