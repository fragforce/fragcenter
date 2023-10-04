package streams

import (
	"context"
	"github.com/fragforce/fragcenter/lib/logs"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type ClusterOption func(options *redis.ClusterOptions)

type EventStream struct {
	logs.BLog
	pool     *redis.ClusterClient
	actReact []*ActionReaction
}

func NewEventStream(log *logrus.Entry, addrs []string, op ...ClusterOption) (*EventStream, error) {
	opts := redis.ClusterOptions{
		Addrs: addrs,
	}

	// Apply options
	for _, fn := range op {
		fn(&opts)
	}

	c := redis.NewClusterClient(&opts)
	if err := c.Ping(context.Background()).Err(); err != nil {
		log.WithError(err).Info("Couldn't ping Redis Cluster")
		return nil, err
	}

	es := EventStream{
		BLog:     *logs.NewBLog(log),
		pool:     nil,
		actReact: make([]*ActionReaction, 0),
	}

	return &es, nil
}
