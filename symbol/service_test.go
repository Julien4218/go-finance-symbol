package symbol

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
)

func TestExecute(t *testing.T) {
	factory, client := givenTestClientFunc()
	restyClientFactoryFunc = factory
	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	baseUrl = "https://example.com"
	httpmock.RegisterResponder("GET", baseUrl+"/TEST",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"chart": map[string]interface{}{
					"result": []map[string]interface{}{
						{
							"meta": map[string]interface{}{
								"regularMarketPrice": 100.0,
							},
						},
					},
				},
			})
			return resp, err
		})

	// // Mock the response for getPrevious
	httpmock.RegisterResponder("GET", baseUrl+"/TEST?interval=1d&range=5d",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"chart": map[string]interface{}{
					"result": []map[string]interface{}{
						{
							"indicators": map[string]interface{}{
								"quote": []map[string]interface{}{
									{
										"close": []float64{100.0, 101.0, 102.0, 103.0, 104.0},
									},
								},
							},
						},
					},
				},
			})
			return resp, err
		})

	var logs []string
	mockLogf := func(format string, args ...interface{}) {
		logs = append(logs, fmt.Sprintf(format, args...))
	}
	mockLog := func(message string) {
		logs = append(logs, message)
	}
	logf = mockLogf
	log = mockLog

	Execute("TEST", []IntervalRange{FiveDay})

	expectedLogs := []string{
		"Stock price for TEST: $100.00",
		"5d average price for TEST: $102.00",
	}

	for i, log := range expectedLogs {
		if logs[i] != log {
			t.Errorf("Expected log %q, got %q", log, logs[i])
		}
	}

}

func givenTestClientFunc() (func() *resty.Client, *resty.Client) {
	client := resty.New()
	defaultLoader := func() *resty.Client {
		return client
	}
	return defaultLoader, client
}
