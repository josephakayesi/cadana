package dto

type GetExchangeRateDto struct {
	CurrencyPair string `json:"currency_pair" validate:"required,ascii,min=2,max=30"`
}

type GetExchangeRateResponseDto struct {
	CurrencyPair string
	ExchangeRate float64
}

func NewGetExchangeRateResponseDto(currencyPair string, exchangeRate float64) *GetExchangeRateResponseDto {
	return &GetExchangeRateResponseDto{
		CurrencyPair: currencyPair,
		ExchangeRate: exchangeRate,
	}
}
