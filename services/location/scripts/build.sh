#!/bin/sh

set -o errexit
set -o nounset
set -o pipefail

# Check if OS, architecture and application version variables are set in Makefile
if [ -z "${OS:-}" ]; then
    echo "OS must be set"
    exit 1
fi
if [ -z "${ARCH:-}" ]; then
    echo "ARCH must be set"
    exit 1
fi
if [ -z "${VERSION:-}" ]; then
    echo "VERSION must be set"
    exit 1
fi

# Disable C code, enable Go modules
export CGO_ENABLED=0
export GOARCH="${ARCH}"
export GOOS="${OS}"
export GOFLAGS="-mod=vendor"

echo "Download dependencies:"
go mod vendor

echo "Compiling proto files:"
protoc                                      \
    --proto_path=./vendor                   \
    --proto_path=.                          \
    --go_out=.                              \
    --go_opt=paths=source_relative          \
    --go-grpc_out=.                         \
    --go-grpc_opt=paths=source_relative     \
    --govalidators_out=.                    \
    api/grpc/v1/location.proto

echo "Building app:"
go install                                  \
    -installsuffix "static"                 \
    ./...

echo "Building API docs:"
swag init                                   \
    -g internal/api/http/v1/server.go       \
    -o ./api/http/v1
