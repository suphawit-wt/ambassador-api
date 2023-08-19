package controllers

import (
	"ambassador/database"
	"ambassador/models"

	"github.com/gofiber/fiber/v2"
)

func GetAllOrders(c *fiber.Ctx) error {
	orders := []models.Order{}

	database.DB.Preload("OrderItems").Find(&orders)

	for i, order := range orders {
		orders[i].Name = order.FullName()
		orders[i].Total = order.GetTotal()
	}

	return c.Status(200).JSON(orders)
}
