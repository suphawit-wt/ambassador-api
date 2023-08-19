package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func IsAuthenticated(c *fiber.Ctx) error {
	accessTokenCookie := c.Cookies("access_token")

	token, err := jwt.ParseWithClaims(accessTokenCookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil || !token.Valid {
		return c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized.",
		})
	}

	return c.Next()
}
