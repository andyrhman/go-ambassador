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

func Stats(c *fiber.Ctx) error {
	id, _ := middlewares.GetUserId(c)

	uid, _ := uuid.Parse(id)

	var links []models.Link

	db.DB.Find(&links, models.Link{
		UserId: uid,
	})

	var result []interface{}

	var orders []models.Order

	for _, link := range links {
		db.DB.Preload("OrderItems").Find(&orders, &models.Order{
			Code:     link.Code,
			Complete: true,
		})

		revenue := 0.0

		for _, order := range orders {
			revenue += order.GetTotal()
		}

		result = append(result, fiber.Map{
			"code":    link.Code,
			"count":   len(orders),
			"revenue": revenue,
		})
	}

	return c.JSON(result)
}

func GetLink(c *fiber.Ctx) error {
	code := c.Params("code")

	var link models.Link

	if err := db.DB.Preload("User").Preload("Products").Find(&link, &models.Link{
		Code: code,
	}).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Code not found"})
	}

	return c.JSON(link)
}
