package main

import (
	"fmt"

	api_controller "go-api/pkg/apis"
	core_controller "go-api/pkg/core"
	time_controller "go-api/pkg/datetime"
	log_controller "go-api/pkg/logging"
)

func main() {
	// Setup Logging
	logService := log_controller.NewLogService("info")
	log := logService.Logger()
	// Setup Services
	coreService := core_controller.NewCoreService(log)
	apiService := api_controller.NewAPIService(log, coreService.EnvironmentVariables())
	timeService := time_controller.NewTimeService(log)

	// Retrieve method parameters
	from, to := timeService.GetTimeRange()

	log.Info().
		Str("Getting data from:", from).
		Str("To", to).
		Msg("Time range")

	// Get Data for each specified endpoint
	for idx, endpoint := range apiService.Endpoints() {

		log.Debug().
			Int("Index:", idx).
			Str("Name", endpoint.Type).
			Str("URL", endpoint.Url).Msg("Logging endpoint")

		resp, err := apiService.GetData(endpoint.Url, to, from)
		if err != nil {
			// Print the error or log it; for now, let's print it
			log.Error().Err(err).Msg("")
			return // Exit or return from the function if there's an error
		}

		formattedJson, totalConsumption := apiService.ParseApiResponse(resp)

		log.Debug().
			Msgf("Response data:\n%s", formattedJson)

		log.Info().
			Str("Endpoint:", endpoint.Type).
			Str("Total consumption", fmt.Sprintf("%.2f kwh", totalConsumption)).
			Msg("Consumption data logged")
	}
}
