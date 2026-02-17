APP_NAME := goMuffin
BIN_DIR := build

VERSION := 8.0.0
CODENAME := tiramisu

CONFIG_PKG := $(shell go list ./internal/configs)

BRANCH := $(shell git branch --show-current)
DATE := $(shell date +%y%m%d%H%M)
COMMIT_HASH := $(shell git rev-parse --short HEAD)
LD_FLAGS := -ldflags "-X '$(CONFIG_PKG).MuffinVersion=$(VERSION)-$(CODENAME)_$(BRANCH).$(DATE).$(COMMIT_HASH)' \
-X '$(CONFIG_PKG).updatedString=$(DATE)'"

EXT :=
ifeq ($(OS),Windows_NT)
	EXT := .exe
endif

BIN := $(BIN_DIR)/$(APP_NAME)$(EXT)
PKG := ./internal/cmd/bot

.PHONY: all build run fmt vet deps

all: build

build:
	@mkdir -p $(BIN_DIR)
	@go build $(LD_FLAGS) -o $(BIN) $(PKG)

run:
	@go run $(PKG)

clean:
	@rm -rf $(BIN_DIR)

deps:
	@go mod tidy

fmt:
	@go fmt $(PKG)

vet:
	@go vet $(PKG)

migration:
	@go run ./internal/cmd/migration
