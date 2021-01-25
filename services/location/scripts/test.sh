#!/bin/sh

set -o errexit
set -o nounset
set -o pipefail

export CGO_ENABLED=1
export GO111MODULE=on
export GOFLAGS="-mod=vendor"

# Collect test targets
SRC_DIRS="cmd internal"
TARGETS=$(for d in ${SRC_DIRS}; do echo ./$d/...; done)

# Run tests
echo "Running unit tests:"
go test -installsuffix "static" -cover -race ${TARGETS} 2>&1
echo

# Collect all `.go` files and `gofmt` against them. If some need formatting - print them.
echo "Checking go fmt:"
ERRS=$(find ${SRC_DIRS} -type f -name \*.go | xargs gofmt -l 2>&1 || true)
if [ -n "${ERRS}" ]; then
    echo "FAIL - the following files need to be gofmt'ed:"
    for e in ${ERRS}; do
        echo "    $e"
    done
    echo
    exit 1
fi
echo "PASS"
echo

# Run `go vet`. If problems are found - print them.
echo "Checking go vet:"
ERRS=$(go vet ${TARGETS} 2>&1 || true)
if [ -n "${ERRS}" ]; then
    echo "FAIL"
    echo "${ERRS}"
    echo
    exit 1
fi
echo "PASS"
echo

# Run golangci-lint. If problems are found - print them.
echo "Checking golangci-lint:"
ERRS=$(golangci-lint run ${TARGETS} 2>&1 || true)
if [ -n "${ERRS}" ]; then
    echo "FAIL"
    echo "${ERRS}"
    echo
    exit 1
fi
echo "PASS"
