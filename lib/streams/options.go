package streams

import "github.com/go-redis/redis/v8"

func ClusterOptUsername(username string) ClusterOption {
	return func(options *redis.ClusterOptions) {
		options.Username = username
	}
}

func ClusterOptPassword(passwd string) ClusterOption {
	return func(options *redis.ClusterOptions) {
		options.Password = passwd
	}
}
