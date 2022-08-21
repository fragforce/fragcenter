package streams

import "github.com/go-redis/redis/v8"

type TrackedKey struct {
}

type EventStream struct {
	pool    *redis.ClusterClient
	tracked map[string]*TrackedKey
}
