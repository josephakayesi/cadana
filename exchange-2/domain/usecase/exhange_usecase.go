package usecase

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/josephakayesi/cadana/exchange-2/application/dto"
)

type ExchangeUsecase interface {
	GetRate(c *fiber.Ctx, r dto.GetExchangeRateDto) (*dto.GetExchangeRateResponseDto, error)
}

type exchangeUsecase struct {
	contextTimeout time.Duration
}

func NewExchangeUsecase(timeout time.Duration) ExchangeUsecase {
	return &exchangeUsecase{
		contextTimeout: timeout,
	}
}

func (uu *exchangeUsecase) GetRate(c *fiber.Ctx, r dto.GetExchangeRateDto) (*dto.GetExchangeRateResponseDto, error) {

	rand := rand.New(rand.NewSource(time.Now().UnixNano()))

	exchangeRates := map[string]float64{
		"USD-JPY": rand.Float64() * 100,
		"USD-EUR": rand.Float64() * 1.2,
		"USD-CAD": rand.Float64() * 1.35,
	}

	if rate, found := exchangeRates[r.CurrencyPair]; found {
		roundedRate := math.Round(rate*100) / 100

		exchangeRateResponseDto := dto.NewGetExchangeRateResponseDto(r.CurrencyPair, roundedRate)

		return exchangeRateResponseDto, nil
	}

	return nil, fmt.Errorf("currency pair %s unsupported", r.CurrencyPair)

}
