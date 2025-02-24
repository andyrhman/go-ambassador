package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var Cache *redis.Client
var CacheChannel chan string

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

func SetupCacheChannel() {
	CacheChannel = make(chan string)

	go func(ch chan string) {
		for {
			time.Sleep(5 * time.Second)

			key := <-ch

			Cache.Del(context.Background(), key)

			fmt.Println("Cache cleared " + key)
		}
	}(CacheChannel)
}

func ClearCache(keys ...string) {
	for _, key := range keys {
		CacheChannel <- key
	}
}
