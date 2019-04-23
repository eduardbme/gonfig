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

var commandLineConfig *string

func init() {
	flags := flag.NewFlagSet("gonfig", flag.ExitOnError)

	flags.ParseErrorsWhitelist.UnknownFlags = true

	commandLineConfig = flags.String("config", "{}", "config json string (default '{}')")

	flags.Parse(os.Args[1:])
}

type Config struct {
	data map[string]interface{}
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

	if err = json.Unmarshal([]byte(*commandLineConfig), &data); err != nil {
		return nil, err
	}

	return &Config{data: data}, nil
}

// GetDirPath checks that `config` directory exists near the executable file
// and returns full path to `config` directory or error.
func GetDirPath() (configDirPath string, err error) {
	var dir string

	dir, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return
	}

	configDirPath = path.Join(dir, "config")

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
