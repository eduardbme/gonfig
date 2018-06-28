package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
)

// Read ...
func Read(configDirFullPath string, config interface{}) (err error) {
	if err = checkConfigFolder(configDirFullPath); err != nil {
		return
	}

	appEnv := getAppEnv()

	defaultConfigData := make(map[string]interface{})
	envConfigData := make(map[string]interface{})
	localEnvConfigData := make(map[string]interface{})

	defaultConfigData = readJSONFile(configDirFullPath, "default.json")

	if len(appEnv) > 0 {
		envConfigData = readJSONFile(configDirFullPath, fmt.Sprintf("%s.json", appEnv))
		localEnvConfigData = readJSONFile(configDirFullPath, fmt.Sprintf("local-%s.json", appEnv))
	}

	joinedConfigData := joinConfigData(defaultConfigData, envConfigData, localEnvConfigData)

	var data []byte
	data, err = json.Marshal(joinedConfigData)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return
	}

	return
}

func checkConfigFolder(configDirFullPath string) (err error) {
	if _, err = os.Stat(configDirFullPath); err != nil {
		if os.IsNotExist(err) {
			err = errors.New("config directory does not exist")
		}
	}

	return
}

func getAppEnv() string {
	appEnv, ok := os.LookupEnv("APP_ENV")

	if ok != true {
		return ""
	}

	return appEnv
}

func readJSONFile(configDirFullPath, filename string) (configData map[string]interface{}) {
	configData = make(map[string]interface{})

	fileFullPath := filepath.Join(configDirFullPath, filename)

	if content, e := ioutil.ReadFile(fileFullPath); e != nil {
		if os.IsNotExist(e) == false {
			panic(e)
		}
	} else {
		json.Unmarshal(content, &configData)
	}

	return configData
}

func joinConfigData(configs ...map[string]interface{}) (result map[string]interface{}) {
	result = make(map[string]interface{})

	for _, configData := range configs {
		result = mergeMaps(result, configData)
	}

	return result
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
