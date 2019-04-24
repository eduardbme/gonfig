package gonfig

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/eduardbcom/gonfig/internal/config"
	"github.com/eduardbcom/gonfig/internal/schema"
)

// Read the configuration json data from
// - config/default.json
// - config/[APP_ENV].json
// - config/local-[APP_END].json
// - CMD --config='{}'
// and merge it into one object.
func Read() ([]byte, error) {
	if configDirPath, err := config.GetDirPath(); err != nil {
		return nil, err
	} else {
		var configs []*config.Config

		if defaultConfig, err := config.NewFromFile(path.Join(configDirPath, "default.json")); err != nil {
			return nil, err
		} else {
			configs = append(configs, defaultConfig)
		}

		if appEnv, ok := os.LookupEnv("APP_ENV"); ok {
			if envConfig, err := config.NewFromFile(path.Join(configDirPath, fmt.Sprintf("%s.json", appEnv))); err != nil {
				return nil, err
			} else {
				configs = append(configs, envConfig)
			}

			if localEnvConfig, err := config.NewFromFile(path.Join(configDirPath, fmt.Sprintf("local-%s.json", appEnv))); err != nil {
				return nil, err
			} else {
				configs = append(configs, localEnvConfig)
			}
		}

		if cmdConfig, err := config.NewFromCMD(); err != nil {
			return nil, err
		} else {
			configs = append(configs, cmdConfig)
		}

		entireConfig := config.Join(configs)

		if entireConfig.IsEmpty() {
			fmt.Fprintf(os.Stderr, "[GONFIG][WARNING] no configuration data\n")
		}

		if configSchema, err := schema.New(); err != nil {
			return nil, err
		} else if configSchema != nil {
			if err := configSchema.Validate(entireConfig); err != nil {
				return nil, err
			}
		}

		return json.Marshal(entireConfig)
	}
}
