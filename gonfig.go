package config

import (
	"encoding/json"
	"fmt"
	"path"

	"github.com/eduardbcom/gonfig/internal"
)

// Read ...
func Read() ([]byte, error) {
	if configDirPath, err := internal.GetConfigDirPath(); err != nil {
		return nil, err
	} else {
		var configs []*internal.Config

		if defaultConfig, err := internal.NewFileConfig(path.Join(configDirPath, "default.json")); err != nil {
			return nil, err
		} else {
			configs = append(configs, defaultConfig)
		}

		appEnv := internal.GetAppEnv()

		if len(appEnv) > 0 {
			if envConfig, err := internal.NewFileConfig(path.Join(configDirPath, fmt.Sprintf("%s.json", appEnv))); err != nil {
				return nil, err
			} else {
				configs = append(configs, envConfig)
			}

			if localEnvConfig, err := internal.NewFileConfig(path.Join(configDirPath, fmt.Sprintf("local-%s.json", appEnv))); err != nil {
				return nil, err
			} else {
				configs = append(configs, localEnvConfig)
			}
		}

		if cmdConfig, err := internal.NewCMDConfig(); err != nil {
			return nil, err
		} else {
			configs = append(configs, cmdConfig)
		}

		entireConfig := internal.JoinConfigs(configs)

		if configSchema, err := internal.NewSchema(path.Join(configDirPath, "schema")); err != nil {
			return nil, err
		} else if configSchema != nil {
			if err := configSchema.Validate(entireConfig); err != nil {
				return nil, err
			}
		}

		return json.Marshal(entireConfig)
	}
}
