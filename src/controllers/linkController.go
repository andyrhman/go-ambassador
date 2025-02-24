package controllers

import (
	"go-ambassador/src/db"
	"go-ambassador/src/models"

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
