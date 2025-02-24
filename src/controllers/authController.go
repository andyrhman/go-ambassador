package controllers

import (
	"go-ambassador/src/db"
	"go-ambassador/src/middlewares"
	"go-ambassador/src/models"
	"go-ambassador/src/validators"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	if data["password"] != data["confirm_password"] {
		return c.Status(400).JSON(fiber.Map{
			"message": "Password do not match",
		})
	}

	user := models.User{
		Fullname:     data["fullname"],
		Username:     data["username"],
		Email:        data["email"],
		Isambassador: strings.Contains(c.Path(), "/api/ambassador"),
	}

	user.SetPassword(data["password"])

	if err := db.DB.Where("email = ?", data["email"]).First(&user).Error; err == nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Email already exists!",
		})
	}

	if err := db.DB.Where("username = ?", data["username"]).First(&user).Error; err == nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Username already exists!",
		})
	}

	db.DB.Create(&user)

	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	var user models.User

	if data["email"] != "" {
		if err := db.DB.Where("LOWER(email) = ?", strings.ToLower(data["email"])).First(&user).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Invalid credentials!",
			})
		}
	} else if data["username"] != "" {
		if err := db.DB.Where("LOWER(username) = ?", strings.ToLower(data["username"])).First(&user).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Invalid credentials!",
			})
		}
	} else {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Invalid credentials!",
		})
	}

	if !user.ComparePassword(data["password"]) {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid credentials!",
		})
	}

	isAmbassador := strings.Contains(c.Path(), "api/ambassador")

	var scope string
	if isAmbassador {
		scope = "ambassador"
	} else {
		scope = "admin"
	}

	// ! Prevent ambassador login into admin login endpoint
	if !isAmbassador && user.Isambassador {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	tokenString, err := middlewares.GenerateJwt(user.Id.String(), scope)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not login",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "user_session",
		Value:    tokenString,
		HTTPOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	return c.JSON(fiber.Map{
		"message": "Successfully Logged In!",
	})
}

func User(c *fiber.Ctx) error {
	id, _ := middlewares.GetUserId(c)

	var user models.User

	db.DB.Where("id = ?", id).First(&user)

	if strings.Contains(c.Path(), "api/ambassador") {
		ambassador := models.Ambassador(user)
		ambassador.CalculateRevenue(db.DB)
		return c.JSON(user)
	}

	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "user_session",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func UpdateInfo(c *fiber.Ctx) error {
	id, _ := middlewares.GetUserId(c)

	// * store the data parsed from the request body
	var input models.User
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse request body",
		})
	}

	// * store the existing user data fetched from the database
	var existingUser models.User

	db.DB.Where("id = ?", id).First(&existingUser)

	if input.Fullname != "" {
		existingUser.Fullname = input.Fullname
	}

	if input.Email != "" && input.Email != existingUser.Email {
		var existingUserByEmail models.User
		if err := db.DB.Where("email = ?", input.Email).First(&existingUserByEmail).Error; err == nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": "Email already exists",
			})
		}
		existingUser.Email = input.Email
	}

	if input.Username != "" && input.Username != existingUser.Username {
		var existingUserByUsername models.User
		if err := db.DB.Where("username = ?", input.Username).First(&existingUserByUsername).Error; err == nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": "Username already exists",
			})
		}
		existingUser.Username = input.Username
	}

	db.DB.Save(&existingUser)

	return c.JSON(existingUser)
}

func UpdatePassword(c *fiber.Ctx) error {
	id, _ := middlewares.GetUserId(c)

	uid, _ := uuid.Parse(id)

	var req validators.PasswordRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse request body",
		})
	}

	if err := validators.ValidatePassword(req); err != nil {
		return c.Status(err.Code).JSON(fiber.Map{
			"error": err.Message,
		})
	}

	user := models.User{}
	user.Id = uid

	user.SetPassword(req.Password)

	if err := db.DB.Model(&user).Updates(user).Error; err != nil {
		return c.JSON(fiber.Map{
			"message": "Cannot update password!",
		})
	}

	return c.JSON(user)
}
