package kv

import (
	"fmt"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(addr string, port string) *RedisClient {
	return &RedisClient{
		Client: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", addr, port),
			Password: "",
			DB:       0,
		}),
	}
}
