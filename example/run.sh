#!/bin/sh

go build main.go

APP_ENV=example ./main --config='{"name": "new-awesome-name"}'
