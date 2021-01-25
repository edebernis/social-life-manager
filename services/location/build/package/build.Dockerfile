FROM golang:1.15-alpine

# Needed for CGO
RUN apk add --no-cache gcc g++

# Install golangci-lint tool for testing
RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b /usr/bin v1.35.2
