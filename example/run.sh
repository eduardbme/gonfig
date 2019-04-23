#!/bin/sh

go build main.go

APP_ENV=example ./main --config_dir='/tmp'
APP_ENV=example ./main --config_dir="`pwd`/config"
APP_ENV=example ./main --config='{"name": "new-awesome-name"}' --config_dir='/tmp'
APP_ENV=example ./main --config='{"name": "new-awesome-name"}' --config_dir="`pwd`/config"
