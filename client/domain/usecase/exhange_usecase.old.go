package usecase

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"strconv"
// 	"strings"
// 	"sync"
// 	"time"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/josephakayesi/cadana/exchange-2/application/dto"
// )

// type ExchangeUsecase interface {
// 	GetRate(c *fiber.Ctx, r dto.GetExchangeRateDto) (*dto.GetExchangeRateResponseDto, error)
// }

// type exchangeUsecase struct {
// 	contextTimeout time.Duration
// }

// func NewExchangeUsecase(timeout time.Duration) ExchangeUsecase {
// 	return &exchangeUsecase{
// 		contextTimeout: timeout,
// 	}
// }

// func (uu *exchangeUsecase) GetRate(c *fiber.Ctx, r dto.GetExchangeRateDto) (*dto.GetExchangeRateResponseDto, error) {
// 	var wg sync.WaitGroup
// 	ch := make(chan string, 2)

// 	url1 := "http://localhost:3001/api/v1/rates"
// 	// url2 := "http://loclhost:3002/api/v1/rates"

// 	wg.Add(1)
// 	go fetchRateFromExchange(url1, r, &wg, ch)
// 	// go fetchRateFromExchange(url2, r, &wg, ch)

// 	wg.Wait()

// 	close(ch)

// 	fmt.Printf("request body: %+v\n", r)

// 	if result, ok := <-ch; ok {
// 		fmt.Println("result: ", result)
// 		rate := result[strings.Index(result, ":")+1:]
// 		roundedRate, _ := strconv.ParseFloat(rate, 64)

// 		exchangeRateResponseDto := dto.NewGetExchangeRateResponseDto(r.CurrencyPair, roundedRate)
// 		fmt.Println(exchangeRateResponseDto)

// 		return nil, nil
// 	}

// 	return nil, fmt.Errorf("currency pair %s unsupported", r.CurrencyPair)

// }

// func fetchRateFromExchange(url string, requestBody dto.GetExchangeRateDto, wg *sync.WaitGroup, ch chan<- string) {
// 	defer wg.Done()

// 	requestBodyBytes, err := json.Marshal(requestBody)
// 	if err != nil {
// 		ch <- fmt.Sprintf("Error marshaling request body for %s: %s", url, err)
// 		return
// 	}

// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBodyBytes))
// 	if err != nil {
// 		ch <- fmt.Sprintf("Error creating request for %s: %s", url, err)
// 		return
// 	}

// 	req.Header.Set("Content-Type", "application/json")

// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		ch <- fmt.Sprintf("Error fetching %s: %s", url, err)
// 		return
// 	}

// 	defer resp.Body.Close()

// 	ch <- fmt.Sprintf("%s: %s", url, resp.Body)
// }
