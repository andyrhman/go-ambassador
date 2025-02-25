package main

import (
	"go-ambassador/src/db"
	"go-ambassador/src/models"
	"context"
	"github.com/redis/go-redis/v9"
)

func main() {
	db.Connect()
	db.SetupRedis()

	ctx := context.Background()

	var users []models.User

	db.DB.Find(&users, models.User{
		Isambassador: true,
	})

	for _, user := range users {
		ambassador := models.Ambassador(user)
		ambassador.CalculateRevenue(db.DB)

		db.Cache.ZAdd(ctx, "rankings", redis.Z{
			Score:  *ambassador.Revenue,
			Member: user.Fullname,
		})
	}
}