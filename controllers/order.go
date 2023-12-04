package controllers

import (
	"ambassador/database"
	"ambassador/models"
	"context"
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/checkout/session"
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

	var lineItems []*stripe.CheckoutSessionLineItemParams

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

		lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name:        stripe.String(product.Title),
					Description: stripe.String(product.Description),
					Images:      []*string{stripe.String(product.Image)},
				},
				UnitAmount: stripe.Int64(100 * int64(product.Price)),
				Currency:   stripe.String("usd"),
			},
			Quantity: stripe.Int64(int64(reqProduct["quantity"])),
		})
	}

	stripe.Key = os.Getenv("STRIPE_KEY")

	params := stripe.CheckoutSessionParams{
		Mode:               stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL:         stripe.String("http://localhost:5000/success?source={CHECKOUT_SESSION_ID}"),
		CancelURL:          stripe.String("http://localhost:5000/error"),
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems:          lineItems,
	}

	source, err := session.New(&params)
	if err != nil {
		tx.Rollback()
		log.Printf("session.New: %v", err)
	}

	order.TransactionId = source.ID

	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	tx.Commit()

	return c.Status(200).JSON(source)
}

func CompleteOrder(c *fiber.Ctx) error {
	req := models.ConfirmOrderRequest{}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Bad Request",
		})
	}

	order := models.Order{}

	if result := database.DB.Preload("OrderItems").First(&order, models.Order{
		TransactionId: req.Source,
	}); result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "Not Found.",
		})
	}

	order.Complete = true
	database.DB.Save(&order)

	go func(order models.Order) {
		ambassadorRevenue := 0.0
		adminRevenue := 0.0

		for _, item := range order.OrderItems {
			ambassadorRevenue += item.AmbassadorRevenue
			adminRevenue += item.AdminRevenue
		}

		user := models.User{}
		user.Id = order.Id

		database.DB.First(&user)

		database.RedisClient.ZIncrBy(context.Background(), "rankings", ambassadorRevenue, user.Name())

		smtpHost := os.Getenv("SMTP_HOST")
		smtpPort := os.Getenv("SMTP_PORT")

		smtpAuth := smtp.PlainAuth("", os.Getenv("SMTP_USERNAME"), os.Getenv("SMTP_PASSWORD"), smtpHost)

		ambassadorMessage := []byte(fmt.Sprintf("You earned $%f from the link #%s", ambassadorRevenue, order.Code))

		smtp.SendMail(smtpHost+":"+smtpPort, smtpAuth, "no-reply@email.com", []string{order.AmbassadorEmail}, ambassadorMessage)

		adminMessage := []byte(fmt.Sprintf("Order #%d with a total of $%f has been completed", order.Id, adminRevenue))

		smtp.SendMail(smtpHost+":"+smtpPort, smtpAuth, "no-reply@email.com", []string{"admin@admin.com"}, adminMessage)
	}(order)

	return c.Status(200).JSON(fiber.Map{
		"message": "Order Confirm Successfully!",
	})
}
