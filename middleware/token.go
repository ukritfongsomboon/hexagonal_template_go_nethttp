package middleware

import (
	"hexagonal/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func ValidToken(c *fiber.Ctx) error {
	var access_token string
	cookie := c.Cookies("Accesstoken")

	authorizationHeader := c.Get("Authorization")
	fields := strings.Fields(authorizationHeader)

	if len(fields) != 0 && fields[0] == "Bearer" {
		access_token = fields[1]
	} else {
		access_token = cookie
	}

	if access_token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    fiber.StatusUnauthorized,
			"status":  false,
			"message": "unauthorized",
			"data":    "",
		})
	}

	sub, err := utils.ValidateToken(access_token, viper.GetString("app.access_token_public_key"))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    fiber.StatusUnauthorized,
			"status":  false,
			"message": err.Error(),
			"data":    "",
		})
	}

	_ = sub

	c.Locals("user_id", sub)
	return c.Next()

	// return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
	// 	"code":    401,
	// 	"status":  false,
	// 	"message": "unauthorized",
	// 	"data":    "",
	// })
}
