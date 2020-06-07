package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"net/http"
	"net/url"
)

var limitsCmd = &cobra.Command{
	Use:   "limits",
	Short: "Prints the current API limits.",
	Run: func(cmd *cobra.Command, args []string) {
		config, loaded := LoadConfig()

		if !loaded {
			log.Fatal().Msg("No configuration found")
		}

		url, err := url.ParseRequestURI(LimitsApiEndpoint)

		if err != nil {
			log.Fatal().Err(err).Msg("Failed to parse limit API endpoint")
		}

		query := url.Query()
		query.Set("token", config.ApiToken)
		url.RawQuery = query.Encode()

		log.Debug().Str("API request target", url.String()).Msg("API test")

		response, err := http.Get(url.String())

		if err != nil {
			log.Fatal().Err(err).Msg("Failed to GET limits API")
		}

		log.Debug().Int("status", response.StatusCode).Msg("API response received")

		apiLimits, err := extractRateLimits(response)

		if err != nil {
			log.Fatal().Err(err).Msg("Could not extract API limit information")
		}

		log.Info().
			Int64("requests-remaining", apiLimits.RequestsRemaining).
			Int64("requests-per-month", apiLimits.RequestsTotalPerMonth).
			Time("reset-at", apiLimits.ResetAt).
			Msg("Rate limit information")

	},
}
