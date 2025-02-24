package controllers

import (
	"go-ambassador/src/db"
	"go-ambassador/src/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func Products(c *fiber.Ctx) error {
	var products []models.Product

	db.DB.Model(&products).Find(&products)

	return c.JSON(products)
}

func CreateProduct(c *fiber.Ctx) error {
	var data map[string]any

	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid body request!",
		})
	}

	product := models.Product{
		Title:       data["title"].(string),
		Description: data["description"].(string),
		Image:       data["image"].(string),
		Price:       data["price"].(float64),
	}

	db.DB.Create(&product)

	return c.JSON(product)
}

func UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	uid, _ := uuid.Parse(id)

	product := models.Product{}
	product.Id = uid

	if err := c.BodyParser(&product); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid body request!",
		})
	}

	db.DB.Model(&product).Updates(&product)

	return c.JSON(product)
}

func GetProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	uid, _ := uuid.Parse(id)

	var product models.Product

	db.DB.Where("id = ?", uid).First(&product)

	return c.JSON(product)
}

func DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	uid, _ := uuid.Parse(id)

	var product models.Product

	db.DB.Where("id = ?", uid).Delete(product)

	return c.Status(204).JSON(nil)
}
