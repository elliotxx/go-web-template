FROM ubuntu:18.04 AS runtime
ENV GIN_MODE=release
ENV PORT=8080
WORKDIR /app
# GoReleaser will automatically generate the binary in the root directory
COPY /standard-sample-go-web .
EXPOSE 8080
ENTRYPOINT ["./standard-sample-go-web"]