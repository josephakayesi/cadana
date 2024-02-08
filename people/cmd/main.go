package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"sync"

	"github.com/josephakayesi/cadana/people/application/api"
	"github.com/josephakayesi/cadana/people/application/dto"
	"github.com/josephakayesi/cadana/people/internal"
)

func main() {
	file, err := os.Open("data/people.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	var people internal.People
	err = json.Unmarshal(byteValue, &people)
	if err != nil {
		fmt.Println(err)
		return
	}

	currencies := internal.GetUniqueCurrencies(people.People)

	var wg sync.WaitGroup

	successCh := make(chan dto.ExchangeRate, len(currencies))
	errorCh := make(chan error, len(currencies))

	currentExchangeRates := make(map[string]float64)

	// Provide the actual URL for the exchange rates API
	url := "http://localhost:3000/api/v1/rates"

	// Create instances of the DefaultExchangeRateGetter and DefaultLogger
	exchangeRateGetter := &api.DefaultExchangeRateGetter{}
	logger := &api.DefaultLogger{}

	for _, currency := range currencies {
		wg.Add(1)
		go api.GetExchangeRatesForCurrency(currency, url, &wg, successCh, exchangeRateGetter, logger)
	}

	go func() {
		wg.Wait()
		close(successCh)
		close(errorCh)
	}()

	for i := 0; i < len(currencies); i++ {
		select {
		case res := <-successCh:
			currentExchangeRates[res.CurrencyPair] = res.Rate
		case err := <-errorCh:
			fmt.Println(err)
		}
	}

	for i, person := range people.People {
		people.People[i].Salary.Value = math.Round(people.People[i].Salary.Value * currentExchangeRates[person.Salary.Currency])
	}

	// internal.PrintPeople(people.People)
	// people.SortBySalaryInDescendingOrder()
	// groupedPeople := people.GroupByCurrency()

	sortedPeople := people.SortBySalaryInAscendingOrder()
	internal.PrintPeople(sortedPeople)
}
