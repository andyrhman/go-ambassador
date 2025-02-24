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

	var link []models.Link

	db.DB.Where("user_id = ?", uid).Find(&link)

	return c.JSON(link)
}
