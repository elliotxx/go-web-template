FROM golang:1.17 AS builder

ADD . /app/
WORKDIR /app/

RUN go mod tidy
RUN go build -o standard-sample-go-web ./cmd


FROM ubuntu:18.04 AS runtime

ENV GIN_MODE=release
ENV PORT=8080

RUN mkdir /app
WORKDIR /app
COPY --from=builder /app/standard-sample-go-web .

EXPOSE 8080
ENTRYPOINT ["./standard-sample-go-web"]