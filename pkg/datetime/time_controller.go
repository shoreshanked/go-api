package time_controller

import (
	"time"

	"github.com/rs/zerolog"
)

type TimeService struct {
	log zerolog.Logger
}

func NewTimeService(log zerolog.Logger) *TimeService {
	return &TimeService{log: log}
}

func (s *TimeService) GetTimeRange() (string, string) {
	s.log.Info().Msg("Calculating timeframe for data retrieval")
	// Define the GMT+1 timezone (British Summer Time)
	loc, err := time.LoadLocation("Europe/London") // Use this location to ensure correct timezone handling
	if err != nil {
		s.log.Error().Err(err).Msg("Error Loading Timezone")
		return "", ""
	}

	// Get today's date at midnight in GMT+1 timezone
	now := time.Now().In(loc)
	to := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc) // Midnight today in GMT+1

	// Subtract 1 day to get yesterday's midnight in GMT+1 timezone
	yesterday := to.AddDate(0, 0, -1)
	dayBeforeYesterday := to.AddDate(0, 0, -2)

	// Return formatted dates as strings in RFC3339 format, with time zone offset
	return dayBeforeYesterday.Format(time.RFC3339), yesterday.Format(time.RFC3339)
}
