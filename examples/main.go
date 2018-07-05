package main

import (
	"encoding/json"
	"fmt"

	gonfig "github.com/eduardbcom/gonfig"
)

func main() {
	type Config struct {
		DbConfig struct {
			Host string `json:"host"`
			Port int    `json:"port"`
		} `json:"dbConfig"`
		Name string `json:"name"`
	}

	appConfig := &Config{}

	if rawData, err := gonfig.Read(); err != nil {
		panic(err)
	} else {
		if err := json.Unmarshal(rawData, appConfig); err != nil {
			panic(err)
		} else {
			fmt.Println(appConfig.DbConfig.Host, appConfig.DbConfig.Port) // prod-db-server 8000
		}
	}
}
