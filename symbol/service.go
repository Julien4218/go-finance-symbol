package symbol

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/Julien4218/go-finance-symbol/observability"
)

var logf = observability.Logf
var log = observability.Log
var restyClientFactoryFunc = getRestyClient
var baseUrl = "https://query1.finance.yahoo.com/v8/finance/chart"

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

func Execute(symbol string, intervalRanges []IntervalRange) {
	price, err := getToday(symbol)
	if err != nil {
		logf("%s", err)
	} else {
		logf("Stock price for %s: $%.2f", symbol, price)
		observability.GetOrCreateGauge(fmt.Sprintf("%s_1d", symbol)).Set(price)
	}

	for _, r := range intervalRanges {
		gather(symbol, r)
	}
}

func gather(symbol string, interval IntervalRange) {
	average, err := getPrevious(symbol, interval)
	if err != nil {
		logf("%s", err)
	} else {
		logf("%s average price for %s: $%.2f", interval, symbol, average)
		observability.GetOrCreateGauge(fmt.Sprintf("%s_%s", symbol, interval)).Set(average)
	}
}

func getToday(symbol string) (float64, error) {
	for {
		url := fmt.Sprintf("%s/%s", baseUrl, symbol)
		request := getClient()
		resp, err := request.Get(url)
		if err != nil {
			return 0, fmt.Errorf("Error fetching stock price: %v", err)
		}

		// Check for "Too Many Requests" response
		if resp.StatusCode() == 429 {
			log("Rate limit exceeded. Retrying after a delay...")
			time.Sleep(10 * time.Second) // Wait for 10 seconds before retrying
			continue
		}

		var data YahooFinanceResponse
		err = json.Unmarshal(resp.Body(), &data)
		if err != nil {
			return 0, fmt.Errorf("Error parsing JSON response: %v", err)
		}

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
		url := fmt.Sprintf("%s/%s?interval=1d&range=%s", baseUrl, symbol, interval)
		request := getClient()
		resp, err := request.Get(url)
		if err != nil {
			return 0, fmt.Errorf("Error fetching stock price: %v", err)
		}

		// Check for "Too Many Requests" response
		if resp.StatusCode() == 429 {
			log("Rate limit exceeded. Retrying after a delay...")
			time.Sleep(10 * time.Second) // Wait for 10 seconds before retrying
			continue
		}

		var data YahooFinanceResponse
		err = json.Unmarshal(resp.Body(), &data)
		if err != nil {
			return 0, fmt.Errorf("Error parsing JSON response: %v", err)
		}

		if len(data.Chart.Result) > 0 && len(data.Chart.Result[0].Indicators.Quote) > 0 {
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

func getClient() *resty.Request {
	client := restyClientFactoryFunc()
	return client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/113.0")
}

func getRestyClient() *resty.Client {
	return resty.New()
}
