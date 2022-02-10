## Introduction


[![GitHub release](https://img.shields.io/github/release/elliotxx/go-web-prototype.svg)](https://github.com/elliotxx/go-web-prototype/releases)
[![Github All Releases](https://img.shields.io/github/downloads/elliotxx/go-web-prototype/total.svg)](https://github.com/elliotxx/go-web-prototype/releases)
[![license](https://img.shields.io/github/license/elliotxx/go-web-prototype.svg)](https://github.com/elliotxx/go-web-prototype/blob/master/LICENSE)
[![Coverage Status](https://coveralls.io/repos/github/elliotxx/go-web-prototype/badge.svg)](https://coveralls.io/github/elliotxx/go-web-prototype)

> This is a web application with go and gin.

## Usage
Local startup:
```
$ go run cmd/main.go
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /ping                     --> main.setupRouter.func1 (3 handlers)
[GIN-debug] GET    /user/:name               --> main.setupRouter.func2 (3 handlers)
[GIN-debug] POST   /admin                    --> main.setupRouter.func3 (4 handlers)
[GIN-debug] Listening and serving HTTP on :8080

$ curl http://127.0.0.1:8080/ping 
pong
```

Local build:
```
$ make build-all
$ ./build/darwin/go-web-prototype
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /ping                     --> main.setupRouter.func1 (3 handlers)
[GIN-debug] GET    /user/:name               --> main.setupRouter.func2 (3 handlers)
[GIN-debug] POST   /admin                    --> main.setupRouter.func3 (4 handlers)
[GIN-debug] Listening and serving HTTP on :8080

$ curl http://127.0.0.1:8080/ping 
pong
```

Run all unit tests:
```
make test
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
```