package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
)

var (
	RedisClient *redis.Client
)

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: "redis:6379", //redis port
	})
	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
}
