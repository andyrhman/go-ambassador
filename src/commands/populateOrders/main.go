package main

import (
	"log"

	"go-ambassador/src/db"
	"go-ambassador/src/models"

	"github.com/brianvoe/gofakeit/v7"
)

func main() {
	db.Connect()

	var users []models.User
	if err := db.DB.Find(&users).Error; err != nil {
		log.Fatalf("Error fetching users: %v", err)
	}

	if len(users) == 0 {
		log.Fatalf("No users found, cannot seed orders")
	}

	var links []models.Link
	if err := db.DB.Find(&links).Error; err != nil {
		log.Fatalf("Error fetching links: %v", err)
	}

	for i := 0; i < 30; i++ {
		user := users[i%len(users)]
		link := links[i%len(links)]

		order := models.Order{
			Code:            link.Code,
			AmbassadorEmail: gofakeit.Email(),
			UserId:          user.Id,
			FullName:        gofakeit.Name(),
			Email:           gofakeit.Email(),
			Address:         gofakeit.Street(),
			Country:         gofakeit.Country(),
			City:            gofakeit.City(),
			Zip:             gofakeit.Zip(),
			Complete:        true,
		}

		for j := 0; j < gofakeit.Number(1, 5); j++ {
			item := models.OrderItem{
				ProductTitle:      gofakeit.ProductName(),
				Price:             float64(gofakeit.Number(100000, 5000000)),
				Quantity:          uint(gofakeit.Number(1, 5)),
				AmbassadorRevenue: float64(gofakeit.Number(10000, 500000)),
				AdminRevenue:      float64(gofakeit.Number(1000, 50000)),
			}
			order.OrderItems = append(order.OrderItems, item)
		}

		if err := db.DB.Create(&order).Error; err != nil {
			log.Fatalf("Error seeding order: %v", err)
		}
	}

	log.Println("🌱 Seeding has been completed successfully!")
}
