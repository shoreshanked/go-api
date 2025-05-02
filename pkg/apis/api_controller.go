package api_controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

// type of retryablefunction = a type that is a function and all it does is return an error
type RetryableFunc func() error

const (
	defaultMaxRetries = 2
	defaultBaseDelay  = 2 * time.Second
)

type APIService struct {
	log        zerolog.Logger
	apiDetails ApiDetails
}

type EndpointInfo struct {
	Type string
	Url  string
}

type ApiDetails struct {
	ApiKey    string
	BaseUri   string
	Endpoints []EndpointInfo
}

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

// Getter for main to access to the Endpoints
func (s *APIService) Endpoints() []EndpointInfo {
	return s.apiDetails.Endpoints
}

func NewAPIService(log zerolog.Logger, envVars map[string]string) *APIService {
	log.Info().Msg("Initialising API details from Env vars")

	elecEndpoint := fmt.Sprintf("%selectricity-meter-points/%s/meters/%s/consumption/",
		envVars["BASE_URI"], envVars["ELEC_MPAN"], envVars["ELEC_SERIAL"])
	gasEndpoint := fmt.Sprintf("%sgas-meter-points/%s/meters/%s/consumptionn/",
		envVars["BASE_URI"], envVars["GAS_MPRN"], envVars["GAS_SERIAL"])

	endpoints := []EndpointInfo{
		{Type: "electric", Url: elecEndpoint},
		{Type: "gas", Url: gasEndpoint},
	}

	apiDetails := ApiDetails{
		ApiKey:    envVars["API_KEY"],
		BaseUri:   envVars["BASE_URI"],
		Endpoints: endpoints,
	}

	return &APIService{
		log:        log,
		apiDetails: apiDetails,
	}
}

func (s *APIService) GetData(endpoint string, to string, from string) ([]byte, error) {
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

		req.SetBasicAuth(s.apiDetails.ApiKey, "")
		client := &http.Client{}

		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode >= 400 {
			return fmt.Errorf(".statusCode %d", resp.StatusCode)
		}

		defer func() {
			_ = resp.Body.Close()
		}()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		s.log.Info().
			Msg("Data recieved succesfully")

		resultBody = body
		return nil
	}

	err := s.doWithRetries(retryableFunc, defaultMaxRetries, defaultBaseDelay) // Passing 0 uses default values
	if err != nil {
		return nil, err // If retries failed, return the error
	}

	return resultBody, nil
}

func (s *APIService) ParseApiResponse(resultBody []byte) (string, float64) {
	response := s.unmarshalApiResponse(resultBody)
	totalConsumption := s.sumConsumption(response)
	prettyJSON := s.marshalIndentResponse(response)
	return string(prettyJSON), totalConsumption
}

func (s *APIService) unmarshalApiResponse(resultBody []byte) []ConsumptionResult {
	var response ApiResponse
	err := json.Unmarshal(resultBody, &response)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to unmarshal JSON")
	}
	return response.Results
}

func (s *APIService) sumConsumption(entries []ConsumptionResult) float64 {
	var total float64
	for _, entry := range entries {
		total += entry.Consumption
	}
	return total
}

func (s *APIService) marshalIndentResponse(entries []ConsumptionResult) string {
	prettyJSON, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Sprintf("Failed to format JSON: %v", err)
	}
	return string(prettyJSON)
}

func (s *APIService) doWithRetries(fn RetryableFunc, max_retries int, base_delay time.Duration) error {
	var last_err error
	for i := 0; i < max_retries; i++ {
		last_err = fn()
		if last_err == nil {
			return nil
		}
		delay := base_delay * (1 << i)
		s.log.Info().
			Str("retrying_in", delay.String()).
			Msg("Retrying...")
		time.Sleep(delay)
	}
	return last_err
}
