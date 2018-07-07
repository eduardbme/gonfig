Configure your Go Applications
===================================

It lets you define a set of default parameters, and extend them for different deployment environments (development, qa, staging, production, etc.).

Quick Start
---------------
```shell
$ mkdir config
$ vi config/default.json
```

```js
{
    "name": "app-name",
    "dbConfig": {
      "host": "localhost",
      "port": 1
    }
}
```

```shell
 $ vi config/production.json
```

```json
{
    "dbConfig": {
      "host": "prod-db-server"
    }
}
```

```shell
 $ vi config/local-production.json
```

```json
{
    "dbConfig": {
      "port": 8000
    }
}
```

**Use configs in your code:**

```go
import (
    gonfig "github.com/eduardbcom/gonfig"
)

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
		fmt.Printf(
            "{\"name\": \"%s\", \"dbConfig\": {\"host\": \"%s\", port: \"%d\"}}\n",
            appConfig.Name,
            appConfig.DbConfig.Host,
            appConfig.DbConfig.Port
        ) // {"name": "new-awesome-name", "dbConfig": {"host": "prod-db-server", port: "1"}}
    }
}
```

```shell
$ export APP_ENV=production
$ go run app.go --config='{"name": "new-app-name"}'
```

Config file formats
---------------
Only `json` file format is supported.

Schema validation
---------------
In order to validate configuration using `json schema` format you need to create `schema` folder within the `config` folder.

`schema/index.json` is the entry point of schema.

Only `draft-04` is supported. Under the hood [validate-json](https://github.com/cesanta/validate-json) is used for schema validation.

Example
---------------
See the example folder.
