package core_controller

import (
	"os"

	"github.com/rs/zerolog"
)

type CoreService struct {
	log     zerolog.Logger
	envVars map[string]string
}

func NewCoreService(log zerolog.Logger) *CoreService {

	// err := godotenv.Load()
	// if err != nil {
	// 	log.Info().Msg("Error loading .env file")
	// }

	requiredKeys := []string{
		"API_KEY",
		"BASE_URI",
		"ELEC_MPAN",
		"ELEC_SERIAL",
		"GAS_MPRN",
		"GAS_SERIAL",
	}

	envVarsMap := make(map[string]string)

	for _, key := range requiredKeys {
		value, ok := os.LookupEnv(key)
		if !ok || value == "" {
			log.Fatal().Str("env", key).Msg("Required environment variable is not set or empty")
		}
		envVarsMap[key] = value
	}

	log.Info().Msg("All required environment variables loaded successfully")

	return &CoreService{
		log:     log,
		envVars: envVarsMap,
	}
}

// Getter for main to access to the Endpoints
func (c *CoreService) EnvironmentVariables() map[string]string {
	return c.envVars
}
