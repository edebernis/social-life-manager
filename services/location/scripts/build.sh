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

echo "Building app:"
go install                      \
    -installsuffix "static"     \
    ./...
echo "OK"

echo "Building API docs:"
swag init                           \
    -g internal/api/http/server.go  \
    -o ./api
echo "OK"
