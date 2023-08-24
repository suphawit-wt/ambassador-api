package controllers

import (
	"ambassador/database"
	"ambassador/models"
	"ambassador/utils"

	"github.com/bxcodec/faker/v4"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func GetUserLinks(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Bad Request",
		})
	}

	links := []models.Link{}

	database.DB.Where("user_id = ?", id).Find(&links)

	for i, link := range links {
		orders := []models.Order{}

		database.DB.Where("code = ? AND complete = true", link.Code).Find(&orders)

		links[i].Orders = orders
	}

	return c.Status(200).JSON(links)
}

func CreateLink(c *fiber.Ctx) error {
	req := models.CreateLinkRequest{}

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

	userId, err := utils.GetUserIdFromToken(c)
	if err != nil {
		panic(err)
	}

	link := models.Link{
		UserId: userId,
		Code:   faker.Username(),
	}

	for _, productId := range req.Products {
		product := models.Product{}
		product.Id = uint(productId)
		link.Products = append(link.Products, product)
	}

	database.DB.Create(&link)

	return c.Status(201).JSON(fiber.Map{
		"message": "Created Link Successfully!",
	})
}

func GetStats(c *fiber.Ctx) error {
	userId, err := utils.GetUserIdFromToken(c)
	if err != nil {
		panic(err)
	}

	links := []models.Link{}

	database.DB.Find(&links, models.Link{
		UserId: userId,
	})

	var result []interface{}

	orders := []models.Order{}

	for _, link := range links {
		database.DB.Preload("OrderItems").Find(&orders, &models.Order{
			Code:     link.Code,
			Complete: true,
		})

		revenue := 0.0
		for _, order := range orders {
			revenue += order.GetTotal()
		}

		result = append(result, fiber.Map{
			"code":    link.Code,
			"count":   len(orders),
			"revenue": revenue,
		})
	}

	return c.Status(200).JSON(result)
}

func GetLink(c *fiber.Ctx) error {
	code := c.Params("code")

	link := models.Link{
		Code: code,
	}

	database.DB.Preload("User").Preload("Products").First(&link)

	return c.Status(200).JSON(link)
}
