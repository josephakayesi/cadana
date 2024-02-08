package usecase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/josephakayesi/cadana/client/application/dto"
	"github.com/josephakayesi/cadana/client/infra/config"
)

// HTTPClient interface for making HTTP requests
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// FiberContext interface for Fiber.Ctx
type FiberContext interface {
	JSON(v interface{}) error
	Status(code int) *fiber.Ctx
	BodyParser(v interface{}) error
}

type ExchangeUsecase interface {
	GetRate(c FiberContext, r dto.GetExchangeRateDto) (*dto.GetExchangeRateResponseDto, []string)
	fetchRateFromExchange(url string, requestBody dto.GetExchangeRateDto, wg *sync.WaitGroup, successCh chan<- dto.GetExchangeRateResponseDto, errorCh chan<- error)
}

type exchangeUsecase struct {
	contextTimeout time.Duration
	httpClient     HTTPClient
}

func NewExchangeUsecase(timeout time.Duration, httpClient HTTPClient) ExchangeUsecase {
	return &exchangeUsecase{
		contextTimeout: timeout,
		httpClient:     httpClient,
	}
}

var cfg = config.NewConfig()

func (uu *exchangeUsecase) GetRate(c FiberContext, r dto.GetExchangeRateDto) (*dto.GetExchangeRateResponseDto, []string) {
	var wg sync.WaitGroup

	var errorsSlice []string

	successCh := make(chan dto.GetExchangeRateResponseDto, 2)
	errorCh := make(chan error, 2)

	urls := []string{
		cfg.EXCHANGE_SERVICE_URL_1,
		cfg.EXCHANGE_SERVICE_URL_2,
	}

	for _, url := range urls {
		wg.Add(1)
		go uu.fetchRateFromExchange(url, r, &wg, successCh, errorCh)
	}

	go func() {
		wg.Wait()
		close(successCh)
		close(errorCh)
	}()

	if getExchangeRateResponseDto, ok := <-successCh; ok {
		return &getExchangeRateResponseDto, nil
	}

	for err := range errorCh {
		errorsSlice = append(errorsSlice, err.Error())
	}

	return nil, errorsSlice
}

func (uu *exchangeUsecase) fetchRateFromExchange(url string, requestBody dto.GetExchangeRateDto, wg *sync.WaitGroup, successCh chan<- dto.GetExchangeRateResponseDto, errorCh chan<- error) {
	defer wg.Done()

	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		errorCh <- fmt.Errorf("error marshaling request body: %s", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		errorCh <- fmt.Errorf("error creating HTTP request: %s", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-access-token", cfg.API_TOKEN)

	resp, err := uu.httpClient.Do(req)
	if err != nil {
		errorCh <- fmt.Errorf("error making HTTP request: %s", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorCh <- fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		return
	}

	var rate dto.GetExchangeRateResponseDto
	err = json.NewDecoder(resp.Body).Decode(&rate)
	if err != nil {
		errorCh <- fmt.Errorf("error decoding response body: %s", err)
		return
	}

	successCh <- rate
}
