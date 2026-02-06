APP_NAME := goMuffin
BIN_DIR := build

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
	@go build -o $(BIN) $(PKG)

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
