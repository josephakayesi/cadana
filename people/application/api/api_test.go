package api

import (
	"bytes"
	"io"
	"net/http"
	"sync"
	"testing"

	"github.com/josephakayesi/cadana/people/application/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockExchangeRateGetter is a mock implementation of the ExchangeRateGetter interface.
type MockExchangeRateGetter struct {
	mock.Mock
}

// Get is the mock implementation for ExchangeRateGetter interface.
func (m *MockExchangeRateGetter) Get(url string, contentType string, body io.Reader) (*http.Response, error) {
	args := m.Called(url, contentType, body)
	return args.Get(0).(*http.Response), args.Error(1)
}

// MockLogger is a mock implementation of the Logger interface.
type MockLogger struct {
	mock.Mock
}

// Fatalf is the mock implementation for Logger interface.
func (m *MockLogger) Fatalf(format string, v ...interface{}) {
	m.Called(format, v)
}

func TestGetExchangeRatesForCurrency(t *testing.T) {
	// Set up mock objects
	mockGetter := new(MockExchangeRateGetter)
	mockLogger := new(MockLogger)

	// Set up wait group and channel for testing
	var wg sync.WaitGroup
	ch := make(chan dto.ExchangeRate, 1)

	// Test when currency is "USD"
	wg.Add(1)
	go GetExchangeRatesForCurrency("USD", "mock-url", &wg, ch, mockGetter, mockLogger)
	wg.Wait()

	// Assert that the exchange rate for USD is sent to the channel
	expectedRate := dto.ExchangeRate{CurrencyPair: "USD", Rate: 1}
	actualRate := <-ch
	assert.Equal(t, expectedRate, actualRate)

	// Test when currency is "JPY"
	wg.Add(1)

	// Mock the HTTP response for the JPY currency
	mockResponse := &http.Response{
		Body:       io.NopCloser(bytes.NewReader([]byte(`{"USD-JPY": 110.5}`))),
		StatusCode: http.StatusOK,
	}
	mockGetter.On("Get", "mock-url", "application/json", mock.Anything).Return(mockResponse, nil)

	go GetExchangeRatesForCurrency("JPY", "mock-url", &wg, ch, mockGetter, mockLogger)
	wg.Wait()

	// Assert that the exchange rate for JPY is sent to the channel
	expectedRateJPY := dto.ExchangeRate{CurrencyPair: "JPY", Rate: 110.5}
	actualRateJPY := <-ch
	assert.Equal(t, expectedRateJPY, actualRateJPY)

	mockLogger.AssertNotCalled(t, "Fatalf")
}
