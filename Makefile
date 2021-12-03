# Default variable
GOIMPORTS ?= goimports
GOCILINT ?= golangci-lint

default:
	@go run cmd/main.go

test:
	@go test -timeout=10m `go list ./pkg/... ./cmd/...`

test-with-coverage:
	@go test -v -coverprofile=profile.cov -timeout=10m `go list ./pkg/... ./cmd/...`

lint:
	@$(GOCILINT) run --no-config --disable=errcheck ./...