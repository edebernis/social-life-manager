#!/bin/sh

set -o errexit
set -o nounset
set -o pipefail

export CGO_ENABLED=1
export GOFLAGS="-mod=vendor"

echo "Download dependencies:"
go mod vendor

# Collect test targets
SRC_DIRS="cmd pkg internal"
TARGETS=$(for d in ${SRC_DIRS}; do echo -n "./$d/... "; done)

# Lint everything
echo "Running linters:"
golangci-lint run -E gofmt --timeout 10m ${TARGETS} 2>&1

# Run tests
echo "Running unit tests:"
go test -installsuffix "static" -cover -race ${TARGETS} 2>&1
