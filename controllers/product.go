package controllers

import (
	"ambassador/database"
	"ambassador/models"

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
