package main

import (
	"fmt"
	"go-ambassador/src/db"
	"go-ambassador/src/models"
	"log"

	"github.com/brianvoe/gofakeit/v7"
)

func main() {
	db.Connect()

	for i := 0; i < 30; i++ {
		imageURL := fmt.Sprintf("https://picsum.photos/seed/%s/200/300", gofakeit.UUID())
		products := models.Product{
			Title:       gofakeit.ProductName(),
			Description: gofakeit.ProductDescription(),
			Image:       imageURL,
			Price:       gofakeit.Product().Price,
		}

		if err := db.DB.Create(&products).Error; err != nil {
			log.Fatalf("Error seeding products: %v", err)
		}
	}

	log.Println("Seeding has been completed 🌱")
}
