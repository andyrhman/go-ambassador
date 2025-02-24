package middlewares

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func IsAuthenticated(c *fiber.Ctx) error {
	cookie := c.Cookies("user_session")

	token, err := jwt.ParseWithClaims(cookie, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_ACCESS")), nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Unauthenticated",
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid token claims",
		})
	}

	scope, ok := claims["scope"].(string)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid token scope",
		})
	}

	isAmbassador := strings.Contains(c.Path(), "/api/ambassador")

	if (scope == "admin" && isAmbassador) || (scope == "ambassador" && !isAmbassador) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthorized",
		})
	}

	return c.Next()
}

func GenerateJwt(userId string, scope string) (string, error) {
	claims := jwt.MapClaims{
		"id":    userId,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
		"scope": scope,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET_ACCESS")))
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

	id, ok := claims["id"].(string)
	if !ok {
		return "", errors.New("user ID not found in token")
	}

	return id, nil
}
