package internal

import (
	"bytes"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestSortBySalaryInAscendingOrderForSameCurrency(t *testing.T) {
	people := People{
		People: []Person{
			{ID: "1", PersonName: "John", Salary: Salary{Value: 5000, Currency: "USD"}},
			{ID: "2", PersonName: "Jane", Salary: Salary{Value: 4000, Currency: "USD"}},
			{ID: "3", PersonName: "Bob", Salary: Salary{Value: 6000, Currency: "USD"}},
		},
	}

	sortedPeople := people.SortBySalaryInAscendingOrder()

	expectedOrder := []float64{4000, 5000, 6000}

	for i, person := range sortedPeople {
		if person.Salary.Value != expectedOrder[i] {
			t.Errorf("Expected order: %v, Actual order: %v", expectedOrder, sortedPeople)
			break
		}
	}
}

func TestSortBySalaryInAscendingOrderForMultipleCurrencies(t *testing.T) {
	people := People{
		People: []Person{
			{ID: "1", PersonName: "John", Salary: Salary{Value: 5000, Currency: "USD"}},
			{ID: "2", PersonName: "Jane", Salary: Salary{Value: 4000, Currency: "EUR"}},
			{ID: "3", PersonName: "Bob", Salary: Salary{Value: 6000, Currency: "USD"}},
		},
	}

	sortedPeople := people.SortBySalaryInAscendingOrder()

	expectedOrder := []float64{5000, 6000, 4000}

	for i, person := range sortedPeople {
		if person.Salary.Value != expectedOrder[i] {
			t.Errorf("Expected order: %v, Actual order: %v", expectedOrder, sortedPeople)
			break
		}
	}
}

func TestSortBySalaryInDescendingOrderForSameCurrency(t *testing.T) {
	people := People{
		People: []Person{
			{ID: "1", PersonName: "John", Salary: Salary{Value: 5000, Currency: "USD"}},
			{ID: "2", PersonName: "Jane", Salary: Salary{Value: 4000, Currency: "USD"}},
			{ID: "3", PersonName: "Bob", Salary: Salary{Value: 6000, Currency: "USD"}},
		},
	}

	sortedPeople := people.SortBySalaryInDescendingOrder()

	expectedOrder := []float64{6000, 5000, 4000}

	for i, person := range sortedPeople {
		if person.Salary.Value != expectedOrder[i] {
			t.Errorf("Expected order: %v, Actual order: %v", expectedOrder, sortedPeople)
			break
		}
	}
}

func TestSortBySalaryInDescendingOrderForMultipleCurrencies(t *testing.T) {
	people := People{
		People: []Person{
			{ID: "1", PersonName: "John", Salary: Salary{Value: 5000, Currency: "USD"}},
			{ID: "2", PersonName: "Jane", Salary: Salary{Value: 4000, Currency: "EUR"}},
			{ID: "3", PersonName: "Bob", Salary: Salary{Value: 6000, Currency: "USD"}},
		},
	}

	sortedPeople := people.SortBySalaryInDescendingOrder()

	expectedOrder := []float64{6000, 5000, 4000}

	for i, person := range sortedPeople {
		if person.Salary.Value != expectedOrder[i] {
			t.Errorf("Expected order: %v, Actual order: %v", expectedOrder, sortedPeople)
			break
		}
	}
}

func TestGroupByCurrency(t *testing.T) {
	people := People{
		People: []Person{
			{ID: "1", PersonName: "John", Salary: Salary{Value: 5000, Currency: "USD"}},
			{ID: "2", PersonName: "Jane", Salary: Salary{Value: 4000, Currency: "EUR"}},
			{ID: "3", PersonName: "Bob", Salary: Salary{Value: 6000, Currency: "USD"}},
		},
	}

	groupedPeople := people.GroupByCurrency()

	expectedGrouping := map[string][]Person{
		"USD": {
			{ID: "1", PersonName: "John", Salary: Salary{Value: 5000, Currency: "USD"}},
			{ID: "3", PersonName: "Bob", Salary: Salary{Value: 6000, Currency: "USD"}},
		},
		"EUR": {
			{ID: "2", PersonName: "Jane", Salary: Salary{Value: 4000, Currency: "EUR"}},
		},
	}

	if !reflect.DeepEqual(groupedPeople, expectedGrouping) {
		t.Errorf("Expected grouping: %v, Actual grouping: %v", expectedGrouping, groupedPeople)
	}
}

func TestGetUniqueCurrencies(t *testing.T) {
	people := []Person{
		{ID: "1", PersonName: "John", Salary: Salary{Value: 5000, Currency: "USD"}},
		{ID: "2", PersonName: "Jane", Salary: Salary{Value: 4000, Currency: "EUR"}},
		{ID: "3", PersonName: "Bob", Salary: Salary{Value: 6000, Currency: "USD"}},
		{ID: "4", PersonName: "Alice", Salary: Salary{Value: 3000, Currency: "USD"}},
		{ID: "5", PersonName: "Smith", Salary: Salary{Value: 7000, Currency: "EUR"}},
		{ID: "6", PersonName: "Emma", Salary: Salary{Value: 5500, Currency: "EUR"}},
	}

	uniqueCurrencies := GetUniqueCurrencies(people)

	expectedCurrencies := []string{"USD", "EUR"}

	if len(uniqueCurrencies) != len(expectedCurrencies) {
		t.Errorf("Expected %d unique currencies, but got %d", len(expectedCurrencies), len(uniqueCurrencies))
		return
	}

	for _, currency := range expectedCurrencies {
		found := false
		for _, unique := range uniqueCurrencies {
			if currency == unique {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected currency %s not found in unique currencies", currency)
		}
	}
}

func TestFilterPeopleByCurrency(t *testing.T) {
	people := []Person{
		{ID: "1", PersonName: "John", Salary: Salary{Value: 5000, Currency: "USD"}},
		{ID: "2", PersonName: "Jane", Salary: Salary{Value: 4000, Currency: "EUR"}},
		{ID: "3", PersonName: "Bob", Salary: Salary{Value: 6000, Currency: "USD"}},
	}

	filteredPeopleUSD := filterPeopleByCurrency(people, "USD")

	expectedPeopleUSD := []Person{
		{ID: "1", PersonName: "John", Salary: Salary{Value: 5000, Currency: "USD"}},
		{ID: "3", PersonName: "Bob", Salary: Salary{Value: 6000, Currency: "USD"}},
	}

	if !reflect.DeepEqual(filteredPeopleUSD, expectedPeopleUSD) {
		t.Errorf("Expected filtered people for USD: %v, Actual filtered people: %v", expectedPeopleUSD, filteredPeopleUSD)
	}

	filteredPeopleEUR := filterPeopleByCurrency(people, "EUR")
	expectedPeopleEUR := []Person{
		{ID: "2", PersonName: "Jane", Salary: Salary{Value: 4000, Currency: "EUR"}},
	}

	if !reflect.DeepEqual(filteredPeopleEUR, expectedPeopleEUR) {
		t.Errorf("Expected filtered people for EUR: %v, Actual filtered people: %v", expectedPeopleEUR, filteredPeopleEUR)
	}

	// Test with empty input
	emptyFilteredPeople := filterPeopleByCurrency([]Person{}, "USD")
	emptyExpectedPeople := []Person{}

	if len(emptyFilteredPeople) != len(emptyExpectedPeople) {
		t.Errorf("Expected empty filtered people: %v, Actual empty filtered people: %v", emptyExpectedPeople, emptyFilteredPeople)
	}

	// Test with no matches
	noMatchesFilteredPeople := filterPeopleByCurrency(people, "GBP")
	noMatchesExpectedPeople := []Person{}

	if len(noMatchesFilteredPeople) != len(noMatchesExpectedPeople) {
		t.Errorf("Expected no matches filtered people: %v, Actual no matches filtered people: %v", noMatchesExpectedPeople, noMatchesFilteredPeople)
	}
}

func TestPrintPeople(t *testing.T) {
	people := []Person{
		{ID: "1", PersonName: "John", Salary: Salary{Value: 5000, Currency: "USD"}},
		{ID: "2", PersonName: "Jane", Salary: Salary{Value: 4000, Currency: "EUR"}},
		{ID: "3", PersonName: "Bob", Salary: Salary{Value: 6000, Currency: "USD"}},
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	PrintPeople(people)

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)

	expectedOutput := "[\n  {\n    \"id\": \"1\",\n    \"personName\": \"John\",\n    \"salary\": {\n      \"value\": 5000,\n      \"currency\": \"USD\"\n    }\n  },\n  {\n    \"id\": \"2\",\n    \"personName\": \"Jane\",\n    \"salary\": {\n      \"value\": 4000,\n      \"currency\": \"EUR\"\n    }\n  },\n  {\n    \"id\": \"3\",\n    \"personName\": \"Bob\",\n    \"salary\": {\n      \"value\": 6000,\n      \"currency\": \"USD\"\n    }\n  }\n]"

	if strings.TrimSpace(buf.String()) != strings.TrimSpace(expectedOutput) {
		t.Errorf("Expected output:\n%v\nActual output:\n%v", expectedOutput, buf.String())
	}
}
