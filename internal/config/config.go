package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"

	flag "github.com/spf13/pflag"
)

var commandLineConfigObject *string
var commandLineConfigPath *string

func init() {
	flags := flag.NewFlagSet("gonfig", flag.ExitOnError)

	flags.ParseErrorsWhitelist.UnknownFlags = true

	commandLineConfigObject = flags.String("config", "{}", "config json string (default '{}')")
	commandLineConfigPath = flags.String("config_dir", "", "custom config dir")

	flags.Parse(os.Args[1:])
}

type Config struct {
	data map[string]interface{}
}

func (c *Config) IsEmpty() bool {
	return len(c.data) == 0
}

func (c *Config) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.data)
}

func NewFromFile(fileFullPath string) (*Config, error) {
	var err error
	var content []byte

	data := make(map[string]interface{})

	content, err = ioutil.ReadFile(fileFullPath)
	if err != nil && os.IsNotExist(err) {
		return &Config{data: data}, nil
	} else if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(content, &data); err != nil {
		return nil, err
	}

	return &Config{data: data}, nil
}

func NewFromCMD() (*Config, error) {
	var err error

	data := make(map[string]interface{})

	if err = json.Unmarshal([]byte(*commandLineConfigObject), &data); err != nil {
		return nil, err
	}

	return &Config{data: data}, nil
}

// GetDirPath checks that `config` directory exists near the executable file
// and returns full path to `config` directory or error.
func GetDirPath() (configDirPath string, err error) {
	configDirPath, err = getConfigDir()
	if err != nil {
		return
	}

	if _, err = os.Stat(configDirPath); err != nil {
		return
	}

	if os.IsNotExist(err) {
		return "", fmt.Errorf("config directory '%s' does not exist", configDirPath)
	}

	return
}

func Join(configs []*Config) *Config {
	data := make(map[string]interface{})

	for _, config := range configs {
		data = mergeMaps(data, config.data)
	}

	return &Config{data: data}
}

func mergeMaps(result, inputMap map[string]interface{}) map[string]interface{} {
	for k, val1 := range inputMap {
		if val2, ok := result[k]; ok == false {
			result[k] = val1
		} else if reflect.TypeOf(val1).Kind() != reflect.TypeOf(val2).Kind() {
			result[k] = val1
		} else if reflect.TypeOf(val1).Kind() == reflect.Map {
			result[k] = mergeMaps(
				result[k].(map[string]interface{}),
				val1.(map[string]interface{}),
			)
		} else {
			result[k] = val1
		}
	}

	return result
}

func getConfigDir() (string, error) {
	if len(*commandLineConfigPath) > 0 {
		return filepath.Abs(*commandLineConfigPath)
	}

	if dir, err := filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
		return "", err
	} else {
		return path.Join(dir, "config"), nil
	}
}
