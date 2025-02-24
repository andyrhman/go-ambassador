package controllers

import (
	"go-ambassador/src/db"
	"go-ambassador/src/models"

	"github.com/gofiber/fiber/v2"
)

func Orders(c *fiber.Ctx) error {
	var orders []models.Order

	db.DB.Preload("OrderItems").Find(&orders)

	for i, order := range orders {
		orders[i].Total = order.GetTotal()
	}

	return c.JSON(orders)
}
