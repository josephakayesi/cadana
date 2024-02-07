package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/josephakayesi/cadana/client/application/dto"
	"github.com/josephakayesi/cadana/client/domain/usecase"
	"github.com/josephakayesi/cadana/client/internal"
	"golang.org/x/exp/slog"
)

type ExchangeController struct {
	ExchangeUsecase usecase.ExchangeUsecase
	Logger          slog.Logger
}

func (ec *ExchangeController) GetRate(c *fiber.Ctx) error {

	getExchangeRateDto := &dto.GetExchangeRateDto{}

	if err := c.BodyParser(getExchangeRateDto); err != nil {
		ec.Logger.Error("unable to parse GetExchangeRateDto", "error", err)
		return err
	}

	getExchangeRateResponseDto, errs := ec.ExchangeUsecase.GetRate(c, *getExchangeRateDto)

	if errs != nil {
		return c.Status(400).JSON(internal.NewErrorResponse("unable to get exchange rate", errs))
	}

	return c.Status(fiber.StatusOK).JSON(getExchangeRateResponseDto)

}
