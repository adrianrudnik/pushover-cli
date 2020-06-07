package cmd

import (
	"bytes"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"
)

var pushCmd = &cobra.Command{
	Use:     "push [message]",
	Short:   "Pushes the given text message",
	Long:    "Pushes the given text message. Will be truncated to a maximum length of 1024 characters of UTF-8.",
	Example: fmt.Sprintf("%s -v push --devices=mobile,workpc -t WARNING --priority=high \"The following error occured: example\"", CommandName),
	Args:    cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		// Validate api endpoint
		_, err := url.ParseRequestURI(PushApiEndpoint)
		if err != nil {
			log.Fatal().Msg("API endpoint url invalid")
		}

		// Validate priority
		_, ok := priorities[PushPriority]
		if !ok {
			log.Fatal().Str("priority", PushPriority).Msg("Unknown priority")
		}

		// Validate sounds
		_, ok = sounds[PushSound]
		if !ok {
			log.Fatal().Str("sound", PushSound).Msg("Unknown sound")
		}

		// Validate attachment
		if PushAttachment != "" {
			finfo, err := os.Stat(PushAttachment)
			if os.IsNotExist(err) {
				log.Fatal().
					Str("attachment", PushAttachment).
					Msg("Attachment file not found")
			}

			if finfo.Size() > maxAttachmentSize {
				log.Fatal().
					Str("attachment", PushAttachment).
					Int64("size-file", finfo.Size()).
					Int64("size-allowed", maxAttachmentSize).
					Msg("Attachment too large")
			}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		config, loaded := LoadConfig()
		if !loaded {
			log.Fatal().Msg("No configuration found")
		}

		// Truncate fields that are too long and would prevent the push from being sent
		if utf8.RuneCountInString(args[0]) > 1024 {
			log.Warn().Msg("Message > 1024 characters, will be truncated")
			args[0] = string([]rune(args[0])[0:1024])
		}

		if utf8.RuneCountInString(PushTitle) > 250 {
			log.Warn().Msg("Title > 250 characters, will be truncated")
			PushTitle = string([]rune(PushTitle)[0:250])
		}

		if utf8.RuneCountInString(PushUrl) > 512 {
			log.Warn().Msg("Link url > 512 characters, will be truncated")
			PushUrl = string([]rune(PushUrl)[0:512])
		}

		if utf8.RuneCountInString(PushUrlTitle) > 100 {
			log.Warn().Msg("Link title > 100 characters, will be truncated")
			PushUrlTitle = string([]rune(PushUrlTitle)[0:100])
		}

		// Prepare form
		form := url.Values{}

		// Required arguments
		form.Add("token", config.ApiToken)
		form.Add("user", config.UserKey)
		form.Add("message", args[0])

		// Optional arguments
		if PushPriority != "normal" {
			form.Add("priority", strconv.FormatInt(int64(priorities[PushPriority]), 10))
		}

		if len(PushToDevices) > 0 {
			form.Add("device", strings.Join(PushToDevices, ","))
		}

		if PushTitle != "" {
			form.Add("title", PushTitle)
		}

		if PushUrl != "" {
			form.Add("url", PushUrl)

			if PushUrlTitle != "" {
				form.Add("url_title", PushUrlTitle)
			}
		}

		if PushSound != "pushover" {
			form.Add("sound", PushSound)
		}

		if PushTimestamp != 0 {
			form.Add("timestamp", strconv.FormatInt(PushTimestamp, 10))
		}

		client := &http.Client{}
		var request *http.Request

		if PushAttachment != "" {
			// Switch to multipart form
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			// Convert basic form values to multipart form
			for f, v := range form {
				err := writer.WriteField(f, v[0])
				if err != nil {
					log.Warn().Err(err).Str("field", f).Msg("Failed to convert form field to multipart field")
				}
			}

			// Check file type
			_, err := parseAttachment(PushAttachment)
			if err != nil {
				log.Fatal().Err(err).Msg("File is not invalid and can not be pushed")
			}
			part, err := writer.CreateFormFile("attachment", filepath.Base(PushAttachment))

			// Load file into form
			file, err := os.Open(PushAttachment)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to open file attachment")
			}

			_, err = io.Copy(part, file)
			if err != nil {
				file.Close()
				log.Fatal().Err(err).Msg("Failed to copy attachment stream")
			}

			err = file.Close()
			if err != nil {
				log.Warn().Err(err).Msg("Failed to close file attachment")
			}

			err = writer.Close()
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to close multipart writer")
			}

			r, err := http.NewRequest("POST", PushApiEndpoint, body)
			if err != nil {
				log.Fatal().Err(err).Msg("HTTP request could not be initialized")
			}

			r.Header.Add("Content-Type", writer.FormDataContentType())

			request = r
		} else {
			// Simple form submission
			encoded := form.Encode()
			r, err := http.NewRequest("POST", PushApiEndpoint, strings.NewReader(encoded))
			if err != nil {
				log.Fatal().Err(err).Msg("HTTP request could not be initialized")
			}

			// Set form headers
			r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			r.Header.Add("Content-Length", strconv.Itoa(len(encoded)))

			request = r
		}

		response, err := client.Do(request)

		if err != nil {
			log.Fatal().Err(err).Msg("HTTP request failed")
		}

		apiResponse, err := extractResponse(response)

		if err != nil {
			log.Fatal().Err(err).Msg("API request failed")
		}

		log.Info().
			Int("status", apiResponse.Status).
			Str("request", apiResponse.Request).
			Msg("Message pushed")

		apiLimits, err := extractRateLimits(response)

		if err != nil {
			log.Warn().Err(err).Msg("Could not extract API limit information")
		}

		if apiLimits != nil {
			log.Info().
				Int64("requests-remaining", apiLimits.RequestsRemaining).
				Int64("requests-per-month", apiLimits.RequestsTotalPerMonth).
				Time("reset-at", apiLimits.ResetAt).
				Msg("Rate limit information")
		}
	},
}
