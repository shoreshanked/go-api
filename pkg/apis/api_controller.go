package api_controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// type of retryablefunction = a type that is a function and all it does is return an error
type RetryableFunc func() error

const defaultMaxRetries = 2
const defaultBaseDelay = 2 * time.Second

type ConsumptionResult struct {
	Consumption   float64 `json:"consumption"`
	IntervalStart string  `json:"interval_start"`
	IntervalEnd   string  `json:"interval_end"`
}

type ApiResponse struct {
	Count    int                 `json:"count"`
	Next     *string             `json:"next"`     // use *string because it can be null
	Previous *string             `json:"previous"` // use *string because it can be null
	Results  []ConsumptionResult `json:"results"`
}

func GetData(apiKey string, endpoint string, to string, from string) ([]byte, error) {

	var resultBody []byte

	// retryableFunc is equal to a function that returns an error
	retryableFunc := func() error {

		req, err := http.NewRequest("GET", endpoint, nil)
		if err != nil {
			return err
		}

		q := req.URL.Query()
		q.Add("period_from", from)
		if to != "" {
			q.Add("period_to", to)
		}
		req.URL.RawQuery = q.Encode()

		req.SetBasicAuth(apiKey, "")
		client := &http.Client{}

		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode >= 400 {
			return fmt.Errorf(".statusCode %d", resp.StatusCode)
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		fmt.Println("Data received successfully")
		resultBody = body
		return nil
	}

	err := doWithRetries(retryableFunc, defaultMaxRetries, defaultBaseDelay) // Passing 0 uses default values
	if err != nil {
		return nil, err // If retries failed, return the error
	}

	return resultBody, nil
}

func ParseApiResponse(resultBody []byte) (string, float64) {
	response := unmarshalApiResponse(resultBody)
	totalConsumption := sumConsumption(response)
	prettyJSON := marshalIndentResponse(response)

	return string(prettyJSON), totalConsumption
}

func unmarshalApiResponse(resultBody []byte) []ConsumptionResult {
	var response ApiResponse
	err := json.Unmarshal(resultBody, &response)
	if err != nil {
		fmt.Printf("Failed to unmarshal JSON: %v", err)
	}
	return response.Results
}

func sumConsumption(entries []ConsumptionResult) float64 {
	var total float64
	for _, entry := range entries {
		total += entry.Consumption
	}
	return total
}

func marshalIndentResponse(entries []ConsumptionResult) string {
	prettyJSON, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Sprintf("Failed to format JSON: %v", err)
	}
	return string(prettyJSON)
}

func doWithRetries(fn RetryableFunc, max_retries int, base_delay time.Duration) error {
	var last_err error
	for i := 0; i < max_retries; i++ {
		last_err = fn()
		if last_err == nil {
			return nil
		}
		delay := base_delay * (1 << i)
		fmt.Printf("Retrying in %v...\n", delay)
		time.Sleep(delay)
	}
	return last_err
}
