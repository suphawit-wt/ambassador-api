package controllers

import (
	"ambassador/database"
	"ambassador/models"

	"github.com/go-playground/validator/v10"
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

func CreateOrder(c *fiber.Ctx) error {
	req := models.CreateOrderRequest{}

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

	link := models.Link{
		Code: req.Code,
	}

	if result := database.DB.Preload("User").First(&link); result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "Not Found.",
		})
	}

	order := models.Order{
		Code:            link.Code,
		UserId:          link.UserId,
		AmbassadorEmail: link.User.Email,
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		Email:           req.Email,
		Address:         req.Address,
		Country:         req.Country,
		City:            req.City,
		Zip:             req.Zip,
	}

	tx := database.DB.Begin()

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	for _, reqProduct := range req.Products {
		product := models.Product{}
		product.Id = uint(reqProduct["product_id"])
		database.DB.First(&product)

		total := product.Price * float64(reqProduct["quantity"])

		item := models.OrderItem{
			OrderId:           order.Id,
			ProductTitle:      product.Title,
			Price:             product.Price,
			Quantity:          uint(reqProduct["quantity"]),
			AmbassadorRevenue: 0.1 * total,
			AdminRevenue:      0.9 * total,
		}

		if err := tx.Create(&item).Error; err != nil {
			tx.Rollback()
			c.Status(500).JSON(fiber.Map{
				"message": "Internal Server Error",
			})
		}
	}

	tx.Commit()

	return c.Status(200).JSON(order)
}
