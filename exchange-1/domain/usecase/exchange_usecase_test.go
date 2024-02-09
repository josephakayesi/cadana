package usecase

import (
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/josephakayesi/cadana/exchange-1/application/dto"
	"github.com/stretchr/testify/assert"
)

func TestGetRate(t *testing.T) {
	mockCurrencyPair := "USD-EUR"
	mockRequest := dto.GetExchangeRateDto{
		CurrencyPair: mockCurrencyPair,
	}

	uu := NewExchangeUsecase(time.Second)

	c := &fiber.Ctx{}

	response, err := uu.GetRate(c, mockRequest)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, mockCurrencyPair, response.CurrencyPair)
	assert.Greater(t, response.ExchangeRate, 0.0)
}

func TestGetRateUnsupportedCurrencyPair(t *testing.T) {
	mockCurrencyPair := "USD-GBP"
	mockRequest := dto.GetExchangeRateDto{
		CurrencyPair: mockCurrencyPair,
	}

	uu := NewExchangeUsecase(time.Second)

	c := &fiber.Ctx{}

	response, err := uu.GetRate(c, mockRequest)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "unsupported")
}
