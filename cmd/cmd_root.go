package cmd

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var (
	VerboseFlag bool

	PushApiEndpoint string
	PushToDevices   []string
	PushTitle       string
	PushUrl         string
	PushUrlTitle    string
	PushPriority    string
	PushSound       string
	PushTimestamp   int64
	PushAttachment  string

	LimitsApiEndpoint string

	priorities = map[string]int{
		"none":    -2,
		"quiet":   -1,
		"normal":  0,
		"high":    1,
		"confirm": 2,
	}

	sounds = map[string]bool{
		"pushover":     true,
		"bike":         true,
		"bugle":        true,
		"cashregister": true,
		"classical":    true,
		"cosmic":       true,
		"falling":      true,
		"gamelan":      true,
		"incoming":     true,
		"intermission": true,
		"magic":        true,
		"mechanical":   true,
		"pianobar":     true,
		"siren":        true,
		"spacealarm":   true,
		"tugboat":      true,
		"alien":        true,
		"climb":        true,
		"persistent":   true,
		"echo":         true,
		"updown":       true,
		"vibrate":      true,
		"none":         true,
	}
)

func Execute() {
	rootCmd.Version = Version
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	pushCmd.ResetFlags()
	pushCmd.ResetCommands()
	pushCmd.PersistentFlags().StringVarP(&PushApiEndpoint, "api-endpoint", "e", "https://api.pushover.net/1/messages.json", "API endpoint for message submission")
	pushCmd.PersistentFlags().StringSliceVarP(&PushToDevices, "devices", "d", []string{}, "devices to limit the push to (comma-separated)")
	pushCmd.PersistentFlags().StringVarP(&PushTitle, "title", "t", "", "message title (max. 250 characters)")
	pushCmd.PersistentFlags().StringVar(&PushUrl, "link-url", "", "supplementary URL (max. 512 characters)")
	pushCmd.PersistentFlags().StringVar(&PushUrlTitle, "link-label", "", "title for the supplementary URL (max. 100 characters)")
	pushCmd.PersistentFlags().StringVarP(&PushPriority, "priority", "p", "normal", "message priority [none, quiet, normal, high, confirm]")
	pushCmd.PersistentFlags().StringVarP(&PushSound, "sound", "s", "pushover", "playback sound [see https://pushover.net/api#sounds]")
	pushCmd.PersistentFlags().Int64Var(&PushTimestamp, "timestamp", 0, "message date and time override as unix timestamp")
	pushCmd.PersistentFlags().StringVarP(&PushAttachment, "attachment", "a", "", "path to image attachment (max size 2.5mb)")

	limitsCmd.ResetFlags()
	limitsCmd.ResetCommands()
	limitsCmd.PersistentFlags().StringVarP(&LimitsApiEndpoint, "api-endpoint", "e", "https://api.pushover.net/1/apps/limits.json", "API endpoint for limit requests")

	configCmd.ResetFlags()
	configCmd.ResetCommands()
	configCmd.AddCommand(configSetupCmd)
	configCmd.AddCommand(configShowPathsCmd)
	configCmd.AddCommand(configClearCmd)

	// root
	rootCmd.ResetFlags()
	rootCmd.ResetCommands()
	rootCmd.PersistentFlags().BoolVarP(&VerboseFlag, "verbose", "v", false, "print debug information")

	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(pushCmd)
	rootCmd.AddCommand(limitsCmd)
}

var rootCmd = &cobra.Command{
	Use:   CommandName,
	Short: fmt.Sprintf("%s is fast and static cli tool to send push notifications over pushover.net.", CommandName),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Configure logger
		if !VerboseFlag {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		} else {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		}

		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	},
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}
