package main

import (
	"go-ambassador/src/db"
	"go-ambassador/src/models"
	"log"

	"github.com/brianvoe/gofakeit/v7"
)

func main() {
	db.Connect()

	for i := 0; i < 30; i++ {
		users := models.User{
			Fullname: gofakeit.Name(),
			Username: gofakeit.Username(),
			Email: gofakeit.Email(),
			Isambassador: true,
		}

		users.SetPassword("123123")

		if err := db.DB.Create(&users).Error; err != nil {
			log.Fatalf("Error seeding users: %v", err)
		}
	}

	log.Println("Seeding has been completed 🌱")
}