package controllers

import (
	"ambassador/database"
	"ambassador/models"

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
