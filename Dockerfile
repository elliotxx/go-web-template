FROM golang:1.17 AS builder

WORKDIR /app
ADD . .
ENV GO111MODULE=on

RUN ["/bin/bash", "-c", "go mod tidy"]
RUN ["/bin/bash", "-c", "go build -o standard-sample-go-web ./cmd"]


FROM ubuntu:18.04 AS runtime

ENV GIN_MODE=release
ENV PORT=8080

WORKDIR /app
COPY --from=builder /app/standard-sample-go-web .

EXPOSE 8080
ENTRYPOINT ["./standard-sample-go-web"]