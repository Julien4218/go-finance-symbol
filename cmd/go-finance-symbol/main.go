package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"time"

	"github.com/Julien4218/go-finance-symbol/observability"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	if err := Command.Execute(); err != nil {
		if err != flag.ErrHelp {
			log.Fatal(err)
		}
	}
}

func init() {
}

func globalInit(cmd *cobra.Command, args []string) {
	observability.Init()
}

var Command = &cobra.Command{
	Use:              "go-finance-symbol",
	Short:            "Test App",
	PersistentPreRun: globalInit,
	Long:             `Execute`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Execute with a symbol as a command argument, for example AAPL for apple")
			return
		}
		for _, symbol := range args {
			price, err := getToday(symbol)
			if err != nil {
				log.Error(err)
			} else {
				log.Infof("Stock price for %s: $%.2f\n", symbol, price)
				observability.GetOrCreateGauge(fmt.Sprintf("%s_1d", symbol)).Set(price)
			}

			average, err := getPrevious(symbol, SixMonth)
			if err != nil {
				log.Error(err)
			} else {
				log.Infof("6-month average price for %s: $%.2f\n", symbol, average)
				observability.GetOrCreateGauge(fmt.Sprintf("%s_6m", symbol)).Set(average)
			}

			observability.HarvestNow()
		}
	},
}

// Define a struct to hold the JSON response
type YahooFinanceResponse struct {
	Chart struct {
		Result []struct {
			Meta struct {
				RegularMarketPrice float64 `json:"regularMarketPrice"`
			} `json:"meta"`
			Timestamp  []int64 `json:"timestamp"`
			Indicators struct {
				Quote []struct {
					Close []float64 `json:"close"`
				} `json:"quote"`
			} `json:"indicators"`
		} `json:"result"`
	} `json:"chart"`
}

func getToday(symbol string) (float64, error) {
	for {
		url := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s", symbol)
		request := getClient(symbol)
		resp, err := request.Get(url)
		if err != nil {
			return 0, fmt.Errorf("Error fetching stock price: %v", err)
		}

		// Check for "Too Many Requests" response
		if resp.StatusCode() == 429 {
			log.Info("Rate limit exceeded. Retrying after a delay...\n")
			time.Sleep(10 * time.Second) // Wait for 10 seconds before retrying
			continue
		}

		// Parse the JSON response
		var data YahooFinanceResponse
		err = json.Unmarshal(resp.Body(), &data)
		if err != nil {
			return 0, fmt.Errorf("Error parsing JSON response: %v", err)
		}

		// Extract and print the stock price
		if len(data.Chart.Result) > 0 {
			price := data.Chart.Result[0].Meta.RegularMarketPrice
			return price, nil
		} else {
			return 0, fmt.Errorf("Stock price not found")
		}
	}
}

func getPrevious(symbol string, interval IntervalRange) (float64, error) {
	for {
		url := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s?interval=1d&range=%s", symbol, interval)
		request := getClient(symbol)
		resp, err := request.Get(url)
		if err != nil {
			return 0, fmt.Errorf("Error fetching stock price: %v", err)
		}

		// Check for "Too Many Requests" response
		if resp.StatusCode() == 429 {
			log.Info("Rate limit exceeded. Retrying after a delay...\n")
			time.Sleep(10 * time.Second) // Wait for 10 seconds before retrying
			continue
		}

		// Parse the JSON response
		var data YahooFinanceResponse
		err = json.Unmarshal(resp.Body(), &data)
		if err != nil {
			return 0, fmt.Errorf("Error parsing JSON response: %v", err)
		}

		// Calculate the 6-month average
		if len(data.Chart.Result) > 0 {
			closes := data.Chart.Result[0].Indicators.Quote[0].Close
			var sum float64
			var count int
			for _, close := range closes {
				if close > 0 {
					sum += close
					count++
				}
			}
			if count > 0 {
				average := sum / float64(count)
				return average, nil
			} else {
				return 0, fmt.Errorf("No valid closing prices found")
			}
		} else {
			return 0, fmt.Errorf("No data found")
		}
	}
}

func getClient(symbol string) *resty.Request {
	// Create a new Resty client
	client := resty.New()
	return client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/113.0")
}
