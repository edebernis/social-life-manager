FROM golang:1.15-alpine

# Needed for CGO and go get
RUN apk add --no-cache gcc g++ git

# Enable Go Modules
ENV GO111MODULE on

# Install golangci-lint tool for testing
RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b /usr/local/bin v1.35.2

# Install swag tool to build HTTP API docs
RUN GOBIN=/usr/local/bin go get -u github.com/swaggo/swag/cmd/swag@v1.7.0
