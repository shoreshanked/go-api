package main

import (
	"fmt"
	"os"

	api_controller "go-api/pkg/apis"
	time_controller "go-api/pkg/datetime"
	log_controller "go-api/pkg/logging"

	"github.com/rs/zerolog"
)

type EndpointInfo struct {
	Type string
	Url  string
}

type ApiDetails struct {
	ApiKey    string
	BaseUri   string
	Endpoints []EndpointInfo
}

func initApiDetails(log zerolog.Logger) ApiDetails {
	// err := godotenv.Load()
	// if err != nil {
	// 	fmt.Println("Error loading .env file")
	// }

	log.Info().Msg("Initialising API details from Env vars")

	apiKey := os.Getenv("API_KEY")
	baseUri := os.Getenv("BASE_URI")
	elecMPAN := os.Getenv("ELEC_MPAN")
	elecSerial := os.Getenv("ELEC_SERIAL")
	gasMPRN := os.Getenv("GAS_MPRN")
	gasSerial := os.Getenv("GAS_SERIAL")

	elecEndpoint := fmt.Sprintf("%selectricity-meter-points/%s/meters/%s/consumption/", baseUri, elecMPAN, elecSerial)
	gasEndpoint := fmt.Sprintf("%sgas-meter-points/%s/meters/%s/consumptionn/", baseUri, gasMPRN, gasSerial)

	endpoints := []EndpointInfo{
		{Type: "electric", Url: elecEndpoint},
		{Type: "gas", Url: gasEndpoint},
	}

	return ApiDetails{
		ApiKey:    apiKey,
		BaseUri:   baseUri,
		Endpoints: endpoints,
	}
}

func main() {

	log := log_controller.SetupLogger()
	api_controller.InitAPIService(log)
	time_controller.InitTimeService(log)

	apiDetails := initApiDetails(log)
	from, to := time_controller.GetTimeRange()

	log.Info().
		Str("Getting data from:", from).
		Str("To", to).
		Msg("Time range")

	for idx, endpoint := range apiDetails.Endpoints {

		log.Info().
			Int("Index:", idx).
			Str("Name", endpoint.Type).
			Str("URL", endpoint.Url).Msg("Logging endpoint")

		resp, err := api_controller.GetData(apiDetails.ApiKey, endpoint.Url, to, from)
		if err != nil {
			// Print the error or log it; for now, let's print it
			log.Error().Err(err).Msg("")
			return // Exit or return from the function if there's an error
		}

		formattedJson, totalConsumption := api_controller.ParseApiResponse(resp)

		log.Info().
			Msgf("Response data:\n%s", formattedJson)

		log.Info().
			Str("Endpoint:", endpoint.Type).
			Str("Total consumption", fmt.Sprintf("%.2f kwh", totalConsumption)).
			Msg("Consumption data logged")
	}
}
