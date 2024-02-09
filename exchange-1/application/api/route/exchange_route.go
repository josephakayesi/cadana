package route

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/josephakayesi/cadana/exchange-1/application/api/controller"
	"github.com/josephakayesi/cadana/exchange-1/domain/usecase"
	"golang.org/x/exp/slog"
)

func NewExchangeRouter(timeout time.Duration, group fiber.Router) {
	ec := &controller.ExchangeController{
		ExchangeUsecase: usecase.NewExchangeUsecase(timeout),
		Logger:          *slog.New(slog.NewJSONHandler(os.Stdout, nil)).With("component", "exhange-1"),
	}

	group.Post("/rates", ec.GetRate)
}
