package cmd

import "time"

type Config struct {
	UserKey     string `json:"user-key"`
	ApiToken    string `json:"api-token"`
}

type ApiResponse struct {
	Status int `json:"status"`
	Request string `json:"request"`
	Errors []string `json:"errors,omitempty"`
}

type ApiRateLimit struct {
	RequestsTotalPerMonth int64
	RequestsRemaining int64
	ResetAt time.Time
}
