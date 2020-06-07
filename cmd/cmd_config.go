package cmd

import (
	"bufio"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/shibukawa/configdir"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strings"
	"syscall"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration commands",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var configSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup the auth configuration for current user.",
	Long:  "Configure the auth configuration and store inside the current users profile settings.",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Enter user key: ")
		userKey, err := reader.ReadString('\n')

		if err != nil {
			log.Fatal().Err(err).Msg("Failed to read input")
		}

		userKey = strings.TrimSpace(userKey)

		fmt.Print("Enter API token: ")
		apiToken, err := terminal.ReadPassword(int(syscall.Stdin))

		if err != nil {
			log.Fatal().Err(err).Msg("Failed to read password input on terminal")
		}

		config := &Config{
			UserKey:  strings.TrimSpace(userKey),
			ApiToken: strings.TrimSpace(string(apiToken)),
		}

		log.Debug().Interface("config", config).Msg("Given config")

		SaveConfig(config)
	},
}

var configShowPathsCmd = &cobra.Command{
	Use:   "paths",
	Short: "Prints the available paths for configuration.",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msgf("Collecting paths that will be used for %s lookup", configFile)

		locations := GetConfigDir()
		folders := locations.QueryFolders(configdir.All)

		for _, v := range folders {
			log.Info().Str("path", v.Path).Msg("Folder found")
		}
	},
}

var configClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Removes the next local configuration file.",
	Run: func(cmd *cobra.Command, args []string) {
		path, err := ClearConfig()

		if err != nil {
			log.Fatal().
				Err(err).
				Str("file", path).
				Msg("Could not removed config file")
		}

		log.Info().Str("file", path).Msg("Config cleared")
	},
}

