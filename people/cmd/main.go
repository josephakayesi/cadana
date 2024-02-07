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
	}

	fmt.Println("successfully opened people.json")
	defer file.Close()

	byteValue, _ := io.ReadAll(file)

	var people internal.People

	json.Unmarshal([]byte(byteValue), &people)

	// people.SortBySalaryInAscendingOrder()
	// people.SortBySalaryInDescendingOrder()
	// groupedPeople := people.GroupByCurrency()
	// internal.PrintGroupedPeople(groupedPeople)

	currencies := internal.GetUniqueCurrencies(people.People)

	var wg sync.WaitGroup

	successCh := make(chan dto.ExchangeRate, 2)
	errorCh := make(chan error, 2)

	currentExchangeRates := make(map[string]float64)

	for _, currency := range currencies {
		wg.Add(1)
		go api.GetExhangeRatesForCurrency(currency, &wg, successCh)
	}

	// wg.Wait()

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
	sortedPeople := people.SortBySalaryInAscendingOrder()
	// people.SortBySalaryInDescendingOrder()
	// groupedPeople := people.GroupByCurrency()
	internal.PrintPeople(sortedPeople)
}
