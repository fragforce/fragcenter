package streams

import "github.com/redis/go-redis/v9"

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
