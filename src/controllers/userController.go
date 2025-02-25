package controllers

import (
	"context"
	"go-ambassador/src/db"
	"go-ambassador/src/models"

	"github.com/gofiber/fiber/v2"
	"github.com/iancoleman/orderedmap"
	"github.com/redis/go-redis/v9"
)

func Ambassadors(c *fiber.Ctx) error {
	var users []models.User

	db.DB.Where("isambassador = ?", true).Find(&users)

	return c.JSON(users)
}

func Rankings(c *fiber.Ctx) error {
	rankings, err := db.Cache.ZRevRangeByScoreWithScores(context.Background(), "rankings", &redis.ZRangeBy{
		Min: "-inf",
		Max: "+inf",
	}).Result()

	if err != nil {
		return err
	}

	orderedResult := orderedmap.New()

	for _, ranking := range rankings {
		orderedResult.Set(ranking.Member.(string), ranking.Score)
	}

	return c.JSON(orderedResult)
}
