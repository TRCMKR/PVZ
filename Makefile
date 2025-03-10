BINARY=app

BUILD_FLAGS=-ldflags="-s -w"

ifneq (,$(wildcard .env))
	include .env
endif

GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)

ifeq ($(POSTGRES_SETUP_TEST),)
	POSTGRES_SETUP_TEST := user=$(DB_USERNAME) password=$(DB_PASSWORD) dbname=$(DB_NAME) host=$(DB_HOST) port=$(DB_PORT) sslmode=disable
endif

MIGRATION_FOLDER=$(CURDIR)/migrations

.PHONY: build
## builds app + clean + fmt + lint
build: fmt lint clean
	go build $(BUILD_FLAGS) -o ./build/$(BINARY) ./cmd/app

.PHONY: run
## runs built app
run:
	./build/$(BINARY)

.PHONY: clean
## cleans previous builds
clean:
	rm -f ./build/*

.PHONY: fmt
## fixes formatting
fmt:
	go fmt ./...

.PHONY: lint
## runs linters + vet
lint: vet
	@if ! command -v golangci-lint &> /dev/null; then \
		echo "golangci-lint не найден, установите его: https://golangci-lint.run/"; \
		exit 1; \
	fi
	golangci-lint run -c .golangci.yml

.PHONY: vet
## runs go vet
vet:
	go vet ./...

.PHONY: install-lint
## installs linter
install-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: deps
## installs dependencies
deps:
	go mod download

.PHONY: update
## updates dependencies
update:
	go get -u ./...

.PHONY: tidy
## runs go mod tidy
tidy:
	go mod tidy

.PHONY: cache
## runs go clean -modcache
cache:
	go clean -modcache

.PHONY: build-windows
## builds app for win
build-windows:
	GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) -o ./build/$(BINARY).exe

## creates migration with first param as name
.PHONY: migration-create
migration-create:
	goose -dir "$(MIGRATION_FOLDER)" create $(name) sql

## applies latest migration
.PHONY: migration-up
migration-up:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" up

## rolls back latest migration
.PHONY: migration-down
migration-down:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" down

## checks migration status
.PHONY: migration-status
migration-status:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" status

.PHONY: swag-init
swag-init:
	swag init -g "./internal/web/router.go" --parseInternal --pd

.PHONY: help
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
		        printf "  %-17s %s\n", target, comment; \
		        comment="";                \
		    }                              \
		}' $(MAKEFILE_LIST)
	@echo ""