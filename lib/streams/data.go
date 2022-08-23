package streams

import (
	"context"
	"github.com/fragforce/fragcenter/lib/logs"
	"github.com/fragforce/fragcenter/lib/msg"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type IsTrackedFunc func(message msg.Message) bool
type ClusterOption func(options *redis.ClusterOptions)

type EventStream struct {
	*logs.BLog
	pool    *redis.ClusterClient
	tracked map[string]IsTrackedFunc
}

func NewEventStream(log *logrus.Entry, addrs []string, o ...ClusterOption) (*EventStream, error) {
	opts := redis.ClusterOptions{
		Addrs: addrs,
	}

	// Apply options
	for _, fn := range o {
		fn(&opts)
	}

	c := redis.NewClusterClient(&opts)
	if err := c.Ping(context.Background()).Err(); err != nil {
		panic("Unable to connect to redis " + err.Error())
	}

	es := EventStream{
		BLog:    logs.NewBLog(log),
		pool:    nil,
		tracked: make(map[string]IsTrackedFunc),
	}

	return &es, nil
}
