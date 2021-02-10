VERSION := $(shell git describe --tags --always --dirty)

ALL_PLATFORMS := linux/amd64 darwin/amd64
OS := $(if $(GOOS),$(GOOS),$(shell go env GOOS))
ARCH := $(if $(GOARCH),$(GOARCH),$(shell go env GOARCH))

BIN := location
BUILD_DIR := bin/$(OS)_$(ARCH)
OUTBIN := $(BUILD_DIR)/$(BIN)

REGISTRY ?= local
IMAGE := $(REGISTRY)/$(BIN)
TAG := $(VERSION)__$(OS)_$(ARCH)
BUILD_IMAGE := build/$(BIN)

build-%:
	@$(MAKE) build                        \
	    --no-print-directory              \
	    GOOS=$(firstword $(subst _, ,$*)) \
	    GOARCH=$(lastword $(subst _, ,$*))

container:
	@$(MAKE) .container                   \
	    --no-print-directory              \
	    GOOS=linux                        \
	    GOARCH=amd64

push:
	@$(MAKE) .push                        \
	    --no-print-directory              \
	    GOOS=linux                        \
	    GOARCH=amd64

all-build: $(addprefix build-, $(subst /,_, $(ALL_PLATFORMS)))

build: $(OUTBIN)
$(OUTBIN):
	@$(MAKE) run-build-script           \
	    --no-print-directory            \
	    SCRIPT=scripts/build.sh

test:
	@$(MAKE) run-build-script          \
	    --no-print-directory           \
	    SCRIPT=scripts/test.sh

build-image:
	@docker build -t $(BUILD_IMAGE) - < build/package/build.Dockerfile

run-build-script: build-image
	@mkdir -p $(BUILD_DIR) .cache
	@docker run                                                 \
	    -i                                                      \
	    --rm                                                    \
	    -u $$(id -u):$$(id -g)                                  \
	    -v $$(pwd):/src                                         \
	    -w /src                                                 \
		-v $$(pwd)/$(BUILD_DIR):/go/bin                         \
	    -v $$(pwd)/$(BUILD_DIR):/go/bin/$(OS)_$(ARCH)           \
		-v $$(pwd)/.cache:/.cache                               \
	    --env HTTP_PROXY=$(HTTP_PROXY)                          \
	    --env HTTPS_PROXY=$(HTTPS_PROXY)                        \
	    $(BUILD_IMAGE)                                          \
	    /bin/sh -c "                                            \
	        ARCH=$(ARCH)                                        \
	        OS=$(OS)                                            \
	        VERSION=$(VERSION)                                  \
	        $(SCRIPT)                                           \
	    "

.container: $(OUTBIN)
	@docker build --build-arg bin=$(BIN) -t $(IMAGE):$(TAG) -t $(IMAGE):latest -f build/package/run.Dockerfile $(BUILD_DIR)

.push: .container
	@docker push $(IMAGE):$(TAG)
	@docker push $(IMAGE):latest

clean-bin:
	@rm -rf bin

clean-cache:
	@rm -rf .cache

clean: clean-bin clean-cache

version:
	@echo $(VERSION)

protobuf:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative api/grpc/v1/location.proto

migrate.up:
	migrate -database $(POSTGRESQL_URL) -path internal/repositories/sql/migrations up

migrate-docker-compose.up:
	make migrate.up POSTGRESQL_URL=postgres://postgres:password@localhost:5432/location?sslmode=disable

migrate.down:
	migrate -database $(POSTGRESQL_URL) -path internal/repositories/sql/migrations down

migrate-docker-compose.down:
	make migrate.down POSTGRESQL_URL=postgres://postgres:password@localhost:5432/location?sslmode=disable

jwt-token:
	jwt encode --alg HS256 --secret secret --sub $(shell uuidgen)
