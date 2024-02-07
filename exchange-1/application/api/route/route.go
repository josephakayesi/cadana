package route

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func Setup(t time.Duration, engine *fiber.App) {
	v1 := engine.Group("/api/v1")

	NewExchangeRouter(t, v1)
}
