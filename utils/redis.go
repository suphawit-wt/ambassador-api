package utils

import (
	"ambassador/database"
	"context"
)

func DeleteCache(key string) {
	database.RedisClient.Del(context.Background(), key)
}
