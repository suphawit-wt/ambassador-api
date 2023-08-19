package utils

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func GetUserIdFromToken(c *fiber.Ctx) (uint, error) {
	cookie := c.Cookies("access_token")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return 0, err
	}

	payload := token.Claims.(*jwt.StandardClaims)

	userId, err := strconv.Atoi(payload.Subject)
	if err != nil {
		return 0, err
	}

	return uint(userId), nil
}

func NewCookie(name string, value string, expires time.Time) *fiber.Cookie {
	cookie := new(fiber.Cookie)
	cookie.Name = name
	cookie.Value = value
	cookie.HTTPOnly = true
	cookie.Expires = expires

	return cookie
}

func SetCookie(c *fiber.Ctx, name string, value string, expire time.Time) {
	c.Cookie(NewCookie(name, value, expire))
}

func ClearCookie(c *fiber.Ctx, name string) {
	c.Cookie(NewCookie(name, "", time.Now().Add(-time.Hour)))
}
