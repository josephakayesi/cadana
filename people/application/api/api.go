package api

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/josephakayesi/cadana/people/application/dto"
	"github.com/opensaucerer/goaxios"
)

func GetExhangeRatesForCurrency(currency string, wg *sync.WaitGroup, ch chan dto.ExchangeRate) {
	exchangeRateMap := map[string]string{
		"EUR": "USD-EUR",
		"JPY": "USD-JPY",
	}

	if currency == "USD" {
		ch <- dto.ExchangeRate{
			CurrencyPair: "USD",
			Rate:         1,
		}

		wg.Done()
		return
	}

	currencyPair := exchangeRateMap[currency]

	requestBody := dto.GetExchangeRateDto{
		CurrencyPair: currencyPair,
	}

	a := goaxios.GoAxios{
		Url:            "http://localhost:3000/api/v1/rates",
		Body:           requestBody,
		Method:         "POST",
		ResponseStruct: &dto.GetExchangeRateResponseDto{},
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	response := a.RunRest()
	if response.Error != nil {
		log.Fatalf("err: %v", response.Error)
	}

	rate, _ := response.Body.(*dto.GetExchangeRateResponseDto)

	b, err := json.MarshalIndent(rate, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(string(b))

	ch <- dto.ExchangeRate{
		CurrencyPair: currency,
		Rate:         (*rate)[currencyPair],
	}

	wg.Done()
}
