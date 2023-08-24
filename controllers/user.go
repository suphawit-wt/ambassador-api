package controllers

import (
	"ambassador/database"
	"ambassador/models"
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

func GetAllAmbassador(c *fiber.Ctx) error {
	users := []models.User{}

	database.DB.Where("is_ambassador = true").Find(&users)

	return c.Status(200).JSON(users)
}

func GetRankings(c *fiber.Ctx) error {
	rankings, err := database.RedisClient.ZRevRangeByScoreWithScores(context.Background(), "rankings", &redis.ZRangeBy{
		Min: "-inf",
		Max: "+inf",
	}).Result()
	if err != nil {
		panic(err)
	}

	result := make(map[string]float64)

	for _, ranking := range rankings {
		result[ranking.Member.(string)] = ranking.Score
	}

	return c.Status(200).JSON(result)
}
