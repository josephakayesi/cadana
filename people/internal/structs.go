package internal

import (
	"encoding/json"
	"fmt"
	"sort"
)

type Person struct {
	ID         string `json:"id"`
	PersonName string `json:"personName"`
	Salary     Salary `json:"salary"`
}

type Salary struct {
	Value    float64 `json:"value"`
	Currency string  `json:"currency"`
}

type People struct {
	People []Person
}

type ExchangeRateResponse struct {
	ExchangeRate map[string]float64
}

func (p *People) SortBySalaryInAscendingOrder() []Person {
	uniqueCurrencies := GetUniqueCurrencies(p.People)

	var sortedPeople []Person

	for _, currency := range uniqueCurrencies {
		currencyPeople := filterPeopleByCurrency(p.People, currency)

		sort.SliceStable(currencyPeople, func(i, j int) bool {
			return currencyPeople[i].Salary.Value < currencyPeople[j].Salary.Value
		})

		// fmt.Printf("\nSorted by salary value in ascending order for currency %s:\n", currency)
		sortedPeople = append(sortedPeople, currencyPeople...)

	}

	// printPeople(sortedPeople)
	return sortedPeople
}

func (p *People) SortBySalaryInDescendingOrder() []Person {
	uniqueCurrencies := GetUniqueCurrencies(p.People)

	var sortedPeople []Person

	for _, currency := range uniqueCurrencies {
		currencyPeople := filterPeopleByCurrency(p.People, currency)

		sort.SliceStable(currencyPeople, func(i, j int) bool {
			return currencyPeople[i].Salary.Value > currencyPeople[j].Salary.Value
		})

		// fmt.Printf("\nSorted by salary value in ascending order for currency %s:\n", currency)
		sortedPeople = append(sortedPeople, currencyPeople...)
	}

	// printPeople(sortedPeople)
	return sortedPeople
}

func (p *People) GroupByCurrency() map[string][]Person {
	groupedPeople := make(map[string][]Person)

	for _, person := range p.People {
		currency := person.Salary.Currency
		groupedPeople[currency] = append(groupedPeople[currency], person)
	}

	return groupedPeople
}

func PrintGroupedPeople(groupedPeople map[string][]Person) {
	for currency, people := range groupedPeople {
		fmt.Printf("Currency: %s\n", currency)
		PrintPeople(people)
		fmt.Println()
	}
}

func GetUniqueCurrencies(people []Person) []string {
	uniqueCurrencies := make(map[string]bool)
	var currencies []string

	for _, person := range people {
		if !uniqueCurrencies[person.Salary.Currency] {
			uniqueCurrencies[person.Salary.Currency] = true
			currencies = append(currencies, person.Salary.Currency)
		}
	}

	return currencies
}

func filterPeopleByCurrency(people []Person, currency string) []Person {
	var filteredPeople []Person

	for _, person := range people {
		if person.Salary.Currency == currency {
			filteredPeople = append(filteredPeople, person)
		}
	}

	return filteredPeople
}

func PrintPeople(people []Person) {
	jsonData, err := json.MarshalIndent(people, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(jsonData))
}
