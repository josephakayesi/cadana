package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/josephakayesi/cadana/client/application/dto"
	"github.com/josephakayesi/cadana/client/domain/usecase"
	"github.com/josephakayesi/cadana/client/internal"
	"golang.org/x/exp/slog"
)

// ExchangeController interface for ExchangeController methods
type ExchangeController interface {
	GetRate(c usecase.FiberContext) error
}

// exchangeController implements ExchangeController
type exchangeController struct {
	ExchangeUsecase usecase.ExchangeUsecase
	Logger          slog.Logger
}

// NewExchangeController creates a new instance of ExchangeController
func NewExchangeController(exchangeUsecase usecase.ExchangeUsecase, logger slog.Logger) ExchangeController {
	return &exchangeController{
		ExchangeUsecase: exchangeUsecase,
		Logger:          logger,
	}
}

// GetRate handles the HTTP request to get the exchange rate
func (ec *exchangeController) GetRate(c usecase.FiberContext) error {
	getExchangeRateDto := &dto.GetExchangeRateDto{}

	if err := c.BodyParser(getExchangeRateDto); err != nil {
		ec.Logger.Error("unable to parse GetExchangeRateDto", "error", err)
		return err
	}

	getExchangeRateResponseDto, errs := ec.ExchangeUsecase.GetRate(c, *getExchangeRateDto)

	if errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(internal.NewErrorResponse("unable to get exchange rate", errs))
	}

	return c.Status(fiber.StatusOK).JSON(getExchangeRateResponseDto)
}
