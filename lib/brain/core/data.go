package core

import (
	"github.com/fragforce/fragcenter/lib/brain/plugin"
	_ "github.com/fragforce/fragcenter/lib/brain/plugins"
	"github.com/fragforce/fragcenter/lib/logs"
	"github.com/fragforce/fragcenter/lib/streams"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Brain struct {
	logs.BLog
	stream *streams.EventStream // Our upstream connection
	cells  map[uuid.UUID]plugin.Cell
}

// NewBrain creates a new Brain object
func NewBrain(log *logrus.Entry) (*Brain, error) {
	stream, err := streams.NewEventStream(
		log,
		viper.GetStringSlice("redis.addrs"),
		streams.ClusterOptPassword(viper.GetString("redis.password")),
	)
	if err != nil {
		return nil, err
	}
	return &Brain{
		BLog:   *logs.NewBLog(log),
		stream: stream,
	}, nil
}

func init() {
	viper.SetDefault("redis.addrs", "127.0.0.1")
	viper.SetDefault("redis.password", "")
}
