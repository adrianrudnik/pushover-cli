package cmd

import (
	"encoding/json"
	"errors"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func extractResponse(response *http.Response) (*ApiResponse, error) {
	apiResponse := &ApiResponse{}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	log.Debug().
		Str("status", response.Status).
		Int("code", response.StatusCode).
		Msg("HTTP response received")

	err = json.Unmarshal(body, &apiResponse)

	if err != nil {
		return nil, err
	}

	if apiResponse.Status != 1 {
		log.Error().
			Int("status", apiResponse.Status).
			Strs("errors", apiResponse.Errors).
			Str("request", apiResponse.Request).
			Msg("API reported errors")

		return nil, errors.New("API reported bad request")
	}

	return apiResponse, nil
}

func extractRateLimits(response *http.Response) (*ApiRateLimit, error) {
	limits := &ApiRateLimit{
		RequestsTotalPerMonth: 0,
		RequestsRemaining:     0,
		ResetAt:               time.Now(),
	}

	rateAppLimit := response.Header.Get("X-Limit-App-Limit")

	if rateAppLimit != "" {
		v, err := strconv.ParseInt(rateAppLimit, 10, 64)

		if err != nil {
			return nil, errors.New("could not parse X-Limit-App-Limit header")
		} else {
			limits.RequestsTotalPerMonth = v
		}
	}

	rateAppRemaining := response.Header.Get("X-Limit-App-Remaining")

	if rateAppRemaining != "" {
		v, err := strconv.ParseInt(rateAppRemaining, 10, 64)

		if err != nil {
			return nil, errors.New("could not parse X-Limit-App-Remaining header")
		} else {
			limits.RequestsRemaining = v
		}
	}

	rateAppReset := response.Header.Get("X-Limit-App-Reset")

	if rateAppReset != "" {
		v, err := strconv.ParseInt(rateAppReset, 10, 64)

		if err != nil {
			return nil, errors.New("could not parse \"X-Limit-App-Reset header")
		} else {
			limits.ResetAt = time.Unix(v, 0).UTC()
		}
	}

	return limits, nil
}
