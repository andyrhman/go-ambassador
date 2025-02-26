package controllers

import (
	"go-ambassador/src/db"
	"go-ambassador/src/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func Orders(c *fiber.Ctx) error {
	var orders []models.Order

	db.DB.Preload("OrderItems").Find(&orders)

	for i, order := range orders {
		orders[i].Total = order.GetTotal()
	}

	return c.JSON(orders)
}

type CreateOrderRequest struct {
	Code     string
	FullName string
	Email    string
	Address  string
	Country  string
	City     string
	Zip      string
	Products []map[string]string
}

func CreateOrder(c *fiber.Ctx) error {
	var orderRequest CreateOrderRequest

	if err := c.BodyParser(&orderRequest); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid body request!"})
	}

	link := models.Link{}
	if err := db.DB.Preload("User").Where("code = ?", orderRequest.Code).First(&link).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid Code!"})
	}
	

	order := models.Order{
		UserId:          link.UserId,
		Code:            link.Code,
		AmbassadorEmail: link.User.Email,
		FullName:        orderRequest.FullName,
		Email:           orderRequest.Email,
		Address:         orderRequest.Address,
		Country:         orderRequest.Country,
		City:            orderRequest.City,
		Zip:             orderRequest.Zip,
	}

	db.DB.Create(&order)

	for _, productRequest := range orderRequest.Products {
		product := models.Product{}
		product.Id = uuid.Must(uuid.Parse(productRequest["product_id"]))
		db.DB.First(&product)

		qty, err := strconv.Atoi(productRequest["quantity"])
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "Invalid quantity"})
		}

		total := product.Price * float64(qty)

		item := models.OrderItem{
			OrderId:           order.Id,
			ProductTitle:      product.Title,
			Price:             product.Price,
			Quantity:          uint(qty),
			AmbassadorRevenue: 0.1 * total,
			AdminRevenue:      0.9 * total,
		}

		db.DB.Create(&item)
	}

	return c.JSON(order)
}
