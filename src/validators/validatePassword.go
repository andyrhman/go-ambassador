package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type PasswordRequest struct {
	Password        string `json:"password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

func ValidatePassword(req PasswordRequest) *fiber.Error {
	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			switch e.Tag() {
			case "eqfield":
				return fiber.NewError(fiber.StatusBadRequest, "Password dan confirm password tidak sama")
			case "required":
				if e.Field() == "Password" {
					return fiber.NewError(fiber.StatusBadRequest, "Password tidak boleh kosong")
				}
				if e.Field() == "ConfirmPassword" {
					return fiber.NewError(fiber.StatusBadRequest, "Confirm password tidak boleh kosong")
				}
			case "min":
				return fiber.NewError(fiber.StatusBadRequest, "Password harus minimal 8 karakter")
			default:
				return fiber.NewError(fiber.StatusBadRequest, e.Error())
			}
		}
	}
	return nil
}
