package redis

import (
	"identity-v2/cmd/config"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(conf config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     conf.Redis.URI,
		Password: conf.Redis.Password,
		DB:       0,
	})
}
