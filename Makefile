default:
	go run cmd/main.go

test:
	go test -timeout=10m `go list ./pkg/... ./cmd/...`