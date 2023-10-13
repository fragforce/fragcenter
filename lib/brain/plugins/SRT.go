package plugins

import (
	"context"
	"github.com/fragforce/fragcenter/lib/auth"
	"github.com/fragforce/fragcenter/lib/brain/plugin"
	"github.com/haivision/srtgo"
	"github.com/sirupsen/logrus"
	"github.com/voc/srtrelay/relay"
	"github.com/voc/srtrelay/srt"
	"os"
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
	s := SRT{
		BrainCell: bc,
		config: &srt.Config{
			Server: srt.ServerConfig{
				Addresses:     bc.Viper().GetStringSlice("listen"),
				PublicAddress: bc.Viper().GetString("public-addy"),
				Latency:       bc.Viper().GetUint("latency"),
				LossMaxTTL:    bc.Viper().GetUint("loss-max-ttl"),
				SyncClients:   bc.Viper().GetBool("sync-clients"),
				Auth:          auth.NewSRTAuth(bc),
			},
			Relay: relay.RelayConfig{
				Buffersize: bc.Viper().GetUint("relay-buff-size"),
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

func (s *SRT) Run(exitC chan os.Signal) error {
	ctx, cancel := context.WithCancel(context.Background())
	s.ctx = ctx
	s.waitForExit(exitC, cancel)
	if err := s.Srv.Listen(ctx); err != nil {
		return err
	}
	return nil
}

func (s *SRT) CleanupIsDone() error {
	// FIXME: Fix this
	return nil
}

func (s *SRT) waitForExit(exitC chan os.Signal, cancelFunc context.CancelFunc) {
	go func(exitC chan os.Signal, cancelFunc context.CancelFunc) {
		for {
			select {
			case <-s.ctx.Done():
				return
			case <-exitC:
				cancelFunc()
			}
		}
	}(exitC, cancelFunc)
}
