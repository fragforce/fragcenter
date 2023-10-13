package plugins

import (
	"context"
	"github.com/fragforce/fragcenter/lib/auth"
	"github.com/fragforce/fragcenter/lib/brain/plugin"
	"github.com/haivision/srtgo"
	"github.com/sirupsen/logrus"
	"github.com/voc/srtrelay/relay"
	"github.com/voc/srtrelay/srt"
)

// SRT is a plugin to add a web api/ui
type SRT struct {
	*plugin.BrainCell
	config *srt.Config
	Srv    srt.Server
	ctx    context.Context
}

func NewSRT(log *logrus.Entry) (plugin.Cell, error) {
	bc, err := plugin.NewBrainCell(log, "SRT")
	if err != nil {
		return nil, err
	}

	v := bc.Viper()
	v.SetDefault("listen", "0.0.0.0:80")
	v.SetDefault("public-addy", "") //FIXME: Set these
	v.SetDefault("latency", "")
	v.SetDefault("loss-max-ttl", "")
	v.SetDefault("sync-clients", "")
	v.SetDefault("relay-buff-size", "")

	s := SRT{
		BrainCell: bc,
		config: &srt.Config{
			Server: srt.ServerConfig{
				Addresses:     v.GetStringSlice("listen"),
				PublicAddress: v.GetString("public-addy"),
				Latency:       v.GetUint("latency"),
				LossMaxTTL:    v.GetUint("loss-max-ttl"),
				SyncClients:   v.GetBool("sync-clients"),
				Auth:          auth.NewSRTAuth(bc),
			},
			Relay: relay.RelayConfig{
				Buffersize: v.GetUint("relay-buff-size"),
			},
		},
	}

	srtgo.InitSRT()
	s.Srv = srt.NewServer(s.config)

	return &s, nil
}

func init() {
	// Register the plugin creator - it gets created when system is ready for it
	plugin.Register(NewSRT)
}

func (s *SRT) Run(ctx context.Context) error {
	s.ctx = ctx
	if err := s.Srv.Listen(ctx); err != nil {
		return err
	}
	return nil
}
