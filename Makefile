# Default variable
APP_NAME 			?= go-web-prototype
GOFORMATER			?= gofumpt
GOFORMATER_VERSION	?= v0.2.0
GOLINTER			?= golangci-lint
GOLINTER_VERSION	?= v1.41.0
COVER_FILE			?= cover.out
SOURCE_PATHS		?= ./cmd/...

.DEFAULT_GOAL := help

help:  ## This help message :)
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# If you encounter an error like "panic: permission denied" on MacOS,
# please visit https://github.com/eisenxp/macos-golink-wrapper to find the solution.
test:  ## Run the tests
	go test -gcflags=all=-l -timeout=10m `go list $(SOURCE_PATHS)` ${TEST_FLAGS}

cover:  ## Generates coverage report
	go test -gcflags=all=-l -timeout=10m `go list $(SOURCE_PATHS)` -coverprofile $(COVER_FILE) ${TEST_FLAGS}

cover-html:  ## Generates coverage report and displays it in the browser
	go tool cover -html=$(COVER_FILE)

format:  ## Format source code
	@which $(GOFORMATER) > /dev/null || (echo "Installing $(GOFORMATER)@$(GOFORMATER_VERSION) ..."; go install mvdan.cc/gofumpt@$(GOFORMATER_VERSION) && echo -e "Installation complete!\n")
	@for path in $(SOURCE_PATHS); do ${GOFORMATER} -l -w -e `echo $${path} | cut -b 3- | rev | cut -b 5- | rev`; done;

lint:  ## Lint, will not fix but sets exit code on error
	@which $(GOLINTER) > /dev/null || (echo "Installing $(GOLINTER)@$(GOLINTER_VERSION) ..."; go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLINTER_VERSION) && echo -e "Installation complete!\n")
	$(GOLINTER) run --deadline=10m $(SOURCE_PATHS)

lint-fix:  ## Lint, will try to fix errors and modify code
	@which $(GOLINTER) > /dev/null || (echo "Installing $(GOLINTER)@$(GOLINTER_VERSION) ..."; go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLINTER_VERSION) && echo -e "Installation complete!\n")
	$(GOLINTER) run --deadline=10m $(SOURCE_PATHS) --fix

doc:  ## Start the documentation server with godoc
	@which godoc > /dev/null || (echo "Installing godoc@latest ..."; go install golang.org/x/tools/cmd/godoc@latest && echo -e "Installation complete!\n")
	godoc -http=:6060

build-changelog:  ## Build the changelog with git-chglog
	@which git-chglog > /dev/null || (echo "Installing git-chglog@v0.15.1 ..."; go install github.com/git-chglog/git-chglog/cmd/git-chglog@v0.15.1 && echo -e "Installation complete!\n")
	mkdir -p ./build
	git-chglog -o ./build/CHANGELOG.md

clean:  ## Clean build bundles
	-rm -rf ./build

gen-version: ## Update version
	cd pkg/version/scripts && go run gen.go

build-all: build-darwin build-linux build-windows ## Build all platforms

build-darwin: gen-version ## Build for MacOS
	-rm -rf ./build/darwin/bin
	# Build kusion
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 \
		go build -o ./build/darwin/bin/$(APP_NAME) \
		./cmd/main.go

build-linux: gen-version ## Build for Linux
	-rm -rf ./build/linux/bin
	# Build kusion
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
		go build -o ./build/linux/bin/$(APP_NAME) \
		./cmd/main.go

build-windows: gen-version ## Build for Windows
	-rm -rf ./build/windows/bin
	# Build kusion
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 \
		go build -o ./build/windows/bin/$(APP_NAME).exe \
		./cmd/main.go

.PHONY: test cover cover-html format lint lint-fix doc build-changelog clean gen-version build-all build-darwin build-linux build-windows
