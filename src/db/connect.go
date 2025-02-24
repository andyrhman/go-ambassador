package db

import (
	"github.com/joho/godotenv"
	"go-ambassador/src/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {

	err := godotenv.Load()

	if err != nil {
		log.Println("Warning: No .env file found, using system environment variables")
	}

	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		log.Fatal("DATABASE_URL is not set in environment variables")
	}

	// Connect to PostgreSQL using GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	DB = db

	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	log.Println("Connected to the database successfully!")
}

func AutoMigrate() {
	DB.AutoMigrate(&models.User{}, &models.Product{}, &models.Link{}, &models.Order{}, &models.OrderItem{})
}
