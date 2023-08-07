FROM ubuntu:18.04 AS runtime
ENV GIN_MODE=release
ENV PORT=8080
WORKDIR /app
# GoReleaser will automatically generate the binary in the root directory
COPY /go-web-template .
EXPOSE 8080
ENTRYPOINT ["./go-web-template"]