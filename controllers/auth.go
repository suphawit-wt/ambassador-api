package controllers

import (
	"ambassador/database"
	"ambassador/models"
	"ambassador/utils"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
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

	user := models.User{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		IsAmbassador: false,
	}

	if err := user.SetPassword(req.Password); err != nil {
		panic(err)
	}

	database.DB.Create(&user)

	return c.Status(201).JSON(fiber.Map{
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

	if err := user.VerifyPassword(req.Password); err != nil {
		return c.Status(401).JSON(fiber.Map{
			"message": "Email or Password is invalid.",
		})
	}

	payload := jwt.StandardClaims{
		Subject:   strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString([]byte("secret"))
	if err != nil {
		panic(err)
	}

	utils.SetCookie(c, "access_token", accessToken, time.Now().Add(time.Hour*24))

	return c.Status(200).JSON(fiber.Map{
		"message": "Login Successfully!",
	})
}

func User(c *fiber.Ctx) error {
	user := models.User{}

	userId, err := utils.GetUserIdFromToken(c)
	if err != nil {
		panic(err)
	}

	database.DB.Where("id = ?", userId).First(&user)

	return c.Status(200).JSON(user)
}

func Logout(c *fiber.Ctx) error {
	utils.ClearCookie(c, "access_token")

	return c.Status(200).JSON(fiber.Map{
		"message": "Logout Successfully!",
	})
}

func UpdateInfo(c *fiber.Ctx) error {
	req := models.UpdateInfoRequest{}

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

	user := models.User{
		Id:        userId,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
	}

	database.DB.Model(&user).Updates(&user)

	return c.Status(200).JSON(fiber.Map{
		"message": "Update User Info Successfully!",
	})
}

func UpdatePassword(c *fiber.Ctx) error {
	req := models.UpdatePasswordRequest{}

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

	user := models.User{
		Id: userId,
	}

	user.SetPassword(req.Password)

	database.DB.Model(&user).Updates(&user)

	return c.Status(200).JSON(fiber.Map{
		"message": "Update User Info Successfully!",
	})
}
