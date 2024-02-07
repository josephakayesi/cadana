package usecase

import (
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/josephakayesi/cadana/client/application/dto"
	"github.com/josephakayesi/cadana/client/infra/config"
	"github.com/opensaucerer/goaxios"
)

type ExchangeUsecase interface {
	GetRate(c *fiber.Ctx, r dto.GetExchangeRateDto) (*dto.GetExchangeRateResponseDto, []string)
}

type exchangeUsecase struct {
	contextTimeout time.Duration
}

func NewExchangeUsecase(timeout time.Duration) ExchangeUsecase {
	return &exchangeUsecase{
		contextTimeout: timeout,
	}
}

var cfg = config.NewConfig()

func (uu *exchangeUsecase) GetRate(c *fiber.Ctx, r dto.GetExchangeRateDto) (*dto.GetExchangeRateResponseDto, []string) {
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
		go fetchRateFromExchange(url, r, &wg, successCh, errorCh)
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

func fetchRateFromExchange(url string, requestBody dto.GetExchangeRateDto, wg *sync.WaitGroup, successCh chan<- dto.GetExchangeRateResponseDto, errorCh chan<- error) {
	defer wg.Done()

	a := goaxios.GoAxios{
		Url:            url,
		Body:           requestBody,
		Method:         "POST",
		ResponseStruct: &dto.GetExchangeRateResponseDto{},
		Headers: map[string]string{
			"Content-Type":   "application/json",
			"x-access-token": cfg.API_TOKEN,
		},
	}
	response := a.RunRest()

	if response.Error != nil {
		errorCh <- fmt.Errorf("error fetching %s: %s", url, response.Error)
		return
	}

	data, _ := response.Body.(*dto.GetExchangeRateResponseDto)

	successCh <- *data
}
