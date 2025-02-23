package controllers

import (
	"go-ambassador/src/db"
	"go-ambassador/src/models"

	"github.com/gofiber/fiber/v2"
)

func AllAmbassadors(c *fiber.Ctx) error {
	var users []models.User

	db.DB.Where("isambassador = ?", true).Find(&users)

	return c.JSON(users)
}
