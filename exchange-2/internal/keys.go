package internal

import (
	"github.com/gofiber/fiber/v2"
	"github.com/josephakayesi/cadana/exchange-2/infra/config"
)

var redis = config.NewRedis()
var db = config.NewDatabase()

func ValidateAPIKey(c *fiber.Ctx) error {
	token := c.Get("x-access-token")

	if t, _ := redis.Get(token); t != "true" {
		if !db.FindOne(token) {
			return c.Status(400).JSON(NewErrorResponse("Invalid API Key"))
		}

		redis.Set(token, "true")
	}

	return c.Next()
}
