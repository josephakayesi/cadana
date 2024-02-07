package dto

type GetExchangeRateDto struct {
	CurrencyPair string `json:"currency_pair" validate:"required,ascii,min=2,max=30"`
}

type GetExchangeRateResponseDto map[string]float64

type ExchangeRate struct {
	CurrencyPair string
	Rate         float64
}
