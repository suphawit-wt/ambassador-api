package controllers

import (
	"ambassador/database"
	"ambassador/models"
	"ambassador/utils"
	"context"
	"encoding/json"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func GetAllProducts(c *fiber.Ctx) error {
	products := []models.Product{}

	database.DB.Find(&products)

	return c.Status(200).JSON(products)
}

func CreateProduct(c *fiber.Ctx) error {
	product := models.Product{}

	if err := c.BodyParser(&product); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Bad Request",
		})
	}

	validate := validator.New()
	if err := validate.Struct(product); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Bad Request",
		})
	}

	database.DB.Create(&product)

	return c.Status(201).JSON(fiber.Map{
		"message": "Created Product Successfully!",
	})
}

func GetProductById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Bad Request",
		})
	}

	product := models.Product{}

	if result := database.DB.First(&product, id); result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "Not Found.",
		})
	}

	return c.Status(200).JSON(product)
}

func UpdateProduct(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Bad Request",
		})
	}

	product := models.Product{}

	if result := database.DB.First(&product, id); result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "Not Found.",
		})
	}

	req := models.Product{}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Bad Request",
		})
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Bad Request",
		})
	}

	product.Title = req.Title
	product.Description = req.Description
	product.Image = req.Image
	product.Price = req.Price

	database.DB.Save(&product)

	go utils.DeleteCache("products_frontend")
	go utils.DeleteCache("products_backend")

	return c.Status(200).JSON(fiber.Map{
		"message": "Updated Product Successfully!",
	})
}

func DeleteProduct(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Bad Request",
		})
	}

	if result := database.DB.Delete(models.Product{}, id); result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "Not Found.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Deleted Product Successfully!",
	})
}

func GetProductsFrontend(c *fiber.Ctx) error {
	products := []models.Product{}
	ctx := context.Background()

	productsCache, err := database.RedisClient.Get(ctx, "products_frontend").Result()
	if err != nil {
		database.DB.Find(&products)

		productsBytes, err := json.Marshal(products)
		if err != nil {
			panic(err)
		}

		database.RedisClient.Set(ctx, "products_frontend", productsBytes, time.Minute*5)
	} else {
		json.Unmarshal([]byte(productsCache), &products)
	}

	return c.Status(200).JSON(products)
}

func GetProductsBackend(c *fiber.Ctx) error {
	products := []models.Product{}
	ctx := context.Background()

	productsCache, err := database.RedisClient.Get(ctx, "products_backend").Result()
	if err != nil {
		database.DB.Find(&products)

		productsBytes, err := json.Marshal(products)
		if err != nil {
			panic(err)
		}

		database.RedisClient.Set(ctx, "products_backend", productsBytes, time.Minute*5)
	} else {
		json.Unmarshal([]byte(productsCache), &products)
	}

	var searchedProducts []models.Product

	if s := c.Query("s"); s != "" {
		srcLower := strings.ToLower(s)
		for _, product := range products {
			if strings.Contains(strings.ToLower(product.Title), srcLower) || strings.Contains(strings.ToLower(product.Description), srcLower) {
				searchedProducts = append(searchedProducts, product)
			}
		}
	} else {
		searchedProducts = products
	}

	if sortQuery := c.Query("sort"); sortQuery != "" {
		sortLower := strings.ToLower(sortQuery)
		if sortLower == "asc" {
			sort.Slice(searchedProducts, func(i, j int) bool {
				return searchedProducts[i].Price < searchedProducts[j].Price
			})
		} else if sortLower == "desc" {
			sort.Slice(searchedProducts, func(i, j int) bool {
				return searchedProducts[i].Price > searchedProducts[j].Price
			})
		}
	}

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Bad Request",
		})
	}

	var total = len(searchedProducts)
	var data []models.Product
	perPage := 9

	if total <= page*perPage && total >= (page-1)*perPage {
		data = searchedProducts[(page-1)*perPage : total]
	} else if total >= page*perPage {
		data = searchedProducts[(page-1)*perPage : page*perPage]
	} else {
		data = []models.Product{}
	}

	return c.Status(200).JSON(fiber.Map{
		"data":      data,
		"total":     total,
		"page":      page,
		"last_page": total/perPage + 1,
	})
}
