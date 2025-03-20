package main

import (
	"go-ambassador/src/db"
	"go-ambassador/src/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	db.Connect()
	db.AutoMigrate()
	db.SetupRedis()
	db.SetupCacheChannel()

	app := fiber.New()
    // this is the cors
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "http://localhost:3000",
	}))

	routes.Setup(app)

	app.Listen(":8000")
}
