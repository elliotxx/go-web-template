## Introduction


[![GitHub release](https://img.shields.io/github/release/elliotxx/go-web-template.svg)](https://github.com/elliotxx/go-web-template/releases)
[![Github All Releases](https://img.shields.io/github/downloads/elliotxx/go-web-template/total.svg)](https://github.com/elliotxx/go-web-template/releases)
[![license](https://img.shields.io/github/license/elliotxx/go-web-template.svg)](https://github.com/elliotxx/go-web-template/blob/master/LICENSE)
[![Coverage Status](https://coveralls.io/repos/github/elliotxx/go-web-template/badge.svg)](https://coveralls.io/github/elliotxx/go-web-template)

> Template application of Domain Driven Design(DDD) in go and gin.

## Usage
Prepare database:
```shell
$ docker pull mariadb:10.5
$ docker run -d --env MARIADB_USER=app --env MARIADB_PASSWORD=app123 --env MARIADB_ROOT_PASSWORD=123456 -p 127.0.0.1:3306:3306 mariadb:10.5
$ mysql -h 127.0.0.1 -u root -p123456
    > CREATE DATABASE IF NOT EXISTS appdb CHARACTER SET utf8mb4;
    > DROP USER 'app';
    > CREATE USER IF NOT EXISTS 'app' IDENTIFIED BY 'app123';
    > GRANT ALL ON appdb.* TO `app`;
    > exit
```

Local startup:
```
$ go run cmd/main.go -f config/local.yaml
```

Local verification:
```
➜ curl http://localhost:80/livez    
OK

➜ curl http://localhost:80/readyz
[+] Database ok
[+] Ping ok
health check passed

➜ curl http://localhost:80/endpoints
DELETE  /api/v1/systemconfig/:id
GET     /api/v1/systemconfig/:id
GET     /api/v1/systemconfig/count
GET     /api/v1/systemconfigs
GET     /debug/pprof/
GET     /debug/pprof/allocs
GET     /debug/pprof/block
GET     /debug/pprof/cmdline
GET     /debug/pprof/goroutine
GET     /debug/pprof/heap
GET     /debug/pprof/mutex
GET     /debug/pprof/profile
GET     /debug/pprof/symbol
GET     /debug/pprof/threadcreate
GET     /debug/pprof/trace
GET     /debug/statsviz/*filepath
GET     /debug/vars
GET     /docs/*any
GET     /livez
GET     /readyz
POST    /api/v1/systemconfig
POST    /debug/pprof/symbol
PUT     /api/v1/systemconfig

➜ curl http://localhost:80/debug/vars
{
    "appOptions": {
        "database": {
            "autoMigrate": true,
            "dbHost": "127.0.0.1",
            "dbName": "appdb",
            "dbPassword": "******",
            "dbPort": 3306,
            "dbUser": "app",
            "migrateFile": "./assets/sql/app.sql"
        },
        "generic": {
            "configFile": "config/local.yaml",
            "dumpEnvs": false,
            "dumpVersion": false
        },
        "logging": {
            "disableText": true,
            "dumpCurrentConfig": true,
            "enableLoggingToFile": true,
            "jsonPretty": true,
            "logLevel": "debug",
            "loggingDirectory": "logs",
            "reportCaller": true
        },
        "network": {
            "port": 80,
            "requestTimeout": 30000000000
        }
    },
    "version": {
        "buildInfo": {
            "GOARCH": "amd64",
            "GOOS": "darwin",
            "buildTime": "2023-08-07 17:29:54",
            "compiler": "gc",
            "goVersion": "go1.19.9",
            "numCPU": 8
        },
        "gitInfo": {
            "commit": "87a0dd63b10473eedf5660f6e319492647d29389",
            "latestTag": "v0.2.0",
            "treeState": "dirty"
        },
        "releaseVersion": "v0.2.0-87a0dd63"
    }
}

➜ curl -s http://localhost:80/api/v1/systemconfig/count | jq
{
  "success": true,
  "code": "00000",
  "message": "OK",
  "data": {
    "total": 0
  },
  "traceID": "754a40b1-1a80-43ff-adc5-2caf65de8db8",
  "startTime": "2023-08-07T17:59:37.74405+08:00",
  "endTime": "2023-08-07T17:59:37.746765+08:00",
  "costTime": "2.715053ms"
}

➜ curl -s --request POST 'http://localhost:80/api/v1/systemconfig' \       
--header 'Content-Type: application/json' \
--data '{
  "tenant": "MAIN_SITE",
  "env": "prod",
  "type": "config",
  "config": "{'\''abc'\'': '\''xxx'\''}",
  "description": "config",
  "creator": "elliotxx",
  "modifier": ""
}' | jq
{
  "success": true,
  "code": "00000",
  "message": "OK",
  "data": {
    "id": 1400004,
    "tenant": "MAIN_SITE",
    "env": "prod",
    "type": "config",
    "config": "{'abc': 'xxx'}",
    "description": "config",
    "creator": "elliotxx",
    "createdAt": "2023-08-07T18:01:05.32+08:00",
    "updatedAt": "2023-08-07T18:01:05.32+08:00"
  },
  "traceID": "129f98e8-cf45-4060-90b2-05b861128ff3",
  "startTime": "2023-08-07T18:01:05.315387+08:00",
  "endTime": "2023-08-07T18:01:05.330384+08:00",
  "costTime": "14.996752ms"
}

➜ curl -s --location --request GET 'http://localhost:80/api/v1/systemconfigs' \
--header 'Content-Type: application/json' \
--data '{    
    "page": 1,
    "perPage": 3,
    "keyword": ""
}' | jq
{
  "success": true,
  "code": "00000",
  "message": "OK",
  "data": [
    {
      "id": 1400004,
      "tenant": "MAIN_SITE",
      "env": "prod",
      "type": "config",
      "config": "{'abc': 'xxx'}",
      "description": "config",
      "creator": "elliotxx",
      "createdAt": "2023-08-07T18:01:05.32+08:00",
      "updatedAt": "2023-08-07T18:01:05.32+08:00"
    }
  ],
  "traceID": "ba332ec3-36ff-4a18-a835-83b28f50d7fa",
  "startTime": "2023-08-07T18:03:31.790936+08:00",
  "endTime": "2023-08-07T18:03:31.798123+08:00",
  "costTime": "7.18668ms"
}
```

Local build:
```
$ make build-all
```

Run all unit tests:
```
make cover
```

All targets:
```
$ make help
help                           This help message :)
test                           Run the tests
cover                          Generates coverage report
cover-html                     Generates coverage report and displays it in the browser
format                         Format source code
lint                           Lint, will not fix but sets exit code on error
lint-fix                       Lint, will try to fix errors and modify code
doc                            Start the documentation server with godoc
clean                          Clean build bundles
build-all                      Build for all platforms
build-darwin                   Build for MacOS
build-linux                    Build for Linux
build-windows                  Build for Windows
gen-api-docs                   Generate API documentation with OpenAPI format
gen-version                    Generate version file
```
