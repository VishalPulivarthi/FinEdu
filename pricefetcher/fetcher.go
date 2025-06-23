package pricefetcher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Fetcher struct {
	APIKey string
}

type PriceResponse struct {
	Name  string `json:"name"`
	Price string `json:"price"`
}

func NewFetcher(apiKey string) *Fetcher {
	return &Fetcher{APIKey: apiKey}
}

func (f *Fetcher) GetPrice(symbol string) (string, error) {
	url := fmt.Sprintf("https://finnhub.io/api/v1/quote?symbol=%s&token=%s", symbol, f.APIKey)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Sprintf("Failed to connect to API for %s: %v", symbol, err), err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("API returned status %d for %s", resp.StatusCode, symbol), fmt.Errorf("API error status: %d", resp.StatusCode)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var data map[string]float64
	if err := json.Unmarshal(body, &data); err != nil {
		return fmt.Sprintf("Failed to read response body for %s: %v", symbol, err), err
	}

	price := data["c"] // current price
	if price == 0 {
		return fmt.Sprintf("No valid price data found for %s", symbol), nil
	}
	return fmt.Sprintf("%.2f", price), nil
}

func (f *Fetcher) FetchPricesFromFile(filename string) ([]PriceResponse, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %v", filename, err)
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
			price = "Sorry, currently don't have the price of the symbol"
		}
		results = append(results, PriceResponse{Name: name, Price: price})
	}

	return results, nil
}
