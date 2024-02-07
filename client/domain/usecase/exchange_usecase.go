package usecase

import (
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/josephakayesi/cadana/exchange-2/application/dto"
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

func (uu *exchangeUsecase) GetRate(c *fiber.Ctx, r dto.GetExchangeRateDto) (*dto.GetExchangeRateResponseDto, []string) {
	var wg sync.WaitGroup

	var errorsSlice []string

	successCh := make(chan dto.GetExchangeRateResponseDto, 2)
	errorCh := make(chan error, 2)

	urls := []string{
		"http://localhost:3001/api/v1/rates",
		"http://localhost:3002/api/v1/rates",
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
		// BearerToken:    token,
	}
	response := a.RunRest()

	if response.Error != nil {
		errorCh <- fmt.Errorf("error fetching %s: %s", url, response.Error)
		return
	}

	data, _ := response.Body.(*dto.GetExchangeRateResponseDto)

	successCh <- *data
}
