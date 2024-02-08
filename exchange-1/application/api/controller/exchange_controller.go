package controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/josephakayesi/cadana/exchange-1/application/dto"
	"github.com/josephakayesi/cadana/exchange-1/domain/usecase"
	"github.com/josephakayesi/cadana/exchange-1/internal"
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

	getExchangeRateResponseDto, err := ec.ExchangeUsecase.GetRate(c, *getExchangeRateDto)

	if err != nil {
		fmt.Println("Error", err.Error())

		return c.Status(400).JSON(internal.NewErrorResponse(err.Error()))
	}

	rateMap := make(map[string]float64)
	rateMap[getExchangeRateResponseDto.CurrencyPair] = getExchangeRateResponseDto.ExchangeRate

	return c.Status(fiber.StatusOK).JSON(rateMap)
}
