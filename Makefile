default:
	go run cmd/main.go

test:
	go test -timeout=10m -v `go list ./pkg/... ./cmd/...`