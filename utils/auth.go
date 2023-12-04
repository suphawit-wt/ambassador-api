package utils

import (
	"ambassador/database"
	"ambassador/models"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

var JwtSecret string = os.Getenv("JWT_SECRET")

func GetUserIdFromToken(c *fiber.Ctx) (uint, error) {
	cookie := c.Cookies("access_token")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JwtSecret), nil
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

func CheckIsAmbassador(userId uint) (bool, error) {
	user := models.User{}

	if result := database.DB.First(&user, userId); result.RowsAffected == 0 {
		return false, &fiber.Error{
			Code: 404,
		}
	}

	if user.IsAmbassador {
		return true, nil
	} else {
		return false, nil
	}
}

func GenerateAccessToken(userId uint) (string, error) {
	payload := jwt.StandardClaims{
		Subject:   strconv.Itoa(int(userId)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString([]byte(JwtSecret))

	return accessToken, err
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
