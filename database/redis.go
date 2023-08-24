package database

import "github.com/go-redis/redis/v8"

var RedisClient *redis.Client

func SetupRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})
}
