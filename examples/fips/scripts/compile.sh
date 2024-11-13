#!/usr/bin/env bash
echo start of compile.sh
go version

cd /mnt
pwd
ls -la
go clean -modcache -cache
go mod download

# Build the Go app
CGO_ENABLED=1 GOOS=linux GOARCH=arm64 GOEXPERIMENT=boringcrypto go build -tags boringcrypto -o fips_web_server_linux ./fips_web_server.go

ls -la