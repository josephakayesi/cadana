package usecase

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/josephakayesi/cadana/client/application/dto"
	"github.com/josephakayesi/cadana/client/infra/config"
	"github.com/stretchr/testify/assert"
)

func TestGetRate(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := dto.GetExchangeRateResponseDto{
			"USD-EUR": 1.2,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	mockConfig := config.Config{
		EXCHANGE_SERVICE_URL_1: mockServer.URL,
		EXCHANGE_SERVICE_URL_2: "http://invalid-url",
		API_TOKEN:              "your-api-token",
	}

	usecase := NewExchangeUsecase(5 * time.Second)
	exchangeUsecase := usecase.(*exchangeUsecase)
	exchangeUsecase.contextTimeout = 5 * time.Second
	cfg = &mockConfig

	request := dto.GetExchangeRateDto{
		CurrencyPair: "USD-EUR",
	}

	ctx := &fiber.Ctx{}

	response, errors := exchangeUsecase.GetRate(ctx, request)

	assert.NotNil(t, response)

	exchangeRate, ok := (*response)["USD-EUR"]
	if !ok {
		assert.True(t, ok)
		assert.Equal(t, 1.2, exchangeRate)
	}

	assert.Empty(t, errors)
}

func TestFetchRateFromExchange(t *testing.T) {
	// Prepare a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "your-api-token", r.Header.Get("x-access-token"))

		// Mock the response based on your requirements
		response := dto.GetExchangeRateResponseDto{"yourKey": 1.2}

		// Encode the response as JSON and write it to the response writer
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))

	defer server.Close()

	// Set up the necessary channels and wait group
	var wg sync.WaitGroup
	successCh := make(chan dto.GetExchangeRateResponseDto, 1)
	errorCh := make(chan error, 1)

	// Call the function with the mock server URL
	wg.Add(1)
	go fetchRateFromExchange(server.URL, dto.GetExchangeRateDto{}, &wg, successCh, errorCh)

	// Wait for the function to complete
	wg.Wait()

	// Verify that the success channel received the expected response
	select {
	case exchangeRate := <-successCh:
		assert.Equal(t, 1.2, exchangeRate["yourKey"])
	case err := <-errorCh:
		t.Fatalf("Unexpected error: %s", err)
	}
}
