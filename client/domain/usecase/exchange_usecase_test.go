package usecase

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/gofiber/fiber/v2"
	"github.com/josephakayesi/cadana/client/application/dto"
)

// MockHTTPClient is a mock implementation of HTTPClient interface
type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

// MockFiberContext is a mock implementation of FiberContext interface
type MockFiberContext struct {
	mock.Mock
}

func (m *MockFiberContext) JSON(v interface{}) error {
	args := m.Called(v)
	return args.Error(0)
}

func (m *MockFiberContext) Status(code int) *fiber.Ctx {
	return &fiber.Ctx{}
}

func TestGetRate(t *testing.T) {
	// Setup mock HTTP client
	mockHTTPClient := new(MockHTTPClient)

	// Setup mock Fiber context
	mockFiberContext := new(MockFiberContext)

	// Create an instance of the exchange usecase with mock dependencies
	uu := NewExchangeUsecase(time.Second, mockHTTPClient)

	// Define the expected request body
	requestBody := dto.GetExchangeRateDto{
		CurrencyPair: "USD-EUR",
	}

	// Define a sample response
	mockResponse := &dto.GetExchangeRateResponseDto{
		"USD-EUR": 1.11,
	}

	// Marshal the mock response to JSON
	jsonResponse, _ := json.Marshal(mockResponse)

	// Mock the HTTP client response
	mockHTTPClient.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBuffer(jsonResponse)),
	}, nil)

	// Mock the Fiber context JSON method
	mockFiberContext.On("JSON", mockResponse).Return(nil)

	// Call the GetRate method
	result, err := uu.GetRate(mockFiberContext, requestBody)

	// Assert the result and error
	assert.NotNil(t, result)
	if len(err) > 0 {
		assert.Fail(t, "Unexpected errors:", strings.Join(err, ", "))
	}

	// Verify that the expected methods were called
	mockHTTPClient.AssertExpectations(t)
}

func TestFetchRateFromExchange(t *testing.T) {
	// Setup mock HTTP client
	mockHTTPClient := new(MockHTTPClient)

	// Create an instance of the exchange usecase with mock dependencies
	uu := NewExchangeUsecase(time.Second, mockHTTPClient)

	// Define the expected request body
	requestBody := dto.GetExchangeRateDto{
		// Your request body fields here
	}

	// Define a sample response
	mockResponse := &dto.GetExchangeRateResponseDto{
		// Your response fields here
	}

	// Marshal the mock response to JSON
	jsonResponse, _ := json.Marshal(mockResponse)

	// Mock the HTTP client response
	mockHTTPClient.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBuffer(jsonResponse)),
	}, nil)

	// Create channels for success and error
	successCh := make(chan dto.GetExchangeRateResponseDto, 1)
	errorCh := make(chan error, 1)

	// Create a WaitGroup
	var wg sync.WaitGroup

	// Increment the WaitGroup counter
	wg.Add(1)

	// Call the fetchRateFromExchange method
	go uu.fetchRateFromExchange("mockURL", requestBody, &wg, successCh, errorCh)

	// Decrement the WaitGroup counter when the goroutine completes
	go func() {
		wg.Wait()
		close(successCh)
		close(errorCh)
	}()

	// Assert the result and error
	select {
	case result := <-successCh:
		assert.NotNil(t, result)
	case err := <-errorCh:
		assert.Fail(t, "Unexpected error: "+err.Error())
	}

	// Verify that the expected methods were called
	mockHTTPClient.AssertExpectations(t)
}
