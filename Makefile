APP_NAME := gh-crossplane
BINARY_NAME := $(APP_NAME)
BUILD_DIR := bin
PKG := github.com/artemlive/$(APP_NAME)
VERSION := $(shell git describe --tags --always --dirty)
LDFLAGS := -s -w -X $(PKG)/cmd.version=$(VERSION)

GO_FILES := $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: all build clean test install lint run

all: build

## Build the binary
build:
	@echo "Building..."
	@go build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/app/main.go

## Run the app (for dev/debug purposes)
run:
	@go run ./cmd/app/main.go

## Format code
fmt:
	@gofmt -w $(GO_FILES)

## Run linters (go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
lint:
	@golangci-lint run

## Run tests
test:
	@go test ./... -v

## Install binary to your system
install:
	@go install -ldflags="$(LDFLAGS)" ./cmd/app/main.go

## Clean build artifacts
clean:
	@rm -rf $(BUILD_DIR)

## Show the current version
version:
	@echo $(VERSION)
