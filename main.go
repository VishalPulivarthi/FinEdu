package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const apiKey = "d10pu2hr01qse6ldaqogd10pu2hr01qse6ldaqp0" // Replace with your Finnhub API key

// Struct to map the response from Finnhub
type FinnhubResponse struct {
	CurrentPrice float64 `json:"c"`
}

// Function to fetch stock price
func getStockPrice(symbol string) (float64, error) {
	url := fmt.Sprintf("https://finnhub.io/api/v1/quote?symbol=%s&token=%s", symbol, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var data FinnhubResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return 0, err
	}

	return data.CurrentPrice, nil
}

// API handler for /price route
func handler(w http.ResponseWriter, r *http.Request) {
	symbol := r.URL.Query().Get("symbol")
	if symbol == "" {
		http.Error(w, "Missing stock symbol!", http.StatusBadRequest)
		return
	}

	price, err := getStockPrice(symbol)
	if err != nil {
		http.Error(w, "Error getting stock price", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"symbol": symbol,
		"price":  price,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/price", handler)
	fmt.Println("✅ API running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
