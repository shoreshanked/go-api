package main

import (
	"fmt"
	api_controller "go-api/pkg/apis"
	time_controller "go-api/pkg/datetime"
	"os"

	"github.com/joho/godotenv"
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

func initApiDetails() ApiDetails {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	apiKey := os.Getenv("API_KEY")
	baseUri := os.Getenv("BASE_URI")
	elecMPAN := os.Getenv("ELEC_MPAN")
	elecSerial := os.Getenv("ELEC_SERIAL")
	gasMPRN := os.Getenv("GAS_MPRN")
	gasSerial := os.Getenv("GAS_SERIAL")

	elecEndpoint := fmt.Sprintf("%selectricity-meter-points/%s/meters/%s/consumption/", baseUri, elecMPAN, elecSerial)
	gasEndpoint := fmt.Sprintf("%sgas-meter-points/%s/meters/%s/consumption/", baseUri, gasMPRN, gasSerial)

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

	apiDetails := initApiDetails()
	from, to := time_controller.GetTimeRange()

	for idx, endpoint := range apiDetails.Endpoints {

		fmt.Printf("%v\t%v\n", idx, endpoint)
		resp, err := api_controller.GetData(apiDetails.ApiKey, endpoint.Url, to, from)

		if err != nil {
			// Print the error or log it; for now, let's print it
			fmt.Println("Error:", err)
			return // Exit or return from the function if there's an error
		}

		formattedJson, totalConsumption := api_controller.ParseApiResponse(resp)

		fmt.Println(formattedJson)
		fmt.Printf("Total %s consumption: %.3f kWh\n", endpoint.Type, totalConsumption)
	}

}
