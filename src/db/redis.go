package db

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var Cache *redis.Client

func SetupRedis() {
	Cache = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		DB:   0,
	})

	ctx := context.Background()
	if err := Cache.Ping(ctx).Err(); err != nil {
		log.Printf("Error connecting to Redis: %v", err)
	} else {
		log.Println("Redis is connected")
	}
}
