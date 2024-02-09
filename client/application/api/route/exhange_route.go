package route

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/josephakayesi/cadana/client/application/api/controller"
	"github.com/josephakayesi/cadana/client/domain/usecase"
)

func NewExchangeRouter(timeout time.Duration, group fiber.Router) {
	ec := &controller.ExchangeController{
		ExchangeUsecase: usecase.NewExchangeUsecase(timeout),
		// Logger:          *slog.New(slog.NewJSONHandler(os.Stdout, nil)).With("component", "client"),
	}

	group.Post("/rates", ec.GetRate)
}
