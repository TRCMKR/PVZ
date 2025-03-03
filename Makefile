BINARY=app

BUILD_FLAGS=-ldflags="-s -w"

GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)

## builds app + clean + fmt + lint
build: clean fmt lint
	go build $(BUILD_FLAGS) -o ./build/$(BINARY) ./cmd/app

## runs built app
run:
	./build/$(BINARY)

## cleans previous builds
clean:
	rm -f ./build/*

## fixes formatting
fmt:
	go fmt ./...

## runs linters + vet
lint: vet
	@if ! command -v golangci-lint &> /dev/null; then \
		echo "golangci-lint не найден, установите его: https://golangci-lint.run/"; \
		exit 1; \
	fi
	golangci-lint run -c .golangci.yml

## runs go vet
vet:
	go vet ./...

## installs linter
install-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

## installs dependencies
deps:
	go mod download

## updates dependencies
update:
	go get -u ./...

## runs go mod tidy
tidy:
	go mod tidy

## runs go clean -modcache
cache:
	go clean -modcache

## builds app for win
build-windows:
	GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) -o ./build/$(BINARY).exe

## prints help about all targets
help:
	@echo ""
	@echo "Usage:"
	@echo "  make <target>"
	@echo ""
	@echo "Targets:"
	@awk '                                \
		BEGIN { comment=""; }             \
		/^\s*##/ {                         \
		    comment = substr($$0, index($$0,$$2)); next; \
		}                                  \
		/^[a-zA-Z0-9_-]+:/ {               \
		    target = $$1;                  \
		    sub(":", "", target);          \
		    if (comment != "") {           \
		        printf "  %-15s %s\n", target, comment; \
		        comment="";                \
		    }                              \
		}' $(MAKEFILE_LIST)
	@echo ""

.PHONY: build run clean fmt deps install-lint lint tidy build-windows help cache