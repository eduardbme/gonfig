package config

import (
	"encoding/json"
	"fmt"
	"path"

	"github.com/eduardbcom/gonfig/internal"
)

// Read ...
func Read() ([]byte, error) {
	var err error

	var configDirPath string

	configDirPath, err = internal.GetConfigDirPath()

	if err != nil {
		return nil, err
	}

	appEnv := internal.GetAppEnv()

	var defaultConfig *internal.Config
	var envConfig *internal.Config
	var localEnvConfig *internal.Config

	defaultConfig, err = internal.NewConfig(path.Join(configDirPath, "default.json"))
	if err != nil {
		return nil, err
	}

	if len(appEnv) > 0 {
		envConfig, err = internal.NewConfig(path.Join(configDirPath, fmt.Sprintf("%s.json", appEnv)))
		if err != nil {
			return nil, err
		}

		localEnvConfig, err = internal.NewConfig(path.Join(configDirPath, fmt.Sprintf("local-%s.json", appEnv)))
		if err != nil {
			return nil, err
		}
	}

	entireConfig := internal.JoinConfigs(defaultConfig, envConfig, localEnvConfig)

	var configSchema *internal.Schema

	configSchema, err = internal.NewSchema(path.Join(configDirPath, "schema"))
	if err != nil {
		return nil, err
	}

	if configSchema != nil {
		err = configSchema.Validate(entireConfig)
		if err != nil {
			return nil, err
		}
	}

	return json.Marshal(entireConfig)
}
