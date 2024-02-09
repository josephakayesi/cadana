package api

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/josephakayesi/cadana/people/application/dto"
)

// ExchangeRateGetter is an interface for getting exchange rates.
type ExchangeRateGetter interface {
	Get(url string, contentType string, body io.Reader) (*http.Response, error)
}

// Logger is an interface for logging.
type Logger interface {
	Fatalf(format string, v ...interface{})
}

// DefaultExchangeRateGetter is the default implementation of ExchangeRateGetter using http.Client.
type DefaultExchangeRateGetter struct{}

// Get implements the ExchangeRateGetter interface using http.Client.
func (d *DefaultExchangeRateGetter) Get(url string, contentType string, body io.Reader) (*http.Response, error) {
	return http.Post(url, contentType, body)
}

// DefaultLogger is the default implementation of Logger using the standard log package.
type DefaultLogger struct{}

// Fatalf implements the Logger interface using the standard log package.
func (d *DefaultLogger) Fatalf(format string, v ...interface{}) {
	log.Fatalf(format, v...)
}

// GetExchangeRatesForCurrency retrieves exchange rates for a given currency.
func GetExchangeRatesForCurrency(currency string, url string, wg *sync.WaitGroup, ch chan dto.ExchangeRate, getter ExchangeRateGetter, logger Logger) {
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

	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		logger.Fatalf("Error marshaling request body: %v", err)
	}

	resp, err := getter.Get(url, "application/json", bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		logger.Fatalf("Error making HTTP request: %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Fatalf("Error reading response body: %v", err)
	}

	var rate dto.GetExchangeRateResponseDto
	err = json.Unmarshal(body, &rate)
	if err != nil {
		logger.Fatalf("Error unmarshaling response body: %v", err)
	}

	ch <- dto.ExchangeRate{
		CurrencyPair: currency,
		Rate:         rate[currencyPair],
	}

	wg.Done()
}
