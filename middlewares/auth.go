package middlewares

import (
	"ambassador/utils"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

var JwtSecret string = os.Getenv("JWT_SECRET")

func IsAmbassador(c *fiber.Ctx) error {
	accessTokenCookie := c.Cookies("access_token")

	token, err := jwt.ParseWithClaims(accessTokenCookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JwtSecret), nil
	})
	if err != nil || !token.Valid {
		return c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized.",
		})
	}

	payload := token.Claims.(*jwt.StandardClaims)

	userId, err := strconv.Atoi(payload.Subject)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized.",
		})
	}

	isAmbassador, err := utils.CheckIsAmbassador(uint(userId))
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized.",
		})
	}

	if !isAmbassador {
		return c.Status(403).JSON(fiber.Map{
			"message": "Forbidden.",
		})
	}

	return c.Next()
}

func IsAdmin(c *fiber.Ctx) error {
	accessTokenCookie := c.Cookies("access_token")

	token, err := jwt.ParseWithClaims(accessTokenCookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JwtSecret), nil
	})
	if err != nil || !token.Valid {
		return c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized.",
		})
	}

	payload := token.Claims.(*jwt.StandardClaims)

	userId, err := strconv.Atoi(payload.Subject)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized.",
		})
	}

	isAmbassador, err := utils.CheckIsAmbassador(uint(userId))
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized.",
		})
	}

	if isAmbassador {
		return c.Status(403).JSON(fiber.Map{
			"message": "Forbidden.",
		})
	}

	return c.Next()
}
