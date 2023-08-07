GOSOURCE_PATHS ?= ./cmd/...

include go.mk


.PHONY: clean
clean:  ## Clean build bundles
	-rm -rf ./_build

.PHONY: build-all
build-all: build-darwin build-linux build-windows ## Build for all platforms

.PHONY: build-darwin
build-darwin: gen-version ## Build for MacOS
	-rm -rf ./_build/darwin
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 \
		go build -o ./_build/darwin/$(APPROOT) \
		./cmd/main.go

.PHONY: build-linux
build-linux: gen-version ## Build for Linux
	-rm -rf ./_build/linux
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
		go build -o ./_build/linux/$(APPROOT) \
		./cmd/main.go

.PHONY: build-windows
build-windows: gen-version ## Build for Windows
	-rm -rf ./_build/windows
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 \
		go build -o ./_build/windows/$(APPROOT).exe \
		./cmd/main.go

.PHONY: gen-api-docs
gen-api-docs: ## Generate API documentation with OpenAPI format
	@which swag > /dev/null || (echo "Installing swag@v1.7.8 ..."; go install github.com/swaggo/swag/cmd/swag@v1.7.8 && echo "Installation complete!\n")
	# Generate API documentation with OpenAPI format
	-swag init --parseDependency --parseDepth 1 -g cmd/main.go -o api/openapispec/
	# Format swagger comments
	-swag fmt -g pkg/**/*.go
	@echo "🎉 Done!"

.PHONY: gen-version
gen-version: ## Generate version file
	# Delete old version file
	-rm -f ./pkg/version/z_update_version.go
	# Update version
	-cd pkg/version/scripts && go run gen/gen.go
