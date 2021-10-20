# Step 1: Modules caching
FROM golang:1.16.4-alpine3.13 as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
ENV GOPROXY=https://goproxy.cn,direct
RUN go mod download

# Step 2: Builder
FROM golang:1.16.4-alpine3.13 as builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN go build -o ./standard-sample-go-web ./cmd

# Step 3: Runtime
FROM ubuntu:18.04 AS runtime
ENV GIN_MODE=release
ENV PORT=8080
WORKDIR /app
COPY --from=builder /app/standard-sample-go-web .
EXPOSE 8080
ENTRYPOINT ["./standard-sample-go-web"]