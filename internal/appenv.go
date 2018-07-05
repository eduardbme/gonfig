package internal

import (
	"os"
)

func GetAppEnv() string {
	if appEnv, ok := os.LookupEnv("APP_ENV"); ok == true {
		return appEnv
	}

	return ""
}
