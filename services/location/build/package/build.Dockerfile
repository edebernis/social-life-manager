FROM golang:1.15-alpine

# Install compilers and git
# - C compilers used for CGO
# - Protobuf compiler
RUN apk add --no-cache gcc g++ git protoc=3.13.0-r2

# Enable Go Modules
ENV GO111MODULE on

# Install Go tools
# - Protobuf Go plugins
# - Swag tool to build HTTP API docs
RUN GOBIN=/usr/local/go/bin go get -u google.golang.org/protobuf/cmd/protoc-gen-go@v1.25.0 \
                                      google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0 \
                                      github.com/swaggo/swag/cmd/swag@v1.7.0

# Install golangci-lint tool for testing
RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b /usr/local/bin v1.35.2
