package internal

import (
	"github.com/gofiber/fiber/v2"
	"github.com/josephakayesi/cadana/exchange-1/infra/config"
)

var redis = config.NewRedis()
var db = config.NewDatabase()

func ValidateAPIKey(c *fiber.Ctx) error {
	token := c.Cookies("access_token")

	if t, _ := redis.Get(token); t != "true" {
		if !db.FindOne(token) {
			return c.Status(400).JSON(NewErrorResponse("Invalid API Key"))
		}

		redis.Set(token, "true")
	}

	return c.Next()
}
