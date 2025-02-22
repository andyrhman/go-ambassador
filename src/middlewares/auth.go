package middlewares

import (
	"errors"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func IsAuthenticated(c *fiber.Ctx) error {
	cookie := c.Cookies("user_session")

	token, err := jwt.ParseWithClaims(cookie, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_ACCESS")), nil
	})

	if err != nil || !token.Valid {
		return c.Status(400).JSON(fiber.Map{
			"message": "Unauthenticated",
		})
	}

	return c.Next()
}

func GetUserId(c *fiber.Ctx) (string, error) {
	cookie := c.Cookies("user_session")

	token, err := jwt.ParseWithClaims(cookie, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_ACCESS")), nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("could not parse claims")
	}

	// Ensure "id" exists and is a string
	id, ok := claims["id"].(string)
	if !ok {
		return "", errors.New("user ID not found in token")
	}

	return id, nil
}
