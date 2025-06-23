package pricefetcher

import (
	"os"
	"strings"
	"testing"
)

type TestFetcher struct{}

func (f *TestFetcher) GetPrice(symbol string) (string, error) {
	switch symbol {
	case "AAPL":
		return "180.00", nil
	case "MSFT":
		return "350.50", nil
	case "GOOG":
		return "1500.25", nil
	default:
		return "Error", nil
	}
}

func (f *TestFetcher) FetchPricesFromFile(filename string) ([]PriceResponse, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	var results []PriceResponse

	for _, line := range lines {
		name := strings.TrimSpace(line)
		if name == "" {
			continue
		}
		price, err := f.GetPrice(name)
		if err != nil {
			price = "Error"
		}
		results = append(results, PriceResponse{Name: name, Price: price})
	}

	return results, nil
}

func TestFetchPricesFromFile(t *testing.T) {
	testSymbols := "AAPL\nMSFT\nGOOG\nAMZN\nUNKNOWN_SYMBOL\n"

	tmpFile, err := os.CreateTemp("", "symbols.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(testSymbols)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	fetcher := &TestFetcher{}
	results, err := fetcher.FetchPricesFromFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedNumResults := 5
	if len(results) != expectedNumResults {
		t.Errorf("Expected %d results, got %d", expectedNumResults, len(results))
	}

	expected := []PriceResponse{
		{Name: "AAPL", Price: "180.00"},
		{Name: "MSFT", Price: "350.50"},
		{Name: "GOOG", Price: "1500.25"},
		{Name: "AMZN", Price: "Error"},
		{Name: "UNKNOWN_SYMBOL", Price: "Error"},
	}

	for i, exp := range expected {
		if i >= len(results) {
			t.Errorf("Result %d missing. Expected %+v", i, exp)
			continue
		}
		actual := results[i]
		if actual.Name != exp.Name || actual.Price != exp.Price {
			t.Errorf("Result %d: Expected %+v, got %+v", i, exp, actual)
		}
	}
}
