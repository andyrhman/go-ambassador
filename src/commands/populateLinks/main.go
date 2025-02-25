package main

import (
	"go-ambassador/src/db"
	"go-ambassador/src/models"
	"log"

	"github.com/brianvoe/gofakeit/v7"
)

func main() {
	db.Connect()

	var users []models.User
	if err := db.DB.Find(&users).Error; err != nil {
		log.Fatalf("Error fetching users: %v", err)
	}

	for i := 0; i < 30; i++ {

		user := users[i%len(users)]

		links := models.Link{
			Code: gofakeit.LetterN(7),
			UserId: user.Id,
		}

		if err := db.DB.Create(&links).Error; err != nil {
			log.Fatalf("Error seeding products: %v", err)
		}
	}

	log.Println("Seeding has been completed 🌱")
}
