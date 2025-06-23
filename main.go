package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"project/pricefetcher"
)

func main() {
	apiKey := os.Getenv("FINNHUB_API_KEY")
	if apiKey == "" {
		fmt.Println("Missing FINNHUB_API_KEY environment variable")
		return
	}

	f := pricefetcher.NewFetcher(apiKey)

	http.HandleFunc("/prices", func(w http.ResponseWriter, r *http.Request) {
		prices, err := f.FetchPricesFromFile("items.txt")
		if err != nil {
			http.Error(w, "Error reading items file", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(prices)
	})

	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
