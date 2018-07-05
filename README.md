Configure your Go Applications
===================================

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
```

```shell
$ export APP_ENV=production
$ go run app.go
```

File formats
---------------
For now only `json` file format is supported.

Schema validation
---------------
In order to validate configuration using `json schema` format you need to create `schema` folder within the `config` folder.

`schema/index.json` is the entry point of schema.

Example
---------------
See the example folder for, well, examples.
