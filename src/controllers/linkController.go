package controllers

import (
	"go-ambassador/src/db"
	"go-ambassador/src/middlewares"
	"go-ambassador/src/models"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func Links(c *fiber.Ctx) error {
	id := c.Params("id")

	uid, _ := uuid.Parse(id)

	var links []models.Link

	db.DB.Where("user_id = ?", uid).Find(&links)

	for i, link := range links {
		var orders []models.Order

		db.DB.Where("code = ? and complete = true", link.Code).Find(&orders)

		links[i].Orders = orders
	}

	return c.JSON(links)
}

func CreateLinks(c *fiber.Ctx) error {
	type ProductsRequest struct {
		Products []uuid.UUID
	}

	id, _ := middlewares.GetUserId(c)

	uid, _ := uuid.Parse(id)
	var req ProductsRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid body request!",
		})
	}

	link := models.Link{
		UserId: uid,
		Code:   gofakeit.LetterN(7),
	}

	for _, productId := range req.Products {
		product := models.Product{}
		product.Id = productId
		link.Products = append(link.Products, product)
	}

	db.DB.Create(&link)

	return c.JSON(link)
}

