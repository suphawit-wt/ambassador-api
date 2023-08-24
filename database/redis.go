package database

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var RedisChannel chan string

func SetupRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})
}

func SetupRedisChannel() {
	RedisChannel = make(chan string)

	go func(ch chan string) {
		for {
			// time.Sleep(5 * time.Second)
			key := <-ch

			RedisClient.Del(context.Background(), key)

			fmt.Println("Cache Cleared" + key)
		}
	}(RedisChannel)
}

func ClearCache(keys ...string) {
	for _, key := range keys {
		RedisChannel <- key
	}
}
