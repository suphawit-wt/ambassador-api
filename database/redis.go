package database

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var RedisChannel chan string

func SetupRedis() {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	RedisClient = redis.NewClient(&redis.Options{
		Addr: redisAddr,
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
