package controllers

import (
	"ambassador/database"
	"ambassador/models"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	req := models.User{}

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

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	user := models.User{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		Password:     passwordHash,
		IsAmbassador: false,
	}

	database.DB.Create(&user)

	return c.JSON(fiber.Map{
		"message": "Register Successfully!",
	})
}

func Login(c *fiber.Ctx) error {
	req := models.LoginRequest{}

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

	user := models.User{}

	if result := database.DB.Where("email = ?", req.Email).First(&user); result.RowsAffected == 0 {
		return c.Status(401).JSON(fiber.Map{
			"message": "Email or Password is invalid.",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Login Successfully!",
	})
}
