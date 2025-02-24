package controllers

import (
	"context"
	"encoding/json"
	"go-ambassador/src/db"
	"go-ambassador/src/models"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func Products(c *fiber.Ctx) error {
	var products []models.Product

	db.DB.Model(&products).Find(&products)

	return c.JSON(products)
}

func CreateProduct(c *fiber.Ctx) error {
	var data map[string]any

	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid body request!",
		})
	}

	product := models.Product{
		Title:       data["title"].(string),
		Description: data["description"].(string),
		Image:       data["image"].(string),
		Price:       data["price"].(float64),
	}

	db.DB.Create(&product)

	go db.ClearCache("products_frontend", "products_backend")

	return c.JSON(product)
}

func UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	uid, _ := uuid.Parse(id)

	product := models.Product{}
	product.Id = uid

	if err := c.BodyParser(&product); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid body request!",
		})
	}

	db.DB.Model(&product).Updates(&product)

	go db.ClearCache("products_frontend", "products_backend")
	
	return c.JSON(product)
}

func GetProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	uid, _ := uuid.Parse(id)

	var product models.Product

	db.DB.Where("id = ?", uid).First(&product)

	return c.JSON(product)
}

func DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	uid, _ := uuid.Parse(id)

	var product models.Product

	db.DB.Where("id = ?", uid).Delete(product)

	go db.ClearCache("products_frontend", "products_backend")

	return c.Status(204).JSON(nil)
}

func ProductsFrontend(c *fiber.Ctx) error {
	var products []models.Product
	var ctx = context.Background()

	result, err := db.Cache.Get(ctx, "products_frontend").Result()

	if err != nil {
		db.DB.Find(&products)

		bytes, err := json.Marshal(products)

		if err != nil {
			panic(err)
		}

		if errKey := db.Cache.Set(ctx, "products_frontend", bytes, 30*time.Minute).Err(); errKey != nil {
			panic(errKey)
		}
	} else {
		json.Unmarshal([]byte(result), &products)
	}

	return c.JSON(products)
}

func ProductsBackend(c *fiber.Ctx) error {
	var products []models.Product
	var ctx = context.Background()

	result, err := db.Cache.Get(ctx, "products_backend").Result()

	if err != nil {
		db.DB.Find(&products)

		bytes, err := json.Marshal(products)
		if err != nil {
			panic(err)
		}

		db.Cache.Set(ctx, "products_backend", bytes, 30*time.Minute)
	} else {
		json.Unmarshal([]byte(result), &products)
	}

	var searchedProducts []models.Product

	if s := c.Query("search"); s != "" {
		lower := strings.ToLower(s)
		for _, product := range products {
			if strings.Contains(strings.ToLower(product.Title), lower) || strings.Contains(strings.ToLower(product.Description), lower) {
				searchedProducts = append(searchedProducts, product)
			}
		}
	} else {
		searchedProducts = products
	}

	if sortParam := c.Query("sort"); sortParam != "" {
		sortLower := strings.ToLower(sortParam)
		if sortLower == "desc" {
			sort.Slice(searchedProducts, func(i, j int) bool {
				return searchedProducts[i].Price < searchedProducts[j].Price
			})
		} else if sortLower == "asc" {
			sort.Slice(searchedProducts, func(i, j int) bool {
				return searchedProducts[i].Price > searchedProducts[j].Price
			})
		}
	}

	var total = len(searchedProducts)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage := 9

	var data []models.Product

	if total <= page*perPage && total >= (page-1)*perPage {
		data = searchedProducts[(page-1)*perPage : total]
	} else if total >= page*perPage {
		data = searchedProducts[(page-1)*perPage : page*perPage]
	} else {
		data = []models.Product{}
	}

	return c.JSON(fiber.Map{
		"data": data,
		"meta": fiber.Map{
			"total":     total,
			"page":      page,
			"last_page": total/perPage + 1,
		},
	})
}
