package cmd

import (
	"encoding/json"
	"errors"
	"github.com/rs/zerolog/log"
	"github.com/shibukawa/configdir"
	"os"
	"path"
)

func GetConfigDir() configdir.ConfigDir {
	return configdir.New("", ApplicationName)
}

// LoadConfig tries to load the config that could be stored in
// many different locations based on the OS.
// It returns the current config and a flag if the config was successfully loaded.
func LoadConfig() (config *Config, loaded bool) {
	config = &Config{
		UserKey:  "",
		ApiToken: "",
	}

	// Probe env first
	if os.Getenv("PUSHOVER_CLI_USER") != "" || os.Getenv("PUSHOVER_CLI_API") != "" {
		log.Debug().Msg("Detected ENV variables, ignoring configuration file")

		config.UserKey = os.Getenv("PUSHOVER_CLI_USER")
		config.ApiToken = os.Getenv("PUSHOVER_CLI_API")

		return config, true
	}

	loaded = false

	// Try to find a location for the config file
	locations := GetConfigDir()
	location := locations.QueryFolderContainsFile(configFile)

	// Debug possible config folder locations
	if VerboseFlag {
		folders := locations.QueryFolders(configdir.All)
		for _, v := range folders {
			log.Debug().Str("path", v.Path).Msg("Possible config folder")
		}
	}

	// If we found a config file, try to parse it
	if location != nil {
		log.Debug().Str("path", path.Join(location.Path, configFile)).Msg("Found config file")

		data, err := location.ReadFile(configFile)

		if err != nil {
			log.Fatal().
				Err(err).
				Msg("Failed to read config")
		}

		err = json.Unmarshal(data, &config)

		if err != nil {
			log.Fatal().
				Err(err).
				Msgf("Failed to parse config json")
		}

		log.Debug().Interface("config", config).Msg("Config loaded")

		return config, true
	}

	// @todo try env?

	return config, false
}

func SaveConfig(config *Config) {
	locations := GetConfigDir()
	folders := locations.QueryFolders(configdir.All)

	if len(folders) == 0 {
		log.Fatal().Msg("Could not guess a fitting config folder")
	}

	data, err := json.Marshal(config)

	if err != nil {
		log.Fatal().Err(err).Msg("Could not convert current config into json")
	}

	err = folders[0].WriteFile(configFile, data)

	if err != nil {
		log.Fatal().Err(err).Str("path", folders[0].Path).Msg("Could not write config file")
	}

	log.Info().Str("path", path.Join(folders[0].Path, configFile)).Msg("Config saved")
}

// ClearConfig will remove the first configuration file that can be found from the system.
// It will return the path that was removed and a possible error.
func ClearConfig() (string, error) {
	locations := GetConfigDir()
	folders := locations.QueryFolders(configdir.All)

	if len(folders) == 0 {
		return "", nil
	}

	location := locations.QueryFolderContainsFile(configFile)

	if location == nil {
		return "", errors.New("no config file found")
	}

	filePath := path.Join(folders[0].Path, configFile)

	err := os.Remove(filePath)

	if err != nil {
		return filePath, err
	}

	return filePath, nil
}
